package mtools

import (
	"strconv"
	"tms-be/internal/casbin"
	"tms-be/internal/category"
	"tms-be/internal/groups"
	"tms-be/internal/pagination"
	"tms-be/internal/uploader"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(db *gorm.DB, casbinSvc *casbin.Service) *Handler {
	repo := NewRepository(db)
	categoryRepo := category.NewRepository(db)
	groupRepo := groups.NewRepository(db)
	fileRepo := uploader.NewUploaderRepository(db)

	return &Handler{
		service: NewService(repo, categoryRepo, groupRepo, fileRepo),
	}
}

func (h *Handler) RegisterTool(c *gin.Context) {
	dto := CreateToolDTO{}
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	res := h.service.Create(c.Request.Context(), dto)
	if res != nil {
		c.JSON(400, gin.H{"error": res.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "tool created successfully"})
}

func (h *Handler) List(c *gin.Context) {
	req := &pagination.DtoPaginationRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) FindByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	dto := UpdateToolDTO{}
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = h.service.Update(c.Request.Context(), id, dto)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "tool updated successfully"})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "tool deleted successfully"})
}

func parseIDParam(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}
