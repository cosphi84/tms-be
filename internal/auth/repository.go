package auth

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	FindByID(ctx context.Context, id uint64) (Model, error)
	FindUser(ctx context.Context, username string) (Model, error)
	IsExists(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, id uint64, user Model) (int, error)
}

type authRepos struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) Repository {
	return &authRepos{db: db}
}

func (usr *authRepos) FindByID(ctx context.Context, id uint64) (Model, error) {
	return gorm.G[Model](usr.db).Where("id = ?", id).Take(ctx)
}

func (usr *authRepos) FindUser(ctx context.Context, username string) (Model, error) {
	return gorm.G[Model](usr.db).
		Where("email = ? OR username = ?", username, username).
		First(ctx)
}

func (usr *authRepos) Update(ctx context.Context, id uint64, user Model) (int, error) {
	return gorm.G[Model](usr.db).Where("id = ?", id).Select("*").Updates(ctx, user)
}

func (usr *authRepos) IsExists(ctx context.Context, email string) (bool, error) {
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
