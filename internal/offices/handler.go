package offices

import (
	"net/http"
	"strconv"
	"tms-be/internal/pagination"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(db *gorm.DB) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return &Handler{service: service}
}

var validLevels = map[string]bool{"cabang": true, "sdss": true, "ssr": true}

// List: GET /offices?level=cabang|sdss|ssr — HANYA 3 level ini, sesuai
// requirement ("get semua Office level cabang, sdss, atau ssr saja").
// Kalau butuh semua level (termasuk hq/sass/tc), pakai GET /offices/tree.
func (h *Handler) List(c *gin.Context) {
	level := c.Query("level")
	if !validLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "level harus salah satu dari: cabang, sdss, ssr"})
		return
	}

	var req pagination.DtoPaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.FindAllOffice(c.Request.Context(), level, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Children: GET /offices/:id/children — anak LANGSUNG dari office tsb.
func (h *Handler) Children(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid office id"})
		return
	}

	children, err := h.service.FindChildren(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": children})
}

func (h *Handler) Create(c *gin.Context) {
	var dto OfficeRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateOffice(c.Request.Context(), dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "office created"})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid office id"})
		return
	}

	var dto OfficeRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "office updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid office id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "office deleted"})
}

// Tree: GET /offices/tree — SEMUA level, nested, sorted by code. Buat
// select/combobox shadcn di FE.
func (h *Handler) Tree(c *gin.Context) {
	tree, err := h.service.GetOfficeTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tree})
}

func parseIDParam(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}
