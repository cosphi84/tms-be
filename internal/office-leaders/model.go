package officeleaders

import (
	"time"
	"tms-be/internal/offices"
	"tms-be/internal/users"

	"gorm.io/gorm"
)

type Model struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	OfficeID uint64         `gorm:"not null;column:office_id;index:idx_office_leaders_office_id" json:"office_id"`
	Office   *offices.Model `gorm:"foreignKey:OfficeID;references:ID" json:"office,omitempty"`

	UserID uint64       `gorm:"not null;column:user_id;index:idx_office_leaders_user_id" json:"user_id"`
	User   *users.Model `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`

	StartDate time.Time  `gorm:"not null;type:date;column:start_date" json:"start_date"`
	EndDate   *time.Time `gorm:"type:date;column:end_date" json:"end_date,omitempty"` // nil = masih menjabat

	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	CreatedBy *uint64        `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	UpdatedBy *uint64        `gorm:"column:updated_by" json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"` // non-pointer, WAJIB biar GORM auto soft-delete
	DeletedBy *uint64        `gorm:"column:deleted_by" json:"deleted_by,omitempty"`
}

func (Model) TableName() string {
	return "office_leaders"
}

// IsActive: true kalau belum ada penerus (EndDate masih nil).
func (m Model) IsActive() bool {
	return m.EndDate == nil
}
