package auth

import "github.com/gin-gonic/gin"

// RegisterPublicRoutes: endpoint yang WAJIB bisa diakses TANPA access_token valid.
func (h *Handler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh", h.Refresh)
	}
}

// RegisterProtectedRoutes: butuh identity (access_token valid), TAPI TIDAK
// butuh authorization/RBAC — semua role yang sedang login boleh logout,
// gak ada permission khusus buat aksi ini.
func (h *Handler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/logout", h.Logout)
	}
}
