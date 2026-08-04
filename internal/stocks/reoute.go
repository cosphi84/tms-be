package stocks

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterPublicRoutes(g *gin.RouterGroup) {}

// Tier 2 -- cukup login.
func (h *Handler) RegisterAuthenticatedRoutes(rg *gin.RouterGroup) {
	stocks := rg.Group("/stocks")
	{
		stocks.GET("", h.List) // ?page=&limit=&tools_id=&sloc_id=&expired_only=
		stocks.GET("/:id", h.Detail)
	}
}

// Tier 3 -- butuh RBAC. Increase/decrease itu operasi transaksional,
// bukan sekadar "edit data", jadi sengaja gak ada PUT generic -- cuma
// dua aksi eksplisit ini.
func (h *Handler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	stocks := rg.Group("/stocks")
	{
		stocks.POST("/increase", h.Increase)
		stocks.POST("/decrease", h.Decrease)
		stocks.DELETE("/:id", h.Delete)
	}
}
