package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	refreshCookieName   = "refresh_token"
	refreshCookiePath   = "/auth/refresh"  // HARUS sama persis dengan path yang di-register di RegisterRoutes()
	refreshCookieMaxAge = 7 * 24 * 60 * 60 // 7 hari, dalam detik — samain sama expiry di jwt.go
)

type Handler struct {
	service Service
}

func NewHandler(db *gorm.DB) *Handler {
	repo := NewAuthRepository(db)
	service := NewAuthenticateService(repo)
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	var dto LoginRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Authenticate(c.Request.Context(), dto, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	setRefreshCookie(c, result.RefreshToken)

	// result.RefreshToken TIDAK ikut ke sini karena json:"-" di dto.go —
	// FE cuma nerima access_token + user, refresh_token murni lewat cookie.
	c.JSON(http.StatusOK, result)
}

// Refresh dipanggil axios interceptor secara TRANSPARAN saat access_token
// expired (401) — user gak pernah tau ini terjadi. Browser otomatis kirim
// cookie refresh_token karena axios pakai withCredentials: true; endpoint
// ini TIDAK baca dari body/header, murni dari cookie.
func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshCookieName)
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found, please login again"})
		return
	}

	result, err := h.service.RefreshToken(c.Request.Context(), RefreshTokenRequestDTO{RefreshToken: refreshToken})
	if err != nil {
		// Refresh token invalid/expired -> hapus cookie, paksa FE redirect ke login
		clearRefreshCookie(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Rotate cookie dengan refresh token yang baru
	setRefreshCookie(c, result.RefreshToken)

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Logout(c *gin.Context) {
	clearRefreshCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func setRefreshCookie(c *gin.Context, token string) {
	secure := os.Getenv("APP_MODE") == "release" // HTTPS-only cookie di production
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		refreshCookieName,
		token,
		refreshCookieMaxAge,
		refreshCookiePath,
		"", // domain kosong = default ke domain request saat ini
		secure,
		true, // httpOnly = true -> JS/axios TIDAK BISA baca cookie ini, cuma browser yang kirim otomatis
	)
}

func clearRefreshCookie(c *gin.Context) {
	c.SetCookie(refreshCookieName, "", -1, refreshCookiePath, "", false, true)
}
