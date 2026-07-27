package uploader

import (
	"time"
	"tms-be/internal/users"

	"gorm.io/gorm"
)

type UploaderModel struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID string `gorm:"type:uuid;not null;default:gen_random_uuid();uniqueIndex:idx_asset_files_uuid" json:"uuid"`

	DiskName     string `gorm:"type:varchar(255);not null" json:"disk_name"`
	OriginalName string `gorm:"type:varchar(255);not null" json:"original_name"`
	MimeType     string `gorm:"type:varchar(100);not null" json:"mime_type"`
	Extension    string `gorm:"type:varchar(20);not null" json:"extension"`
	Size         int64  `gorm:"not null" json:"size"`
	Checksum     string `gorm:"type:varchar(64);not null;index:idx_asset_files_checksum" json:"checksum"`

	Path    string `gorm:"type:varchar(500);not null" json:"path"`
	Storage string `gorm:"type:varchar(50);not null;default:local" json:"storage"`

	IsArchived   bool    `gorm:"not null;default:false;index:idx_asset_files_is_archived" json:"is_archived"`
	ArchivedPath *string `gorm:"type:varchar(500)" json:"archived_path,omitempty"`

	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_asset_files_deleted_at" json:"deleted_at"`
	CreatedBy *int64         `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedBy *int64         `gorm:"column:updated_by" json:"updated_by,omitempty"`
	DeletedBy *int64         `gorm:"column:deleted_by" json:"deleted_by,omitempty"`

	// Relations — mengikuti pola audit-user di StorageLocation
	CreatedByUser *users.Model `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	UpdatedByUser *users.Model `gorm:"foreignKey:UpdatedBy" json:"updated_by_user,omitempty"`
	DeletedByUser *users.Model `gorm:"foreignKey:DeletedBy" json:"deleted_by_user,omitempty"`
}

func (UploaderModel) TableName() string {
	return "asset_files"
}
