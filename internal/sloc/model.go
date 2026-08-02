package sloc

import (
	"slices"
	"time"
	"tms-be/internal/offices"

	"gorm.io/gorm"
)

type MemberType string

const (
	MemberSEID MemberType = "SEID"
	MemberOTS  MemberType = "OTS"
	MemberSASS MemberType = "SASS"
)

// AllMembers = satu-satunya tempat daftar member type didefinisikan --
// sama pola kayak conf.AllRoles. dto.go baca dari sini lewat custom
// validator, gak hardcode ulang.
var AllMembers = []MemberType{MemberSEID, MemberOTS, MemberSASS}

func (m MemberType) Valid() bool {
	return slices.Contains(AllMembers, m)
}

type Model struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	Code string `gorm:"type:varchar(20);not null;index:idx_sloc_code" json:"code"`
	Name string `gorm:"type:varchar(100);not null" json:"name"`

	OfficeID uint64         `gorm:"not null;column:office_id;index:idx_sloc_office_id" json:"office_id"`
	Office   *offices.Model `gorm:"foreignKey:OfficeID;references:ID" json:"office,omitempty"`

	Member string `gorm:"type:varchar(10);not null;check:member IN ('SEID','OTS','SASS')" json:"member"`

	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	CreatedBy *uint64        `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	UpdatedBy *uint64        `gorm:"column:updated_by" json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"` // non-pointer, wajib biar GORM auto soft-delete
	DeletedBy *uint64        `gorm:"column:deleted_by" json:"deleted_by,omitempty"`
}

func (Model) TableName() string {
	return "storage_locations"
}
