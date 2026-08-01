package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"tms-be/internal/conf"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Authenticate HANYA membuktikan identity (siapa kamu) lewat access_token.
// TIDAK memeriksa permission apapun — itu tugas casbin.Authorize() yang
// dipasang setelah middleware ini di route group yang sama.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is empty"})
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}
		strToken := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))

		claims, err := ParseToken(strToken)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				// Ini SINYAL buat axios interceptor FE: access_token expired,
				// waktunya panggil POST /auth/refresh (yang otomatis kirim
				// httpOnly cookie refresh_token) — user gak perlu tau ini terjadi.
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims.TokenType != "access" {
			// Endpoint yang dilindungi middleware ini gak boleh menerima
			// refresh_token — refresh_token cuma valid di POST /auth/refresh.
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
			return
		}

		// Gin context — buat handler yang butuh akses cepat tanpa GetClaims()
		c.Set("user_id", claims.UserID)

		// Request context — dibaca oleh casbin.Authorize() dan GetClaims()
		requestContext := context.WithValue(c.Request.Context(), conf.AuthContextKey, claims)
		c.Request = c.Request.WithContext(requestContext)

		c.Next()
	}
}
