package stocks

import (
	"time"
	mtools "tms-be/internal/mtools"
	sloc "tms-be/internal/sloc"

	"gorm.io/gorm"
)

type Model struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	ToolsID uint64        `gorm:"not null;column:tools_id;index:idx_stocks_tools_id" json:"tools_id"`
	Tool    *mtools.Model `gorm:"foreignKey:ToolsID;references:ID" json:"tool,omitempty"`

	// Kolom fisik "slock_id" -- sesuai migrasi yang kamu kasih.
	SlocID uint64      `gorm:"not null;column:slock_id;index:idx_stocks_slock_id" json:"sloc_id"`
	Sloc   *sloc.Model `gorm:"foreignKey:SlocID;references:ID" json:"sloc,omitempty"`

	Qty int `gorm:"not null;column:qty" json:"qty"`

	StartDate   time.Time  `gorm:"not null;column:start_date;default:now()" json:"start_date"`
	ExpiredDate *time.Time `gorm:"column:expired_date" json:"expired_date,omitempty"`

	// Berapa kali stok ini di-top-up/kedatangan barang baru (mulai dari 1).
	ContStock int `gorm:"not null;default:1;column:cont_stock" json:"cont_stock"`

	Remarks *string `gorm:"column:remarks" json:"remarks,omitempty"`

	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	CreatedBy *uint64        `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	UpdatedBy *uint64        `gorm:"column:updated_by" json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
	DeletedBy *uint64        `gorm:"column:deleted_by" json:"deleted_by,omitempty"`
}

func (Model) TableName() string {
	return "stocks"
}

// IsExpired: helper buat cek status expired tanpa perlu logic tambahan di
// service/handler.
func (m Model) IsExpired() bool {
	return m.ExpiredDate != nil && m.ExpiredDate.Before(time.Now())
}
