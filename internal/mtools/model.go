package mtools

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
	"tms-be/internal/groups"
	toolcategories "tms-be/internal/tool-categories"
	"tms-be/internal/uploader"

	"gorm.io/gorm"
)

type Model struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	Code      string `gorm:"type:varchar(20);not null;index:idx_mtools_code" json:"code"`
	Name      string `gorm:"type:varchar(150);not null" json:"name"`
	Brand     string `gorm:"type:varchar(100);not null" json:"brand"`
	Type      string `gorm:"type:varchar(100);not null" json:"type"`
	SerialNum string `gorm:"type:varchar(100);column:serial_num" json:"serial_num"`

	CategoryID uint32                `gorm:"not null;column:category_id;index:idx_mtools_category_id" json:"category_id"`
	Category   *toolcategories.Model `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`

	GroupID *uint32       `gorm:"column:group_id;references:ID" json:"group_id,omitempty"`
	Group   *groups.Model `gorm:"foreignKey:GroupID;references:ID" json:"group,omitempty"`

	PhotoID *uint64         `gorm:"type:varchar(36);column:photo_uuid" json:"photo_id,omitempty"`
	Photo   *uploader.Model `gorm:"foreignKey:PhotoID;references:UUID" json:"photo,omitempty"`

	Price float64 `gorm:"type:numeric(15,2);not null" json:"price"`

	// Format singkat: "2y" (2 tahun), "6m" (6 bulan), "30d" (30 hari).
	// Dipakai module lain (stock tools) buat hitung expiry lewat
	// CalculateExpiry() di bawah -- BUKAN buat di-parse ulang di tempat lain.
	UsagePeriod string `gorm:"type:varchar(10);column:usage_period" json:"usage_period"`

	IsActive bool `gorm:"not null;default:true;index:idx_mtools_is_active" json:"is_active"`

	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	CreatedBy *uint64        `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	UpdatedBy *uint64        `gorm:"column:updated_by" json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
	DeletedBy *uint64        `gorm:"column:deleted_by" json:"deleted_by,omitempty"`
}

func (Model) TableName() string {
	return "master_tools"
}

var usagePeriodPattern = regexp.MustCompile(`^(\d+)([ymd])$`)

// ParseUsagePeriod: "2y" -> (2, 0, 0), "6m" -> (0, 6, 0), "30d" -> (0, 0, 30).
// Dipakai internal (CalculateExpiry) dan bisa diimpor langsung dari module
// lain (misal stock tools) kalau butuh validasi/parsing terpisah.
func ParseUsagePeriod(s string) (years, months, days int, err error) {
	match := usagePeriodPattern.FindStringSubmatch(s)
	if match == nil {
		return 0, 0, 0, fmt.Errorf("invalid usage_period format %q, expected like \"2y\", \"6m\", or \"30d\"", s)
	}

	n, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, 0, 0, err
	}

	switch match[2] {
	case "y":
		return n, 0, 0, nil
	case "m":
		return 0, n, 0, nil
	case "d":
		return 0, 0, n, nil
	default:
		return 0, 0, 0, fmt.Errorf("unknown usage_period unit %q", match[2])
	}
}

// CalculateExpiry: dipakai module stock tools buat hitung tanggal expired
// stok dari tanggal tertentu (misal tanggal masuk stok ke gudang), pakai
// AddDate biar kalkulasi kalender akurat (bukan cuma 24 jam x hari).
func (m Model) CalculateExpiry(from time.Time) (time.Time, error) {
	years, months, days, err := ParseUsagePeriod(m.UsagePeriod)
	if err != nil {
		return time.Time{}, err
	}
	return from.AddDate(years, months, days), nil
}
