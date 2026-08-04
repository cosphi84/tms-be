package stocks

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tms-be/internal/pagination"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ListFilter struct {
	ToolsID     *uint64
	SlocID      *uint64
	ExpiredOnly bool
}

type Repository interface {
	FindByID(ctx context.Context, id uint64) (Model, error)
	FindByToolAndSloc(ctx context.Context, toolsID, slocID uint64) (Model, error)
	FindAllPaginated(ctx context.Context, filter ListFilter, req pagination.DtoPaginationRequest) ([]Model, int64, error)
	IncreaseStock(ctx context.Context, toolsID, slocID uint64, qty int, expiredDate time.Time, remarks *string) (Model, error)
	DecreaseStock(ctx context.Context, toolsID, slocID uint64, qty int) error
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

func (r *repositoryImpl) FindByToolAndSloc(ctx context.Context, toolsID, slocID uint64) (Model, error) {
	return gorm.G[Model](r.db).
		Where("tools_id = ? AND slock_ic = ?", toolsID, slocID).
		Take(ctx)
}

func (r *repositoryImpl) FindAllPaginated(ctx context.Context, filter ListFilter, req pagination.DtoPaginationRequest) ([]Model, int64, error) {
	query := r.db.WithContext(ctx).Model(&Model{})

	if filter.ToolsID != nil {
		query = query.Where("tools_id = ?", *filter.ToolsID)
	}
	if filter.SlocID != nil {
		query = query.Where("slock_ic = ?", *filter.SlocID)
	}
	if filter.ExpiredOnly {
		query = query.Where("expired_date IS NOT NULL AND expired_date < ?", time.Now())
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit

	var result []Model
	err := query.
		Preload("Tool").
		Preload("Sloc").
		Order("start_date DESC").
		Limit(req.Limit).
		Offset(offset).
		Find(&result).Error

	return result, total, err
}

// IncreaseStock = INTI logic modul ini. Row-lock di pasangan (tools_id,
// slock_ic) biar 2 request "barang masuk" bersamaan ke kombinasi yang sama
// gak keduanya baca qty lama sebelum salah satu selesai nulis (lost update).
//
// Kalau row belum ada -> insert baru (kombinasi ini emang unik, dijamin
// juga sama UNIQUE index di migrasi sebagai backstop terakhir).
// Kalau udah ada -> qty ditambah, cont_stock naik, start_date & expired_date
// di-refresh (barang baru masuk = masa pakai dihitung ulang dari sekarang).
func (r *repositoryImpl) IncreaseStock(ctx context.Context, toolsID, slocID uint64, qty int, expiredDate time.Time, remarks *string) (Model, error) {
	var result Model

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing Model
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("tools_id = ? AND slock_ic = ?", toolsID, slocID).
			First(&existing).Error

		switch {
		case err == nil:
			newQty := existing.Qty + qty
			updates := map[string]interface{}{
				"qty":          newQty,
				"cont_stock":   existing.ContStock + 1,
				"start_date":   time.Now(),
				"expired_date": expiredDate,
			}
			if remarks != nil {
				updates["remarks"] = *remarks
			}
			if err := tx.Model(&Model{}).Where("id = ?", existing.ID).Updates(updates).Error; err != nil {
				return err
			}
			// re-fetch biar result yang di-return up-to-date
			return tx.Where("id = ?", existing.ID).First(&result).Error

		case errors.Is(err, gorm.ErrRecordNotFound):
			newStock := Model{
				ToolsID:     toolsID,
				SlocID:      slocID,
				Qty:         qty,
				StartDate:   time.Now(),
				ExpiredDate: &expiredDate,
				ContStock:   1,
				Remarks:     remarks,
			}
			if err := tx.Create(&newStock).Error; err != nil {
				return err
			}
			result = newStock
			return nil

		default:
			return err
		}
	})

	return result, err
}

// DecreaseStock = pemakaian barang. TIDAK menyentuh start_date/expired_date/
// cont_stock -- itu cuma berubah kalau ada barang BARU masuk (IncreaseStock).
func (r *repositoryImpl) DecreaseStock(ctx context.Context, toolsID, slocID uint64, qty int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing Model
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("tools_id = ? AND slock_ic = ?", toolsID, slocID).
			First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("stock record not found for this tool at this location")
			}
			return err
		}

		if existing.Qty < qty {
			return fmt.Errorf("insufficient stock: available %d, requested %d", existing.Qty, qty)
		}

		return tx.Model(&Model{}).Where("id = ?", existing.ID).
			Update("qty", existing.Qty-qty).Error
	})
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
