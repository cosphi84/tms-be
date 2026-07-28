package auth

import (
	"tms-be/internal/conf"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID    uint64          `json:"user_id"`
	OfficeID  uint64          `json:"office_id"`
	TokenType string          `json:"token_type"`
	Role      []conf.RoleType `json:"role"`
	jwt.RegisteredClaims
}
