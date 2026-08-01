package users

import "github.com/gin-gonic/gin"

// RegisterRoutes — semua endpoint di sini WAJIB Tier 3 (authMW + casbinMW),
// dipasang dari internal/app/route.go lewat protected group. Gak ada
// endpoint publik di module ini.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("", h.List)
		users.POST("", h.Create)
		users.PUT("/:id", h.Update)
		users.PATCH("/:id/activate", h.Activate)
		users.PATCH("/:id/deactivate", h.Deactivate)
		users.DELETE("/:id", h.Delete)
	}
}
