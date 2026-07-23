package uploader

import (
	"context"

	"gorm.io/gorm"
)

type UploaderRepository interface {
	Create(ctx context.Context, file *UploaderModel) error
	FindByUUID(ctx context.Context, uuid string) (UploaderModel, error)
	Update(ctx context.Context, id int64, file UploaderModel) error
	SoftDelete(ctx context.Context, id int64) error

	// WithTx mengembalikan instance FileRepository baru yang beroperasi
	// di dalam transaction milik caller (mis. ToolsService.ReplaceIcon).
	// Ini memungkinkan File Manager Module ikut serta dalam atomic
	// transaction lintas-modul tanpa File Manager perlu tahu apa-apa
	// soal domain caller — caller yang mengontrol Begin/Commit/Rollback.
	WithTx(tx *gorm.DB) UploaderRepository
}

type uploaderRepository struct {
	db *gorm.DB
}

func NewUploaderRepository(db *gorm.DB) UploaderRepository {
	return &uploaderRepository{db: db}
}

func (r *uploaderRepository) Create(ctx context.Context, file *UploaderModel) error {
	err := gorm.G[UploaderModel](r.db).Create(ctx, file)
	return err
}

func (r *uploaderRepository) FindByUUID(ctx context.Context, uuid string) (UploaderModel, error) {
	file, err := gorm.G[UploaderModel](r.db).Where("uuid = ?", uuid).First(ctx)
	return file, err
}

func (r *uploaderRepository) Update(ctx context.Context, id int64, file UploaderModel) error {
	rows, err := gorm.G[UploaderModel](r.db).Where("id = ?", id).Updates(ctx, file)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *uploaderRepository) SoftDelete(ctx context.Context, id int64) error {
	rows, err := gorm.G[UploaderModel](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *uploaderRepository) WithTx(tx *gorm.DB) UploaderRepository {
	return &uploaderRepository{db: tx}
}
