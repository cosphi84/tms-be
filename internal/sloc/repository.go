package sloc

import (
	"context"
	"fmt"
	"strings"
	"tms-be/internal/pagination"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, m *Model) error
	FindByID(ctx context.Context, id uint64) (Model, error)
	FindAllPaginated(ctx context.Context, officeID *uint64, req pagination.DtoPaginationRequest) ([]Model, int64, error)
	FindAllForSelect(ctx context.Context) ([]Model, error)
	IsCodeExists(ctx context.Context, code string) (bool, error)
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint64) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) Create(ctx context.Context, m *Model) error {
	return gorm.G[Model](r.db).Create(ctx, m)
}

func (r *repositoryImpl) FindByID(ctx context.Context, id uint64) (Model, error) {
	return gorm.G[Model](r.db).Where("id = ?", id).Take(ctx)
}

var slocSortableColumns = map[string]bool{
	"code":       true,
	"name":       true,
	"member":     true,
	"created_at": true,
}

func (r *repositoryImpl) FindAllPaginated(ctx context.Context, officeID *uint64, req pagination.DtoPaginationRequest) ([]Model, int64, error) {
	query := r.db.WithContext(ctx).Model(&Model{})

	if officeID != nil {
		query = query.Where("office_id = ?", *officeID)
	}
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
		Preload("Office").
		Order(fmt.Sprintf("%s %s", sortCol, sortDir)).
		Limit(req.Limit).
		Offset(offset).
		Find(&result).Error

	return result, total, err
}

// FindAllForSelect: SEMUA sloc aktif, diurutkan by code, tanpa pagination --
// dipakai khusus buat dropdown/select FE, bukan buat tabel data.
func (r *repositoryImpl) FindAllForSelect(ctx context.Context) ([]Model, error) {
	var all []Model
	err := r.db.WithContext(ctx).Order("code ASC").Find(&all).Error
	return all, err
}

func (r *repositoryImpl) IsCodeExists(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Model{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

func (r *repositoryImpl) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Model{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
