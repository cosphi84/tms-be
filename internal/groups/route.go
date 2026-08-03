package groups

import "github.com/gin-gonic/gin"

// public routes -- tidak butuh login.
func (h *Handler) RegisterPublicRoutes(g *gin.RouterGroup) {}

// general route, butuh login.
func (h *Handler) RegisterAuthenticatedRoutes(rg *gin.RouterGroup) {
	categories := rg.Group("/categories")
	{
		categories.GET("", h.List)               // ?page=&limit=&search=
		categories.GET("/options", h.GetOptions) // buat select FE: {value, label}
		categories.GET("/:id", h.GetDetail)
	}
}

// special auth -- butuh login + role tertentu.
func (h *Handler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	categories := rg.Group("/categories")
	{
		categories.POST("", h.Create)
		categories.PUT("/:id", h.Update)
		categories.DELETE("/:id", h.Delete)
	}
}
