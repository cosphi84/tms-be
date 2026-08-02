package sloc

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterPublicRoutes(g *gin.RouterGroup) {}

// Tier 2 -- cukup login.
func (h *Handler) RegisterAuthenticatedRoutes(rg *gin.RouterGroup) {
	slocs := rg.Group("/slocs")
	{
		slocs.GET("", h.List)            // ?page=&limit=&search=&office_id=
		slocs.GET("/options", h.Options) // buat select FE: {value, label}
		slocs.GET("/:id", h.Detail)
	}
}

// Tier 3 -- butuh RBAC.
func (h *Handler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	slocs := rg.Group("/slocs")
	{
		slocs.POST("", h.Create)
		slocs.PUT("/:id", h.Update)
		slocs.DELETE("/:id", h.Delete)
	}
}
