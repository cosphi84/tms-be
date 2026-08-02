package officeleaders

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tms-be/internal/pagination"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ListFilter struct {
	OfficeID   *uint64
	ActiveOnly bool
}

type Repository interface {
	FindByID(ctx context.Context, id uint64) (Model, error)
	FindActiveByOffice(ctx context.Context, officeID uint64) (Model, error)
	FindActiveByUser(ctx context.Context, userID uint64) (Model, error)
	FindAllPaginated(ctx context.Context, filter ListFilter, req pagination.DtoPaginationRequest) ([]Model, int64, error)
	AssignLeader(ctx context.Context, newLeader *Model) error
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	EndTerm(ctx context.Context, id uint64, endDate time.Time) error
	Delete(ctx context.Context, id uint64) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) FindByID(ctx context.Context, id uint64) (Model, error) {
	return gorm.G[Model](r.db).Where("id = ?", id).Take(ctx)
}

func (r *repositoryImpl) FindActiveByOffice(ctx context.Context, officeID uint64) (Model, error) {
	return gorm.G[Model](r.db).
		Where("office_id = ? AND end_date IS NULL", officeID).
		Take(ctx)
}

func (r *repositoryImpl) FindActiveByUser(ctx context.Context, userID uint64) (Model, error) {
	return gorm.G[Model](r.db).
		Where("user_id = ? AND end_date IS NULL", userID).
		Take(ctx)
}

var leaderSortableColumns = map[string]bool{
	"start_date": true,
	"end_date":   true,
	"created_at": true,
}

func (r *repositoryImpl) FindAllPaginated(ctx context.Context, filter ListFilter, req pagination.DtoPaginationRequest) ([]Model, int64, error) {
	query := r.db.WithContext(ctx).Model(&Model{})

	if filter.OfficeID != nil {
		query = query.Where("office_id = ?", *filter.OfficeID)
	}
	if filter.ActiveOnly {
		query = query.Where("end_date IS NULL")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortCol := "start_date"
	if leaderSortableColumns[req.SortedBy] {
		sortCol = req.SortedBy
	}
	sortDir := "DESC"
	if strings.EqualFold(req.SortDir, "asc") {
		sortDir = "ASC"
	}
	offset := (req.Page - 1) * req.Limit

	var result []Model
	err := query.
		Preload("Office").
		Preload("User").
		Order(fmt.Sprintf("%s %s", sortCol, sortDir)).
		Limit(req.Limit).
		Offset(offset).
		Find(&result).Error

	return result, total, err
}

// AssignLeader = INTI logic modul ini. Dalam SATU transaction + row lock:
//  1. Kalau office target udah punya leader aktif -> tutup masa jabatannya
//     (end_date = tanggal mulai leader baru). Ini "suksesi", history kejaga.
//  2. Kalau user yang mau di-assign udah aktif mimpin office LAIN -> TOLAK.
//     Gak auto-close, karena itu keputusan admin yang harus eksplisit
//     (end-term dulu di office lama, baru assign ke office baru).
//  3. Insert record baru dengan end_date NULL.
//
// Row lock (SELECT ... FOR UPDATE) WAJIB ada di sini -- tanpa ini, 2 request
// assign-leader ke office yang sama bisa keduanya lolos cek "belum ada
// leader aktif" sebelum salah satu selesai nulis, dan berakhir dengan 2
// leader aktif sekaligus (baru ketauan pas partial unique index di DB
// nolak salah satunya dengan error yang membingungkan).
func (r *repositoryImpl) AssignLeader(ctx context.Context, newLeader *Model) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var currentOfficeLeader Model
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("office_id = ? AND end_date IS NULL", newLeader.OfficeID).
			First(&currentOfficeLeader).Error

		switch {
		case err == nil:
			// Ada leader aktif di office ini -> tutup masa jabatannya.
			if err := tx.Model(&Model{}).
				Where("id = ?", currentOfficeLeader.ID).
				Update("end_date", newLeader.StartDate).Error; err != nil {
				return err
			}
		case errors.Is(err, gorm.ErrRecordNotFound):
			// Office belum pernah punya leader, gak masalah, lanjut.
		default:
			return err
		}

		var userActive Model
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND end_date IS NULL", newLeader.UserID).
			First(&userActive).Error

		switch {
		case err == nil:
			return fmt.Errorf("user %d is already actively leading office %d, end that term first", newLeader.UserID, userActive.OfficeID)
		case errors.Is(err, gorm.ErrRecordNotFound):
			// User belum aktif mimpin office manapun, aman buat di-assign.
		default:
			return err
		}

		return tx.Create(newLeader).Error
	})
}

func (r *repositoryImpl) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&Model{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repositoryImpl) EndTerm(ctx context.Context, id uint64, endDate time.Time) error {
	result := r.db.WithContext(ctx).Model(&Model{}).
		Where("id = ? AND end_date IS NULL", id).
		Update("end_date", endDate)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found or term already ended")
	}
	return nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Model{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
