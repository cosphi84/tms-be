package uploader

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterProtectedRoutes(g *gin.RouterGroup) {
	uploadGroup := g.Group("/upload")
	{
		uploadGroup.POST("", h.Upload)                 // upload file baru
		uploadGroup.GET("/:uuid", h.GetMetadata)       // ambil metadata by uuid
		uploadGroup.GET("/:uuid/download", h.Download) // download file by uuid
		uploadGroup.PATCH("/:uuid", h.Patch)           // ganti/edit file by uuid
		uploadGroup.DELETE("/:uuid", h.Delete)         // soft delete file by uuid
	}
}
