package uploader

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
	db      *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	repo := NewUploaderRepository(db)
	storage, _ := NewLocalStorage()
	validator := NewFileValidator()
	service := NewUploadService(repo, storage, validator)
	return &Handler{
		service: service,
		db:      db,
	}
}

// Upload menangani POST /upload
// Body: multipart/form-data { folder, access, file }
func (h *Handler) Upload(c *gin.Context) {
	var dto UploadRequestDTO
	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file: field 'file' is required"})
		return
	}

	resp, err := h.service.UploadFile(c.Request.Context(), dto, fileHeader)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMetadata menangani GET /upload/:uuid
func (h *Handler) GetMetadata(c *gin.Context) {
	fileUUID := c.Param("uuid")
	if fileUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid: parameter is required"})
		return
	}

	resp, err := h.service.GetMetadata(c.Request.Context(), fileUUID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Download menangani GET /upload/:uuid/download
// Streaming langsung ke response writer, bukan load full file ke memory.
func (h *Handler) Download(c *gin.Context) {
	fileUUID := c.Param("uuid")
	if fileUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid: parameter is required"})
		return
	}

	reader, meta, err := h.service.OpenForDownload(c.Request.Context(), fileUUID)
	if err != nil {
		h.respondError(c, err)
		return
	}
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, meta.OriginalName))
	c.Header("Content-Type", meta.MimeType)
	c.Header("Content-Length", fmt.Sprintf("%d", meta.Size))

	if _, err := io.Copy(c.Writer, reader); err != nil {
		// Header (dan mungkin sebagian body) sudah terkirim ke client,
		// jadi tidak bisa lagi balas dengan JSON error di titik ini.
		// Cukup dicatat lewat gin error collector untuk logger/middleware.
		_ = c.Error(fmt.Errorf("file: failed to stream file body: %w", err))
	}
}

// Patch menangani PATCH /upload/:uuid — mengganti file lama dengan file baru,
// mengembalikan metadata file baru (UploadResponseDTO).
// Body: multipart/form-data { folder, access, file }
func (h *Handler) Patch(c *gin.Context) {
	oldUUID := c.Param("uuid")
	if oldUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid: parameter is required"})
		return
	}

	var dto ReplaceRequestDTO
	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file: field 'file' is required"})
		return
	}

	tx := h.db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "file: failed to start transaction"})
		return
	}

	resp, err := h.service.ReplaceWithTx(c.Request.Context(), tx, dto.Folder, fileHeader, oldUUID)
	if err != nil {
		tx.Rollback()
		h.respondError(c, err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "file: failed to commit transaction"})
		return
	}

	// File baru sudah aman ter-commit. Hapus fisik file lama dilakukan
	// setelah commit — kalau gagal, itu bukan alasan untuk fail-kan
	// request (metadata baru sudah valid), cukup dilog.
	if err := h.service.FinalizeReplace(c.Request.Context(), oldUUID); err != nil {
		_ = c.Error(fmt.Errorf("file: failed to finalize replace for old file %s: %w", oldUUID, err))
	}

	c.JSON(http.StatusOK, resp)
}

// Delete menangani DELETE /upload/:uuid — soft delete metadata file.
// Penghapusan fisik dilakukan belakangan oleh retention job (lihat service.DeleteFile).
func (h *Handler) Delete(c *gin.Context) {
	fileUUID := c.Param("uuid")
	if fileUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid: parameter is required"})
		return
	}

	if err := h.service.DeleteFile(c.Request.Context(), fileUUID); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}

func (h *Handler) respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrFileMetaNotFound), errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
