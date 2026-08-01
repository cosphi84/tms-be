package offices

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id" nestedset:"id"`

	// *uint64, BUKAN sql.NullInt64 — konsisten sama CreatedBy/UpdatedBy/DeletedBy
	// di bawah, dan lebih aman buat GORM self-referencing FK (ID juga uint64).
	// nil = root (HQ), non-nil = office biasa yang punya parent.
	ParentID *uint64 `gorm:"column:parent_id;index:idx_offices_parent_id" json:"parent_id" nestedset:"parent_id"`

	Code string `gorm:"type:varchar(10);unique;not null;index:idx_offices_code" json:"code"`
	Name string `gorm:"type:varchar(100);not null" json:"name"`
	Type string `gorm:"type:varchar(50);not null;check:type IN ('cabang', 'sdss', 'ssr', 'sass', 'tc', 'hq')" json:"type"`

	Depth         int `gorm:"column:depth;type:integer" json:"depth" nestedset:"depth"`
	Rgt           int `gorm:"column:rgt;type:integer" json:"rgt" nestedset:"rgt"`
	Lft           int `gorm:"column:lft;type:integer" json:"lft" nestedset:"lft"`
	ChildrenCount int `gorm:"column:children_count;type:integer" json:"children_count" nestedset:"children_count"`

	CreatedAt time.Time  `gorm:"not null;default:now();type:timestamp with time zone" json:"created_at"`
	CreatedBy *uint64    `gorm:"column:created_by;type:integer" json:"created_by,omitempty"`
	UpdatedAt *time.Time `gorm:"type:timestamp with time zone" json:"updated_at,omitempty"`
	UpdatedBy *uint64    `gorm:"column:updated_by;type:integer" json:"updated_by,omitempty"`

	// gorm.DeletedAt TANPA pointer — GORM cuma ngenalin tipe ini persis buat
	// auto soft-delete. Kalau di-pointer-kan (*gorm.DeletedAt kayak
	// sebelumnya), plugin soft-delete GORM gak akan detect field ini sama
	// sekali, dan .Delete() akan HARD delete tanpa kamu sadar.
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
	DeletedBy *uint64        `gorm:"column:deleted_by;type:integer" json:"deleted_by,omitempty"`

	Parent   *Model  `gorm:"foreignKey:ParentID;references:ID" json:"parent,omitempty"`
	Children []Model `gorm:"foreignKey:ParentID;references:ID" json:"children,omitempty"`
}

func (Model) TableName() string {
	return "offices"
}
