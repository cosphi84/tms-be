package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RoleType string
type Key string

const (
	RoleSuperadmin RoleType = "superadmin"
	RoleAdminHQ    RoleType = "admin_hq"
	RoleSvcHead    RoleType = "service_head"
	RoleManagement RoleType = "management"
	RoleAuditor    RoleType = "auditor"
	RoleAuditorCS  RoleType = "auditor_cs"
	RoleTechnician RoleType = "technician"
)

const (
	AuthContextKey Key = "auth"
)

type JWTClaims struct {
	UserID    int64      `json:"user_id"`
	OfficeID  int64      `json:"office_id"`
	Role      []RoleType `json:"role"`
	TokenType string     `json:"token_type"`
	jwt.RegisteredClaims
}

func GetClaims(ctx context.Context) (*JWTClaims, error) {
	value := ctx.Value(AuthContextKey)
	if value == nil {
		return nil, errors.New("no auth claims found in context")
	}
	claims, ok := value.(*JWTClaims)

	if !ok {
		return nil, errors.New("auth claims in context have invalid type")
	}
	return claims, nil
}

func GenerateTokenPair(UserID int64, OfficeID int64, Role []RoleType) (string, string, error) {
	secret := os.Getenv("APP_JWT_SECRET")
	if secret == "" {
		panic("APP_JWT_SECRET environment variable is not set")
	}

	accessClaims := JWTClaims{
		UserID:    UserID,
		OfficeID:  OfficeID,
		Role:      Role,
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
		Role:      Role,
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
