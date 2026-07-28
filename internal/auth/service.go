package auth

import (
	"context"
	"errors"
	"time"
	"tms-be/internal/casbin"
	"tms-be/internal/conf"
	"tms-be/internal/helpers"

	"gorm.io/gorm"
)

type Service interface {
	Authenticate(ctx context.Context, dto LoginRequestDTO, ip string) (*LoginResponseSTO, error)
	RefreshToken(ctx context.Context, dto RefreshTokenRequestDTO) (*LoginResponseSTO, error)
}
type authServiceImpl struct {
	repo Repository
	role *casbin.Service
}

func NewAuthenticateService(authRepo Repository) Service {
	return &authServiceImpl{
		repo: authRepo,
	}
}
func (a *authServiceImpl) Authenticate(ctx context.Context, dto LoginRequestDTO, ip string) (*LoginResponseSTO, error) {
	maxLoginAttempts := 3
	lockoutDuration := 1 // in hours

	usr, err := a.repo.FindUser(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email atau password salah")
		}
		return nil, err
	}

	if !usr.IsActive {
		return nil, errors.New("invalid credentials")
	}

	now := time.Now()
	if usr.LockedUntil != nil {
		if now.Before(*usr.LockedUntil) {
			return nil, errors.New("account locked")
		}

		usr.LockedUntil = nil
		usr.FailedLoginAttempt = 0
	}

	valid := helpers.VerifyPassword(dto.Password, usr.Password)
	if !valid {
		nTry := usr.FailedLoginAttempt + 1

		usr.FailedLoginAttempt = nTry
		if nTry >= maxLoginAttempts {
			lockedUntil := time.Now().Add(time.Duration(lockoutDuration) * time.Hour)
			usr.LockedUntil = &lockedUntil
			usr.FailedLoginAttempt = 0
		}

		err = a.repo.Udate(ctx, usr)

		if err != nil {
			return nil, err
		}
		return nil, errors.New("invalid credentials")
	}

	usr.LastLoginFrom = &ip
	usr.LastLoginAt = &now
	usr.FailedLoginAttempt = 0
	usr.LockedUntil = nil

	_ = a.repo.Update(usr)

	usrRolesRaw, _ := a.role.GetRoleForUser(usr.Email)
	usrRoles := make([]conf.RoleType, len(usrRolesRaw))
	for i, r := range usrRolesRaw {
		usrRoles[i] = conf.RoleType(r)
	}
	accessToken, refreshToken, err := GenerateTokenPair(usr.ID, usr.OfficeID, usrRoles)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         usr,
	}, nil

}

func (a *authServiceImpl) RefreshToken(ctx context.Context, dto RefreshTokenRequestDTO) (*LoginResponseSTO, error) {
}
