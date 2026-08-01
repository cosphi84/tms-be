package casbin

import (
	"net/http"
	"strconv"

	"tms-be/internal/auth"
	"tms-be/internal/conf"

	"github.com/gin-gonic/gin"
)

// Authorize HARUS dipasang SETELAH auth.AuthenticationMiddleware
func Authorize(authz *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		v := c.Request.Context().Value(conf.AuthContextKey)
		if v == nil {
			// Identity belum terbukti sama sekali -> 401
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		claims, ok := v.(*auth.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		obj := c.FullPath()
		act := c.Request.Method

		// SATU pemanggilan Enforce, subjek = user_id (identity).
		// Casbin resolve user_id -> role lewat g policy secara otomatis,
		sub := strconv.FormatUint(claims.UserID, 10)

		allowed, err := authz.Enforce(sub, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !allowed {
			// Identity valid, tapi gak punya izin -> 403
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
