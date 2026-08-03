package mtools

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	tools := rg.Group("/tools")
	{
		tools.POST("", h.RegisterTool)
		tools.PUT("/:id", h.Update)
		tools.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) RegisterAuthenticatedRoutes(rg *gin.RouterGroup) {
	tools := rg.Group("/tools")
	{
		tools.GET("", h.List)
		tools.GET("/:id", h.FindByID)
	}
}
