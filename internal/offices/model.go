package offices

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type OfficeModel struct {
	ID       uint64        `gorm:"primaryKey;autoIncrement" json:"id" nestedset:"id"`
	ParentID sql.NullInt64 `gorm:"column:parent_id;index:idx_offices_parent_id" json:"parent_id" nestedset:"parent_id"`
	Code     string        `gorm:"type:varchar(10);unique;not null;index:idx_offices_code" json:"code"`
	Name     string        `gorm:"type:varchar(100);not null" json:"name"`
	Type     string        `gorm:"type:varchar(50);not null;check:type IN ('cabang', 'sdss', 'ssr', 'sass', 'tc', 'hq')" json:"type"`

	Depth         int `gorm:"column:depth;type:integer" json:"depth" nestedset:"depth"`
	Rgt           int `gorm:"column:rgt;type:integer" json:"rgt" nestedset:"rgt"`
	Lft           int `gorm:"column:lft;type:integer" json:"lft" nestedset:"lft"`
	ChildrenCount int `gorm:"column:children_count;type:integer" json:"children_count" nestedset:"children_count"`

	CreatedAt time.Time       `gorm:"not null;default:now();type:timestamp with time zone" json:"created_at"`
	CreatedBy *uint64         `gorm:"column:created_by;type:integer" json:"created_by,omitempty"`
	UpdatedAt *time.Time      `gorm:"type:timestamp with time zone" json:"updated_at,omitempty"`
	UpdatedBy *uint64         `gorm:"column:updated_by;type:integer" json:"updated_by,omitempty"`
	DeletedAt *gorm.DeletedAt `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
	DeletedBy *uint64         `gorm:"column:deleted_by;type:integer" json:"deleted_by,omitempty"`

	Parent   *OfficeModel  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []OfficeModel `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (OfficeModel) TableName() string {
	return "offices"
}
