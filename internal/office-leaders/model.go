package officeleaders

import (
	"time"
	"tms-be/internal/offices"
	"tms-be/internal/users"
)

type OfficeLeaderModel struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OfficeID  uint64     `gorm:"column:office_id;not null;index:idx_office_leaders_office_id" json:"office_id"`
	UserID    uint64     `gorm:"column:user_id;not null;index:idx_office_leaders_user_id" json:"user_id"`
	StartDate *time.Time `gorm:"column:start_date;type:date;not null" json:"start_date,omitempty"`
	EndDate   *time.Time `gorm:"column:end_date;type:date" json:"end_date,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp with time zone" json:"updated_at,omitempty"`

	Office offices.OfficeModel `gorm:"foreignKey:OfficeID" json:"office"`
	Leader []users.UserModel   `gorm:"foreignKey:UserID" json:"leader"`
}

func (OfficeLeaderModel) TableName() string {
	return "office_leaders"
}
