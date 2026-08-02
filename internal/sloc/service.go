package sloc

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"tms-be/internal/offices"
	"tms-be/internal/pagination"

	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, dto CreateSlocDTO) error
	List(ctx context.Context, officeID *uint64, req pagination.DtoPaginationRequest) (*pagination.DtoPaginationResponse, error)
	FindByID(ctx context.Context, id uint64) (*Model, error)
	Options(ctx context.Context) ([]SlocOption, error)
	Update(ctx context.Context, id uint64, dto UpdateSlocDTO) error
	Delete(ctx context.Context, id uint64) error
}

type serviceImpl struct {
	repo       Repository
	officeRepo offices.Repository
}

func NewService(repo Repository, officeRepo offices.Repository) Service {
	return &serviceImpl{repo: repo, officeRepo: officeRepo}
}

func (s *serviceImpl) Create(ctx context.Context, dto CreateSlocDTO) error {
	if _, err := s.officeRepo.FindByID(ctx, dto.OfficeID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("office %d not found", dto.OfficeID)
		}
		return err
	}

	exists, err := s.repo.IsCodeExists(ctx, dto.Code)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("sloc with code %s already exists", dto.Code)
	}

	m := &Model{
		Code:     dto.Code,
		Name:     dto.Name,
		OfficeID: dto.OfficeID,
		Member:   dto.Member,
	}
	return s.repo.Create(ctx, m)
}

func (s *serviceImpl) List(ctx context.Context, officeID *uint64, req pagination.DtoPaginationRequest) (*pagination.DtoPaginationResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	list, total, err := s.repo.FindAllPaginated(ctx, officeID, req)
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

// Options: "value" = id, "label" = "code - name", sesuai kontrak yang diminta.
func (s *serviceImpl) Options(ctx context.Context) ([]SlocOption, error) {
	all, err := s.repo.FindAllForSelect(ctx)
	if err != nil {
		return nil, err
	}

	options := make([]SlocOption, 0, len(all))
	for _, m := range all {
		options = append(options, SlocOption{
			Value: strconv.FormatUint(m.ID, 10),
			Label: fmt.Sprintf("%s - %s", m.Code, m.Name),
		})
	}
	return options, nil
}

func (s *serviceImpl) Update(ctx context.Context, id uint64, dto UpdateSlocDTO) error {
	updates := map[string]interface{}{}

	if dto.Code != "" {
		exists, err := s.repo.IsCodeExists(ctx, dto.Code)
		if err != nil {
			return err
		}
		// Guard sederhana: kalau code ini udah dipakai row LAIN, tolak.
		// (Belum cek "punya row ini sendiri" secara eksplisit -- partial
		// unique index di DB tetap jadi backstop terakhir kalau race.)
		current, err := s.repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if exists && current.Code != dto.Code {
			return fmt.Errorf("sloc with code %s already exists", dto.Code)
		}
		updates["code"] = dto.Code
	}
	if dto.Name != "" {
		updates["name"] = dto.Name
	}
	if dto.OfficeID != 0 {
		if _, err := s.officeRepo.FindByID(ctx, dto.OfficeID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("office %d not found", dto.OfficeID)
			}
			return err
		}
		updates["office_id"] = dto.OfficeID
	}
	if dto.Member != "" {
		updates["member"] = dto.Member
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}
	return s.repo.Update(ctx, id, updates)
}

func (s *serviceImpl) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
