package users

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"tms-be/internal/pagination"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *Model) error
	FindAllPaginated(ctx context.Context, req pagination.DtoPaginationRequest) ([]Model, int64, error)
	FindByID(ctx context.Context, id uint64) (Model, error)
	FindUser(ctx context.Context, username string) (Model, error)
	IsExists(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, id uint64, user Model) (int, error)
	Delete(ctx context.Context, id uint64) (int, error)
}

type userRepos struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &userRepos{db: db}
}

func (usr *userRepos) Create(ctx context.Context, user *Model) error {
	return gorm.G[Model](usr.db).Create(ctx, user)
}

// sortableColumns = whitelist kolom yang boleh dipakai buat ORDER BY.
// PENTING: req.SortedBy itu input dari user lewat query param — kalau
// langsung disambung mentah ke Order(), itu celah SQL injection (Gorm
// gak otomatis sanitize nama kolom di Order(), beda sama Where() yang
// pakai placeholder ?). Whitelist ini WAJIB ada.
var sortableColumns = map[string]bool{
	"id":         true,
	"username":   true,
	"email":      true,
	"created_at": true,
	"updated_at": true,
	"is_active":  true,
}

func (usr *userRepos) FindAllPaginated(ctx context.Context, req pagination.DtoPaginationRequest) ([]Model, int64, error) {
	query := usr.db.WithContext(ctx).Model(&Model{})

	if req.Search != "" {
		like := "%" + req.Search + "%"
		query = query.Where("username ILIKE ? OR email ILIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortCol := "created_at" // default aman kalau sorted_by kosong/gak valid
	if sortableColumns[req.SortedBy] {
		sortCol = req.SortedBy
	}
	sortDir := "DESC"
	if strings.EqualFold(req.SortDir, "asc") {
		sortDir = "ASC"
	}

	offset := (req.Page - 1) * req.Limit

	var users []Model
	err := query.
		Order(fmt.Sprintf("%s %s", sortCol, sortDir)).
		Limit(req.Limit).
		Offset(offset).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (usr *userRepos) FindByID(ctx context.Context, ID uint64) (Model, error) {
	return gorm.G[Model](usr.db).Where("id = ?", ID).First(ctx)
}

func (usr *userRepos) FindUser(ctx context.Context, username string) (Model, error) {
	return gorm.G[Model](usr.db).
		Where("email = ? OR username = ?", username, username).
		First(ctx)
}

// Update nulis SELURUH kolom (Select("*")) — TANPA ini, GORM skip field
// ber-zero-value pas Updates(struct). Contoh nyata bug yang ke-fix di sini:
// Deactivate() set IsActive = false, tapi false itu zero-value buat bool,
// jadi TANPA Select("*") kolom is_active gak akan pernah ke-update di SQL —
// user "dinonaktifkan" tapi tetap aktif selamanya di DB.
func (usr *userRepos) Update(ctx context.Context, id uint64, user Model) (int, error) {
	return gorm.G[Model](usr.db).Where("id = ?", id).Select("*").Updates(ctx, user)
}

func (usr *userRepos) Delete(ctx context.Context, id uint64) (int, error) {
	return gorm.G[Model](usr.db).Where("id = ?", id).Delete(ctx)
}

func (usr *userRepos) IsExists(ctx context.Context, email string) (bool, error) {
	_, err := gorm.G[Model](usr.db).
		Where("email = ?", email).
		Select("id").
		First(ctx)

	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}
