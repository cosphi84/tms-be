package auth

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tms-be/internal/helpers"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Service interface {
	Authenticate(ctx context.Context, dto LoginRequestDTO, ip string) (*LoginResponseDTO, error)
	RefreshToken(ctx context.Context, dto RefreshTokenRequestDTO) (*LoginResponseDTO, error)
}

type authServiceImpl struct {
	repo Repository
}

func NewAuthenticateService(authRepo Repository) Service {
	return &authServiceImpl{repo: authRepo}
}

func (a *authServiceImpl) Authenticate(ctx context.Context, dto LoginRequestDTO, ip string) (*LoginResponseDTO, error) {
	const maxLoginAttempts = 3
	const lockoutDuration = 1 // in hours

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
			lockedUntil := time.Now().Add(lockoutDuration * time.Hour)
			usr.LockedUntil = &lockedUntil
			usr.FailedLoginAttempt = 0
		}

		if _, err := a.repo.Update(ctx, usr.ID, usr); err != nil {
			return nil, fmt.Errorf("auth: failed to persist login attempt: %w", err)
		}
		return nil, errors.New("invalid credentials")
	}

	usr.LastLoginFrom = &ip
	usr.LastLoginAt = &now
	usr.FailedLoginAttempt = 0
	usr.LockedUntil = nil

	if _, err := a.repo.Update(ctx, usr.ID, usr); err != nil {
		return nil, fmt.Errorf("auth: failed to persist login success: %w", err)
	}

	accessToken, refreshToken, err := GenerateTokenPair(usr.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         usr,
	}, nil
}

func (a *authServiceImpl) RefreshToken(ctx context.Context, dto RefreshTokenRequestDTO) (*LoginResponseDTO, error) {
	claims, err := ParseToken(dto.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			// Refresh token juga expired -> user BENERAN harus login ulang.
			// Ini satu-satunya titik di mana "user harus tau" sesi habis.
			return nil, errors.New("session expired, please login again")
		}
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	usr, err := a.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if !usr.IsActive {
		return nil, errors.New("user is inactive")
	}

	// Rotate: refresh token lama otomatis "mati" begitu client pakai yang baru
	// (client-side, cookie lama ditimpa). Ini praktik standar refresh rotation.
	accessToken, newRefreshToken, err := GenerateTokenPair(usr.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         usr,
	}, nil
}
