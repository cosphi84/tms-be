package users

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
	"tms-be/internal/auth"
	"tms-be/internal/casbin"
	"tms-be/internal/pagination"

	"gorm.io/gorm"
)

type UserService interface {
	Create(context.Context, *CreateUserDTO) error
	AllUsers(context.Context, pagination.DtoPaginationRequest) (*pagination.DtoPaginationResponse, error)
	Update(context.Context, uint64, *UpdateUserDTO) error
	Activate(context.Context, uint64) error
	Deactivate(context.Context, uint64) error
	Delete(context.Context, uint64) error
}

type userService struct {
	usrRepo   Repository
	authorize *casbin.Service
}

func NewUserService(r Repository, a *casbin.Service) UserService {
	return &userService{
		usrRepo:   r,
		authorize: a,
	}
}

func (u *userService) Create(ctx context.Context, usr *CreateUserDTO) error {
	exists, err := u.usrRepo.IsExists(ctx, usr.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user with email %s already registered", usr.Email)
	}

	loggedUser, err := auth.GetClaims(ctx)
	if err != nil {
		return err
	}

	// auth.HashPassword — dipindah dari helpers ke auth buat hindarin
	// dependency cycle.
	hashed, err := auth.HashPassword(usr.Password)
	if err != nil {
		return err
	}

	newUsr := Model{
		Username:    usr.Name,
		Email:       usr.Email,
		Password:    hashed,
		OfficeID:    usr.OfficeID,
		FotoProfile: &usr.Image,
		IsActive:    true,
		CreatedAt:   time.Now(),
		CreatedBy:   &loggedUser.UserID,
	}

	// Create user DULU — baru assign role pakai ID yang baru ke-generate.
	// Sebelumnya GrantRole dipanggil SEBELUM Create(), padahal ID user
	// belum ada sama sekali di titik itu (auto-increment dari DB).
	if err := u.usrRepo.Create(ctx, &newUsr); err != nil {
		return err
	}

	// Subjek Casbin = user_id (string), BUKAN email — konsisten sama
	// keputusan di casbin/service.go & seluruh flow authorization.
	sub := strconv.FormatUint(newUsr.ID, 10)
	if _, err := u.authorize.GrantRole(sub, usr.Role); err != nil {
		// Compensating action: user udah kepalanjur ke-insert tapi gagal
		// dikasih role -> hapus lagi user-nya biar gak ada akun "yatim"
		// tanpa role sama sekali (gak akan pernah lolos authorization).
		_, _ = u.usrRepo.Delete(ctx, newUsr.ID)
		return fmt.Errorf("failed to grant role to new user: %w", err)
	}

	return nil
}

func (u *userService) AllUsers(ctx context.Context, req pagination.DtoPaginationRequest) (*pagination.DtoPaginationResponse, error) {
	// Clamp di sini (service/business layer), BUKAN di repository — biar
	// repository cuma eksekusi query, gak nanggung business rule kayak
	// "limit maksimal berapa".
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100 // guard biar gak ada yang iseng minta limit=999999
	}

	users, total, err := u.usrRepo.FindAllPaginated(ctx, req)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	if totalPages < 1 {
		totalPages = 1
	}

	return &pagination.DtoPaginationResponse{
		Data:       users,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalRows:  total,
		TotalPages: totalPages,
	}, nil
}

func (u *userService) Update(ctx context.Context, id uint64, usr *UpdateUserDTO) error {
	theUser, err := u.usrRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user %d not found", id)
		}
		return err // sebelumnya gak ada return di sini -> lanjut jalan dengan theUser kosong
	}

	// Ganti password: ADMIN yang edit user lain, jadi TIDAK perlu verifikasi
	// "password lama" (itu logic buat self-service change-password, beda
	// use case). Sebelumnya kode ini bandingin hash(password_baru) dengan
	// hash lama di DB -- itu SELALU mismatch karena bcrypt selalu hasilin
	// hash berbeda tiap kali di-generate (salted), walau plaintext sama.
	if usr.Password != "" {
		hashed, err := auth.HashPassword(usr.Password)
		if err != nil {
			return err
		}
		theUser.Password = hashed
	}

	if usr.Name != "" && usr.Name != theUser.Username {
		theUser.Username = usr.Name
	}

	if usr.OfficeID != 0 && usr.OfficeID != theUser.OfficeID {
		theUser.OfficeID = usr.OfficeID
	}

	// Guard nil pointer -- FotoProfile bisa nil kalau user belum pernah
	// upload foto. Sebelumnya *theUser.FotoProfile langsung di-deref
	// tanpa cek, bisa panic.
	if usr.Image != "" && (theUser.FotoProfile == nil || usr.Image != *theUser.FotoProfile) {
		theUser.FotoProfile = &usr.Image
	}

	theUser.UpdatedAt = time.Now()

	if _, err := u.usrRepo.Update(ctx, theUser.ID, theUser); err != nil {
		return err
	}

	// Role di-update TERAKHIR, setelah data user berhasil disimpan --
	// dan pakai user_id (bukan email dari DTO yang bisa aja kosong/beda).
	if usr.Role != "" {
		sub := strconv.FormatUint(theUser.ID, 10)

		// Revoke semua role lama dulu sebelum grant yang baru, supaya
		// user gak numpuk lebih dari 1 role gara-gara role diganti
		// berkali-kali (GrantRole cuma nambah, gak otomatis replace).
		oldRoles, err := u.authorize.GetRolesForUser(sub)
		if err != nil {
			return err
		}
		for _, r := range oldRoles {
			if _, err := u.authorize.RevokeRole(sub, r); err != nil {
				return err
			}
		}

		if _, err := u.authorize.GrantRole(sub, usr.Role); err != nil {
			return err
		}
	}

	return nil
}

func (u *userService) Activate(ctx context.Context, ID uint64) error {
	usr, err := u.usrRepo.FindByID(ctx, ID)
	if err != nil {
		return err
	}

	usr.IsActive = true
	usr.UpdatedAt = time.Now()

	_, err = u.usrRepo.Update(ctx, ID, usr)
	return err
}

func (u *userService) Deactivate(ctx context.Context, ID uint64) error {
	usr, err := u.usrRepo.FindByID(ctx, ID)
	if err != nil {
		return err
	}

	usr.IsActive = false
	usr.UpdatedAt = time.Now()

	_, err = u.usrRepo.Update(ctx, ID, usr)
	// CATATAN: role Casbin user ini SENGAJA gak ikut di-revoke di sini.
	// Proteksi buat user non-aktif udah di-handle di auth.RefreshToken()
	// (cek usr.IsActive sebelum terbitin token baru) -- jadi access_token
	// yang masih hidup (maks 15 menit) bakal otomatis "mati" pas expired,
	// tanpa perlu utak-atik role assignment di sini.
	return err
}

func (u *userService) Delete(ctx context.Context, ID uint64) error {
	sub := strconv.FormatUint(ID, 10)

	// Beda dari Deactivate: Delete itu permanen, jadi role di Casbin
	// juga dibersihkan sekalian (bukan sekadar nunggu token expired).
	roles, err := u.authorize.GetRolesForUser(sub)
	if err != nil {
		return err
	}
	for _, r := range roles {
		if _, err := u.authorize.RevokeRole(sub, r); err != nil {
			return err
		}
	}

	_, err = u.usrRepo.Delete(ctx, ID)
	return err
}
