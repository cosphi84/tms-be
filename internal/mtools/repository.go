package mtools

import (
	"context"
	"fmt"
	"strings"
	"tms-be/internal/pagination"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, model *Model) error
	FindAllPaginated(ctx context.Context, req *pagination.DtoPaginationRequest) ([]Model, int64, error)
	FindByID(ctx context.Context, id uint64) (Model, error)
	ToolsIsExists(ctx context.Context, toolID string) (bool, error)
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint64) error
}

type reposImpl struct {
	db *gorm.DB
}

func NewMToolsRepository(db *gorm.DB) Repository {
	return &reposImpl{db: db}
}

// Create implements [Repository].
func (r *reposImpl) Create(ctx context.Context, model *Model) error {
	return gorm.G[Model](r.db).Create(ctx, model)
}

// Delete implements [Repository].
func (r *reposImpl) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Model{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

var slocSortableColumns = map[string]bool{
	"code":       true,
	"name":       true,
	"brand":      true,
	"type":       true,
	"category":   true,
	"group":      true,
	"created_at": true,
}

// FindAllPaginated implements [Repository].
func (r *reposImpl) FindAllPaginated(ctx context.Context, req *pagination.DtoPaginationRequest) ([]Model, int64, error) {
	query := r.db.WithContext(ctx).Model(&Model{})

	if req.Search != "" {
		like := "%" + req.Search + "%"
		query = query.Where("code ILIKE ? OR name ILIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortCol := "code"
	if slocSortableColumns[req.SortedBy] {
		sortCol = req.SortedBy
	}
	sortDir := "ASC"
	if strings.EqualFold(req.SortDir, "desc") {
		sortDir = "DESC"
	}
	offset := (req.Page - 1) * req.Limit

	var result []Model
	err := query.
		Preload("Photo").
		Preload("Category").
		Preload("Group").
		Order(fmt.Sprintf("%s %s", sortCol, sortDir)).
		Limit(req.Limit).
		Offset(offset).
		Find(&result).Error

	return result, total, err
}

// FindByID implements [Repository].
func (r *reposImpl) FindByID(ctx context.Context, id uint64) (Model, error) {
	var model Model
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		return Model{}, result.Error
	}
	return model, nil
}

// ToolsIsExists implements [Repository].
func (r *reposImpl) ToolsIsExists(ctx context.Context, toolID string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&Model{}).Where("name = ?", toolID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// Update implements [Repository].
func (r *reposImpl) Update(ctx context.Context, id uint64, updates map[string]interface{}) (err error) {
	result := r.db.WithContext(ctx).
		Model(&Model{}).
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
