package category

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID          uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:varchar(255);" json:"description"`

	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	CreatedBy *uint64        `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	UpdatedBy *uint64        `gorm:"column:updated_by" json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"` // non-pointer, wajib biar GORM auto soft-delete
	DeletedBy *uint64        `gorm:"column:deleted_by" json:"deleted_by,omitempty"`
}

func (Model) TableName() string {
	return "tool_categories"
}
