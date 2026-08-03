package groups

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, c *Model) error
	FindByID(ctx context.Context, id uint32) (Model, error)
	GroupExists(ctx context.Context, cat string) (bool, error)
	FindAll(ctx context.Context) ([]Model, error)
	Update(ctx context.Context, id uint32, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint32) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db: db,
	}
}

func (r *repositoryImpl) Create(ctx context.Context, c *Model) error {
	return gorm.G[Model](r.db).Create(ctx, c)
}

func (r *repositoryImpl) FindByID(ctx context.Context, id uint32) (Model, error) {
	return gorm.G[Model](r.db).Where("id = ?", id).Take(ctx)
}

func (r *repositoryImpl) GroupExists(ctx context.Context, cat string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(Model{}).
		Where("name = ?", cat).
		Count(&count).Error
	return count > 0, err
}

func (r *repositoryImpl) FindAll(ctx context.Context) ([]Model, error) {
	var all []Model
	err := r.db.WithContext(ctx).Find(&all).Error
	return all, err
}

func (r *repositoryImpl) Update(ctx context.Context, id uint32, updates map[string]interface{}) error {
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

func (r *repositoryImpl) Delete(ctx context.Context, id uint32) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&Model{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
