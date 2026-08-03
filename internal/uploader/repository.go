package uploader

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, file *Model) error
	FindByUUID(ctx context.Context, uuid string) (Model, error)
	FindByID(ctx context.Context, id uint64) (Model, error)
	Update(ctx context.Context, id uint64, file Model) (int, error)
	SoftDelete(ctx context.Context, mdl Model) (int, error)

	// WithTx mengembalikan instance FileRepository baru yang beroperasi
	// di dalam transaction milik caller (mis. ToolsService.ReplaceIcon).
	// Ini memungkinkan File Manager Module ikut serta dalam atomic
	// transaction lintas-modul tanpa File Manager perlu tahu apa-apa
	// soal domain caller — caller yang mengontrol Begin/Commit/Rollback.
	WithTx(tx *gorm.DB) Repository
}

type uploadRepos struct {
	db *gorm.DB
}

func NewUploaderRepository(db *gorm.DB) Repository {
	return &uploadRepos{db: db}
}

func (upd *uploadRepos) Create(ctx context.Context, file *Model) error {
	err := gorm.G[Model](upd.db).Create(ctx, file)
	return err
}

func (upd *uploadRepos) FindByUUID(ctx context.Context, uuid string) (Model, error) {
	file, err := gorm.G[Model](upd.db).Where("uuid = ?", uuid).Take(ctx)
	return file, err
}

func (upd *uploadRepos) Update(ctx context.Context, id uint64, file Model) (int, error) {
	return gorm.G[Model](upd.db).Where("id = ?", id).Select("*").Updates(ctx, file)
}

func (upd *uploadRepos) SoftDelete(ctx context.Context, mdl Model) (int, error) {
	return gorm.G[Model](upd.db).
		Select("deleted_at", "deleted_by").
		Where("id = ?", mdl.ID).
		Updates(ctx, mdl)

}

func (upd *uploadRepos) FindByID(ctx context.Context, id uint64) (Model, error) {
	file, err := gorm.G[Model](upd.db).Where("id = ?", id).Take(ctx)
	return file, err
}

func (upd *uploadRepos) WithTx(tx *gorm.DB) Repository {
	return &uploadRepos{db: tx}
}
