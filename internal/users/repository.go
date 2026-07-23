package users

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *Model) error
	FindAll(ctx context.Context) ([]*Model, error)
	FindByID(ctx context.Context, id int32) (*Model, error)
	FindByEmail(ctx context.Context, email string) (*Model, error)
	Update(ctx context.Context, user *Model) error
	Delete(ctx context.Context, id int32) error
}

type userRepos struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (usr *userRepos) Create(ctx context.Context, user *Model) error {
	err := gorm.G[Model](usr.db).Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
