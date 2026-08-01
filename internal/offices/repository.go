package offices

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"tms-be/internal/pagination"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateHQ(ctx context.Context, office *Model) error
	CreateBranch(ctx context.Context, office *Model, parentID uint64) error
	FindAllByLevel(ctx context.Context, level string, req pagination.DtoPaginationRequest) ([]Model, int64, error)
	FindAllFlat(ctx context.Context) ([]Model, error)
	FindChildren(ctx context.Context, parentID uint64) ([]Model, error)
	FindByID(ctx context.Context, id uint64) (Model, error)
	FindOffice(ctx context.Context, officeCode string) (Model, error)
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint64) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) FindOffice(ctx context.Context, officeCode string) (Model, error) {
	return gorm.G[Model](r.db).
		Where("code = ?", officeCode).
		Take(ctx)
}

// CreateHQ bikin root node. HANYA boleh ada 1 root (parent_id IS NULL).
// Dipanggil cuma dari seeder (lihat seed.go) — bukan lewat HTTP endpoint.
func (r *repositoryImpl) CreateHQ(ctx context.Context, office *Model) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&Model{}).Where("parent_id IS NULL").Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("HQ root already exists, cannot create another root")
		}

		office.ParentID = nil
		office.Lft = 1
		office.Rgt = 2
		office.Depth = 0
		office.ChildrenCount = 0

		return tx.Create(office).Error
	})
}

// CreateBranch = INTI algoritma nested set insert. Harus jalan dalam
// transaction + row lock di parent, karena kalau ada 2 insert bersamaan
// ke parent yang sama tanpa lock, lft/rgt seluruh tabel bisa KORUP
// (dua goroutine baca parent.Rgt yang sama, geser interval yang tumpang
// tindih -> data hierarki jadi gak konsisten).
func (r *repositoryImpl) CreateBranch(ctx context.Context, office *Model, parentID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var parent Model
		// SELECT ... FOR UPDATE -- kunci row parent sampai transaction ini
		// selesai, biar insert lain ke parent yang sama harus antri.
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", parentID).
			First(&parent).Error; err != nil {
			return err
		}

		insertPoint := parent.Rgt

		// Geser rgt semua node (termasuk parent & seluruh ancestor-nya)
		// yang rgt-nya >= titik insert, kasih ruang 2 slot buat node baru.
		if err := tx.Model(&Model{}).
			Where("rgt >= ?", insertPoint).
			UpdateColumn("rgt", gorm.Expr("rgt + 2")).Error; err != nil {
			return err
		}
		// Geser lft semua node yang lft-nya SETELAH titik insert (sibling
		// yang lebih "kanan" dari parent).
		if err := tx.Model(&Model{}).
			Where("lft > ?", insertPoint).
			UpdateColumn("lft", gorm.Expr("lft + 2")).Error; err != nil {
			return err
		}

		office.ParentID = &parentID
		office.Lft = insertPoint
		office.Rgt = insertPoint + 1
		office.Depth = parent.Depth + 1
		office.ChildrenCount = 0

		if err := tx.Create(office).Error; err != nil {
			return err
		}

		return tx.Model(&Model{}).
			Where("id = ?", parent.ID).
			UpdateColumn("children_count", gorm.Expr("children_count + 1")).Error
	})
}

var officeSortableColumns = map[string]bool{
	"code":       true,
	"name":       true,
	"type":       true,
	"created_at": true,
}

func (r *repositoryImpl) FindAllByLevel(ctx context.Context, level string, req pagination.DtoPaginationRequest) ([]Model, int64, error) {
	query := r.db.WithContext(ctx).Model(&Model{}).Where("type = ?", level)

	if req.Search != "" {
		like := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortCol := "code"
	if officeSortableColumns[req.SortedBy] {
		sortCol = req.SortedBy
	}
	sortDir := "ASC"
	if strings.EqualFold(req.SortDir, "desc") {
		sortDir = "DESC"
	}
	offset := (req.Page - 1) * req.Limit

	var result []Model
	err := query.
		Order(fmt.Sprintf("%s %s", sortCol, sortDir)).
		Limit(req.Limit).
		Offset(offset).
		Find(&result).Error

	return result, total, err
}

// FindAllFlat ambil SEMUA office (semua level, termasuk HQ), diurutkan
// lft ASC (urutan preorder alami nested set). Dipakai buat bangun tree
// di service.go, BUKAN buat ditampilkan langsung.
func (r *repositoryImpl) FindAllFlat(ctx context.Context) ([]Model, error) {
	var all []Model
	err := r.db.WithContext(ctx).Order("lft ASC").Find(&all).Error
	return all, err
}

// FindChildren = anak LANGSUNG (bukan seluruh descendant), diurutkan by code.
func (r *repositoryImpl) FindChildren(ctx context.Context, parentID uint64) ([]Model, error) {
	var children []Model
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("code ASC").
		Find(&children).Error
	return children, err
}

func (r *repositoryImpl) FindByID(ctx context.Context, id uint64) (Model, error) {
	return gorm.G[Model](r.db).Where("id = ?", id).Take(ctx)
}

// Update SENGAJA cuma terima map[string]interface{} (bukan full struct
// Updates) -- field yang boleh diubah cuma code/name/type. parent_id TIDAK
// bisa diubah lewat sini, karena pindah posisi node di nested set butuh
// re-kalkulasi lft/rgt seluruh subtree yang dipindah, itu operasi terpisah
// yang lebih kompleks (di luar scope endpoint Update biasa).
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

// Delete cuma boleh buat LEAF node (children_count == 0). Soft-delete
// (lewat gorm.DeletedAt) -- lft/rgt node yang dihapus SENGAJA gak
// di-compact/di-geser ulang, cukup ditinggal sebagai celah numerik.
// Ini aman karena semua query lain otomatis exclude row yang soft-deleted,
// jadi celah itu gak pernah kebaca sebagai bagian tree yang valid -- cuma
// "buang" sedikit ruang angka, bukan bug.
func (r *repositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var office Model
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", id).First(&office).Error; err != nil {
			return err
		}

		if office.ChildrenCount > 0 {
			return fmt.Errorf("cannot delete office %d: still has %d child office(s), move or delete them first", id, office.ChildrenCount)
		}

		if err := tx.Delete(&office).Error; err != nil {
			return err
		}

		if office.ParentID != nil {
			if err := tx.Model(&Model{}).
				Where("id = ?", *office.ParentID).
				UpdateColumn("children_count", gorm.Expr("children_count - 1")).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
