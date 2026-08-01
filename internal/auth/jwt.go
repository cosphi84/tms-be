package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"tms-be/internal/conf"

	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret = satu-satunya tempat baca APP_JWT_SECRET dari env.
// Dipakai oleh ParseToken, GenerateTokenPair, dan middleware.go
func getJWTSecret() []byte {
	secret := os.Getenv("APP_JWT_SECRET")
	if secret == "" {
		panic("APP_JWT_SECRET environment variable is not set")
	}
	return []byte(secret)
}

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

// ParseToken parse & validasi JWT string apapun (access ATAU refresh —
// pengecekan TokenType dilakukan oleh caller, bukan di sini, supaya
// fungsi ini reusable buat keduanya).
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return getJWTSecret(), nil
		},
	)
	if err != nil {
		return nil, err // biarkan error asli (termasuk jwt.ErrTokenExpired) tembus ke caller
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

// GenerateTokenPair menerbitkan access + refresh token.
func GenerateTokenPair(userID uint64) (accessToken string, refreshToken string, err error) {
	secret := getJWTSecret()

	accessClaims := JWTClaims{
		UserID:    userID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := JWTClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
