package sloc

import (
	"net/http"
	"strconv"
	"tms-be/internal/offices"
	"tms-be/internal/pagination"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(db *gorm.DB) *Handler {
	repo := NewRepository(db)
	officeRepo := offices.NewRepository(db)
	service := NewService(repo, officeRepo)
	return &Handler{service: service}
}

// List: GET /slocs?page=&limit=&search=&office_id=
func (h *Handler) List(c *gin.Context) {
	var req pagination.DtoPaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var officeID *uint64
	if officeIDStr := c.Query("office_id"); officeIDStr != "" {
		id, err := strconv.ParseUint(officeIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid office_id"})
			return
		}
		officeID = &id
	}

	result, err := h.service.List(c.Request.Context(), officeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Detail(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	m, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

// Options: GET /slocs/options — buat select/combobox FE.
func (h *Handler) Options(c *gin.Context) {
	options, err := h.service.Options(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": options})
}

func (h *Handler) Create(c *gin.Context) {
	var dto CreateSlocDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(c.Request.Context(), dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "sloc created"})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var dto UpdateSlocDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sloc updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sloc deleted"})
}

func parseIDParam(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}
