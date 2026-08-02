package officeleaders

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tms-be/internal/offices"
	"tms-be/internal/pagination"
	"tms-be/internal/users"

	"gorm.io/gorm"
)

type Service interface {
	AssignLeader(ctx context.Context, dto AssignLeaderDTO) error
	List(ctx context.Context, filter ListFilter, req pagination.DtoPaginationRequest) (*pagination.DtoPaginationResponse, error)
	FindByID(ctx context.Context, id uint64) (*Model, error)
	ActiveByOffice(ctx context.Context, officeID uint64) (*Model, error)
	EndTerm(ctx context.Context, id uint64, dto EndTermDTO) error
	Update(ctx context.Context, id uint64, dto UpdateOfficeLeaderDTO) error
	Delete(ctx context.Context, id uint64) error
}

type serviceImpl struct {
	repo Repository
	// Dependensi LANGSUNG ke Repository module lain (bukan Service mereka),
	// karena di sini cuma butuh "apakah office/user ini ada", bukan business
	// rule kompleks dari module tersebut. Kalau nanti butuh validasi lebih
	// dalam (misal office harus aktif/gak di-nonaktifkan), ganti ke Service.
	officeRepo offices.Repository
	userRepo   users.Repository
}

func NewService(repo Repository, officeRepo offices.Repository, userRepo users.Repository) Service {
	return &serviceImpl{repo: repo, officeRepo: officeRepo, userRepo: userRepo}
}

func (s *serviceImpl) AssignLeader(ctx context.Context, dto AssignLeaderDTO) error {
	startDate, err := time.Parse(DateLayout, dto.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date format, expected YYYY-MM-DD: %w", err)
	}

	if _, err := s.officeRepo.FindByID(ctx, dto.OfficeID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("office %d not found", dto.OfficeID)
		}
		return err
	}

	if _, err := s.userRepo.FindByID(ctx, dto.UserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user %d not found", dto.UserID)
		}
		return err
	}

	newLeader := &Model{
		OfficeID:  dto.OfficeID,
		UserID:    dto.UserID,
		StartDate: startDate,
	}
	return s.repo.AssignLeader(ctx, newLeader)
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

func (s *serviceImpl) ActiveByOffice(ctx context.Context, officeID uint64) (*Model, error) {
	m, err := s.repo.FindActiveByOffice(ctx, officeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (s *serviceImpl) EndTerm(ctx context.Context, id uint64, dto EndTermDTO) error {
	endDate, err := time.Parse(DateLayout, dto.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end_date format, expected YYYY-MM-DD: %w", err)
	}

	rec, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if endDate.Before(rec.StartDate) {
		return errors.New("end_date cannot be before start_date")
	}

	return s.repo.EndTerm(ctx, id, endDate)
}

// Update cuma buat koreksi tanggal (typo input, dsb) -- office_id/user_id
// TIDAK bisa diubah lewat sini. Butuh ganti orang/office? Assign baru,
// itu secara konsep memang assignment yang berbeda, bukan "koreksi".
func (s *serviceImpl) Update(ctx context.Context, id uint64, dto UpdateOfficeLeaderDTO) error {
	updates := map[string]interface{}{}

	if dto.StartDate != "" {
		sd, err := time.Parse(DateLayout, dto.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start_date format: %w", err)
		}
		updates["start_date"] = sd
	}
	if dto.EndDate != "" {
		ed, err := time.Parse(DateLayout, dto.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end_date format: %w", err)
		}
		updates["end_date"] = ed
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}
	return s.repo.Update(ctx, id, updates)
}

func (s *serviceImpl) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
