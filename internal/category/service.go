package category

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

type Service interface {
	Create(ctx context.Context, dto CreateCategoryDTO) error
	List(ctx context.Context) ([]Model, error)
	FindByID(ctx context.Context, id uint32) (*Model, error)
	Options(ctx context.Context) ([]CategoryOption, error)
	Update(ctx context.Context, id uint32, dto UpdateCategoryDTO) error
	Delete(ctx context.Context, id uint32) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (s *serviceImpl) Create(ctx context.Context, dto CreateCategoryDTO) error {
	exist, err := s.repo.CategoryExists(ctx, dto.Name)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("category %s already exists", dto.Name)
	}

	c := &Model{
		Name:        dto.Name,
		Description: dto.Description,
	}

	return s.repo.Create(ctx, c)
}

func (s *serviceImpl) List(ctx context.Context) ([]Model, error) {
	return s.repo.FindAll(ctx)
}

func (s *serviceImpl) FindByID(ctx context.Context, id uint32) (*Model, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *serviceImpl) Options(ctx context.Context) ([]CategoryOption, error) {
	all, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	options := make([]CategoryOption, 0, len(all))
	for _, option := range all {
		options = append(options, CategoryOption{
			Value: strconv.FormatInt(int64(option.ID), 10),
			Label: option.Name,
		})
	}

	return options, nil
}

func (s *serviceImpl) Update(ctx context.Context, id uint32, dto UpdateCategoryDTO) error {
	updates := map[string]interface{}{}

	if dto.Name != "" {
		exist, err := s.repo.CategoryExists(ctx, dto.Name)
		if err != nil {
			return err
		}
		rec, err := s.repo.FindByID(ctx, id)
		if err != nil {
			return err
		}

		if exist && rec.Name != dto.Name {
			return fmt.Errorf("category %s already exists", dto.Name)
		}
		updates["name"] = dto.Name
	}
	if dto.Description != "" {
		updates["description"] = dto.Description
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}
	return s.repo.Update(ctx, id, updates)
}

func (s *serviceImpl) Delete(ctx context.Context, id uint32) error {
	return s.repo.Delete(ctx, id)
}
