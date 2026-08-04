package stocks

import (
	"context"
	"errors"
	"fmt"
	"time"
	mtools "tms-be/internal/master-tools"
	"tms-be/internal/pagination"
	sloc "tms-be/internal/storage-locations"

	"gorm.io/gorm"
)

type Service interface {
	IncreaseStock(ctx context.Context, dto IncreaseStockDTO) error
	DecreaseStock(ctx context.Context, dto DecreaseStockDTO) error
	List(ctx context.Context, filter ListFilter, req pagination.DtoPaginationRequest) (*pagination.DtoPaginationResponse, error)
	FindByID(ctx context.Context, id uint64) (*Model, error)
	Delete(ctx context.Context, id uint64) error
}

type serviceImpl struct {
	repo     Repository
	toolRepo mtools.Repository
	slocRepo sloc.Repository
}

func NewService(repo Repository, toolRepo mtools.Repository, slocRepo sloc.Repository) Service {
	return &serviceImpl{repo: repo, toolRepo: toolRepo, slocRepo: slocRepo}
}

func (s *serviceImpl) IncreaseStock(ctx context.Context, dto IncreaseStockDTO) error {
	tool, err := s.toolRepo.FindByID(ctx, dto.ToolsID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("tool %d not found", dto.ToolsID)
		}
		return err
	}

	if _, err := s.slocRepo.FindByID(ctx, dto.SlocID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("storage location %d not found", dto.SlocID)
		}
		return err
	}

	// expired_date dihitung dari usage_period tool ini, DIHITUNG ULANG
	// setiap kali barang baru masuk (bukan diakumulasi dari expired_date lama).
	expiredDate, err := tool.CalculateExpiry(time.Now())
	if err != nil {
		return fmt.Errorf("failed to calculate expiry: %w", err)
	}

	var remarks *string
	if dto.Remarks != "" {
		remarks = &dto.Remarks
	}

	_, err = s.repo.IncreaseStock(ctx, dto.ToolsID, dto.SlocID, dto.Qty, expiredDate, remarks)
	return err
}

func (s *serviceImpl) DecreaseStock(ctx context.Context, dto DecreaseStockDTO) error {
	return s.repo.DecreaseStock(ctx, dto.ToolsID, dto.SlocID, dto.Qty)
}

func (s *serviceImpl) List(ctx context.Context, filter ListFilter, req pagination.DtoPaginationRequest) (*pagination.DtoPaginationResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	list, total, err := s.repo.FindAllPaginated(ctx, filter, req)
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

func (s *serviceImpl) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
