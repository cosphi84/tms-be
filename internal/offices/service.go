package offices

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"tms-be/internal/pagination"

	"gorm.io/gorm"
)

type Service interface {
	CreateHQ(ctx context.Context, dto HQRequestDTO) (*Model, error)
	CreateOffice(ctx context.Context, dto OfficeRequestDTO) error
	FindAllOffice(ctx context.Context, level string, req pagination.DtoPaginationRequest) (*pagination.DtoPaginationResponse, error)
	FindChildren(ctx context.Context, parentID uint64) ([]Model, error)
	FindOffice(ctx context.Context, code string) (*Model, error)
	FindByID(ctx context.Context, id uint64) (*Model, error)
	GetOfficeTree(ctx context.Context) ([]*OfficeTreeOption, error)
	Update(ctx context.Context, id uint64, dto OfficeRequestDTO) error
	Delete(ctx context.Context, id uint64) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

// CreateHQ idempotent -- kalau kode HQ udah ada, return yang lama (bukan
// error), biar seeder aman dijalankan berkali-kali tanpa perlu dicek
// manual dulu dari caller.
func (s *serviceImpl) CreateHQ(ctx context.Context, dto HQRequestDTO) (*Model, error) {
	existing, err := s.repo.FindOffice(ctx, dto.Code)
	if err == nil {
		return &existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	office := &Model{
		Code: dto.Code,
		Name: dto.Name,
		Type: dto.Type,
	}
	if err := s.repo.CreateHQ(ctx, office); err != nil {
		return nil, err
	}
	return office, nil
}

func (s *serviceImpl) CreateOffice(ctx context.Context, dto OfficeRequestDTO) error {
	if dto.ParentID == nil {
		return errors.New("parent_id is required")
	}

	parent, err := s.repo.FindByID(ctx, *dto.ParentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("parent office %d not found", *dto.ParentID)
		}
		return err
	}

	office := &Model{
		Code: dto.Code,
		Name: dto.Name,
		Type: dto.Type,
	}
	return s.repo.CreateBranch(ctx, office, parent.ID)
}

func (s *serviceImpl) FindAllOffice(ctx context.Context, level string, req pagination.DtoPaginationRequest) (*pagination.DtoPaginationResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	list, total, err := s.repo.FindAllByLevel(ctx, level, req)
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

func (s *serviceImpl) FindChildren(ctx context.Context, parentID uint64) ([]Model, error) {
	return s.repo.FindChildren(ctx, parentID)
}

func (s *serviceImpl) FindOffice(ctx context.Context, code string) (*Model, error) {
	rec, err := s.repo.FindOffice(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

func (s *serviceImpl) FindByID(ctx context.Context, id uint64) (*Model, error) {
	rec, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

// GetOfficeTree bangun struktur nested (buat shadcn select) dari flat list,
// di-grouping by ParentID (O(n)), tiap level diurutkan by Code (O(n log n)
// total). Root = office yang ParentID-nya nil (HQ).
func (s *serviceImpl) GetOfficeTree(ctx context.Context) ([]*OfficeTreeOption, error) {
	flat, err := s.repo.FindAllFlat(ctx)
	if err != nil {
		return nil, err
	}

	childrenByParent := make(map[uint64][]Model)
	var roots []Model
	for _, o := range flat {
		if o.ParentID == nil {
			roots = append(roots, o)
		} else {
			childrenByParent[*o.ParentID] = append(childrenByParent[*o.ParentID], o)
		}
	}

	var assemble func(nodes []Model) []*OfficeTreeOption
	assemble = func(nodes []Model) []*OfficeTreeOption {
		sort.Slice(nodes, func(i, j int) bool { return nodes[i].Code < nodes[j].Code })

		result := make([]*OfficeTreeOption, 0, len(nodes))
		for _, n := range nodes {
			result = append(result, &OfficeTreeOption{
				Value:    strconv.FormatUint(n.ID, 10),
				Label:    n.Name,
				Code:     n.Code,
				Type:     n.Type,
				Children: assemble(childrenByParent[n.ID]),
			})
		}
		return result
	}

	return assemble(roots), nil
}

// Update cuma bisa ubah code/name/type -- parent_id TIDAK bisa diubah lewat
// endpoint ini (lihat catatan di repository.go).
func (s *serviceImpl) Update(ctx context.Context, id uint64, dto OfficeRequestDTO) error {
	updates := map[string]interface{}{}
	if dto.Code != "" {
		updates["code"] = dto.Code
	}
	if dto.Name != "" {
		updates["name"] = dto.Name
	}
	if dto.Type != "" {
		updates["type"] = dto.Type
	}
	if len(updates) == 0 {
		return errors.New("no fields to update")
	}
	return s.repo.Update(ctx, id, updates)
}

func (s *serviceImpl) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
