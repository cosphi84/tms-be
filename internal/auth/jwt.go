package auth

import (
	"context"
	"errors"
	"os"
	"time"
	"tms-be/internal/conf"

	"github.com/golang-jwt/jwt/v5"
)

func GetClaims(ctx context.Context) (*JWTClaims, error) {
	value := ctx.Value(conf.AuthContextKey)
	if value == nil {
		return nil, errors.New("no auth claims found in context")
	}
	claims, ok := value.(*JWTClaims)

	if !ok {
		return nil, errors.New("auth claims in context have invalid type")
	}
	return claims, nil
}

func GenerateTokenPair(UserID uint64, OfficeID uint64, Roles []conf.RoleType) (string, string, error) {
	secret := os.Getenv("APP_JWT_SECRET")
	if secret == "" {
		panic("APP_JWT_SECRET environment variable is not set")
	}

	accessClaims := JWTClaims{
		UserID:    UserID,
		OfficeID:  OfficeID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)), // Access token expires in 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := JWTClaims{
		UserID:    UserID,
		OfficeID:  OfficeID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // Refresh token expires in 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}
