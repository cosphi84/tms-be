package mtools

import (
	"context"
	"errors"
	"fmt"
	"tms-be/internal/category"
	"tms-be/internal/groups"
	"tms-be/internal/pagination"
	"tms-be/internal/uploader"

	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, dto CreateToolDTO) error
	List(ctx context.Context, req *pagination.DtoPaginationRequest) (*pagination.DtoPaginationResponse, error)
	FindByID(ctx context.Context, id uint64) (*Model, error)
	Update(ctx context.Context, id uint64, dto UpdateToolDTO) error
	Delete(ctx context.Context, id uint64) error
}

type serviceImpl struct {
	repo         Repository
	categoryRepo category.Repository
	groupRepo    groups.Repository
	fileRepo     uploader.Repository
}

func NewService(repo Repository, categoryRepo category.Repository, groupRepo groups.Repository, fileRepo uploader.Repository) Service {
	return &serviceImpl{
		repo:         repo,
		categoryRepo: categoryRepo,
		groupRepo:    groupRepo,
		fileRepo:     fileRepo,
	}
}

func (s *serviceImpl) Create(ctx context.Context, dto CreateToolDTO) error {
	if _, err := s.categoryRepo.FindByID(ctx, dto.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("category %d not found", dto.CategoryID)
		}
		return err
	}

	if _, err := s.groupRepo.FindByID(ctx, dto.GroupID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("group %d not found", dto.GroupID)
		}
		return err
	}

	if dto.PhotoID != 0 {
		if _, err := s.fileRepo.FindByID(ctx, dto.PhotoID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("photo %d not found", dto.PhotoID)
			}
			return err
		}
	}

	exists, err := s.repo.ToolsIsExists(ctx, dto.Code, dto.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("tool with code %s already exists", dto.Code)
	}

	m := &Model{
		Code:        dto.Code,
		Name:        dto.Name,
		Brand:       dto.Brand,
		Type:        dto.Type,
		SerialNum:   dto.SerialNum,
		CategoryID:  dto.CategoryID,
		GroupID:     dto.GroupID,
		Price:       dto.Price,
		UsagePeriod: dto.UsagePeriod,
		IsActive:    true,
	}
	if dto.PhotoID != 0 {
		m.PhotoID = &dto.PhotoID
	}

	return s.repo.Create(ctx, m)
}

func (s *serviceImpl) List(ctx context.Context, req *pagination.DtoPaginationRequest) (*pagination.DtoPaginationResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	list, total, err := s.repo.FindAllPaginated(ctx, req)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	if totalPages < 1 {
		totalPages = 1
	}

	return &pagination.DtoPaginationResponse{
		Data:       list,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalRows:  total,
		TotalPages: totalPages,
	}, nil
}

func (s *serviceImpl) FindByID(ctx context.Context, id uint64) (*Model, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *serviceImpl) Update(ctx context.Context, id uint64, dto UpdateToolDTO) error {
	updates := map[string]interface{}{}

	if dto.Code != "" && dto.Name != "" {
		exists, err := s.repo.ToolsIsExists(ctx, dto.Code, dto.Name)
		if err != nil {
			return err
		}
		current, err := s.repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if exists && current.Code != dto.Code {
			return fmt.Errorf("tool with code %s and name %s already exists", dto.Code, dto.Name)
		}
		updates["code"] = dto.Code
	}
	if dto.Name != "" {
		updates["name"] = dto.Name
	}
	if dto.Brand != "" {
		updates["brand"] = dto.Brand
	}
	if dto.Type != "" {
		updates["type"] = dto.Type
	}
	if dto.SerialNum != "" {
		updates["serial_num"] = dto.SerialNum
	}

	if dto.CategoryID != 0 {
		if _, err := s.categoryRepo.FindByID(ctx, dto.CategoryID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("category %d not found", dto.CategoryID)
			}
			return err
		}
		updates["category_id"] = dto.CategoryID
	}
	if dto.GroupID != 0 {
		updates["group_id"] = dto.GroupID
	}
	if dto.PhotoID != 0 {
		updates["photo_id"] = dto.PhotoID
	}
	if dto.Price != 0 {
		updates["price"] = dto.Price
	}
	if dto.UsagePeriod != "" {
		updates["usage_period"] = dto.UsagePeriod
	}
	if dto.IsActive != nil {
		updates["is_active"] = *dto.IsActive
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}
	return s.repo.Update(ctx, id, updates)
}

func (s *serviceImpl) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
