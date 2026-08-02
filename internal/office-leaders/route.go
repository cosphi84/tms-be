package officeleaders

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterPublicRoutes(g *gin.RouterGroup) {}

// Tier 2 -- baca data, cukup login, gak butuh RBAC.
func (h *Handler) RegisterAuthenticatedRoutes(rg *gin.RouterGroup) {
	leaders := rg.Group("/office-leaders")
	{
		leaders.GET("", h.List) // ?page=&limit=&office_id=&active_only=
		leaders.GET("/:id", h.Detail)
		leaders.GET("/office/:officeId/active", h.ActiveByOffice)
	}
}

// Tier 3 -- create/update/end-term/delete, WAJIB lolos RBAC.
func (h *Handler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	leaders := rg.Group("/office-leaders")
	{
		leaders.POST("", h.Assign)
		leaders.PUT("/:id", h.Update)
		leaders.PATCH("/:id/end-term", h.EndTerm)
		leaders.DELETE("/:id", h.Delete)
	}
}
