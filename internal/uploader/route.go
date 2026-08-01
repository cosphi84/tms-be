package uploader

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterProtectedRoutes(g *gin.RouterGroup) {
	uploadGroup := g.Group("/upload")
	{
		uploadGroup.POST("")
	}
}
