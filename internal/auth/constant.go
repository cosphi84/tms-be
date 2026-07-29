package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID    uint64 `json:"user_id"`
	OfficeID  uint64 `json:"office_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}
