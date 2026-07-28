package users

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *Model) error
	FindAll(ctx context.Context) ([]Model, error)
	FindByID(ctx context.Context, id uint64) (Model, error)
	FindUser(ctx context.Context, username string) (Model, error)
	IsExists(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, id uint64, user Model) (int, error)
	Delete(ctx context.Context, id uint64) (int, error)
}

type userRepos struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &userRepos{db: db}
}

func (usr *userRepos) Create(ctx context.Context, user *Model) error {
	return gorm.G[Model](usr.db).Create(ctx, user)
}

func (usr *userRepos) FindAll(ctx context.Context) ([]Model, error) {
	return gorm.G[Model](usr.db).Find(ctx)
}

func (usr *userRepos) FindByID(ctx context.Context, ID uint64) (Model, error) {
	return gorm.G[Model](usr.db).Where("id = ?", ID).First(ctx)
}

func (usr *userRepos) FindUser(ctx context.Context, username string) (Model, error) {
	return gorm.G[Model](usr.db).
		Where("email = ? OR username = ?", username, username).
		First(ctx)
}

func (usr *userRepos) Update(ctx context.Context, id uint64, user Model) (int, error) {
	return gorm.G[Model](usr.db).Where("id = ?", id).Updates(ctx, user)
}

func (usr *userRepos) Delete(ctx context.Context, id uint64) (int, error) {
	return gorm.G[Model](usr.db).Where("id = ?", id).Delete(ctx)
}

func (usr *userRepos) IsExists(ctx context.Context, email string) (bool, error) {
	_, err := gorm.G[Model](usr.db).
		Where("email = ?", email).
		Select("id").
		First(ctx)

	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}
