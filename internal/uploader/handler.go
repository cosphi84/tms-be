package uploader

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(db *gorm.DB) *Handler {
	repo := NewUploaderRepository(db)
	storage, _ := NewLocalStorage()
	validator := NewFileValidator()
	service := NewUploadService(repo, storage, validator)
	return &Handler{
		service: service,
	}
}

func (h *Handler) Upload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
