package officeleaders

import (
	"net/http"
	"strconv"
	"tms-be/internal/offices"
	"tms-be/internal/pagination"
	"tms-be/internal/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(db *gorm.DB) *Handler {
	repo := NewRepository(db)
	officeRepo := offices.NewRepository(db)
	userRepo := users.NewUserRepository(db)
	service := NewService(repo, officeRepo, userRepo)
	return &Handler{service: service}
}

// List: GET /office-leaders?page=&limit=&office_id=&active_only=true
func (h *Handler) List(c *gin.Context) {
	var req pagination.DtoPaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var filter ListFilter
	if officeIDStr := c.Query("office_id"); officeIDStr != "" {
		officeID, err := strconv.ParseUint(officeIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid office_id"})
			return
		}
		filter.OfficeID = &officeID
	}
	filter.ActiveOnly = c.Query("active_only") == "true"

	result, err := h.service.List(c.Request.Context(), filter, req)
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

// ActiveByOffice: GET /office-leaders/office/:officeId/active
func (h *Handler) ActiveByOffice(c *gin.Context) {
	officeID, err := strconv.ParseUint(c.Param("officeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid office id"})
		return
	}

	m, err := h.service.ActiveByOffice(c.Request.Context(), officeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if m == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "office has no active leader"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *Handler) Assign(c *gin.Context) {
	var dto AssignLeaderDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignLeader(c.Request.Context(), dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "leader assigned"})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var dto UpdateOfficeLeaderDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) EndTerm(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var dto EndTermDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.EndTerm(c.Request.Context(), id, dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "term ended"})
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
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func parseIDParam(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}
