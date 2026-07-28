package auth

import (
	"time"
)

type Model struct {
	ID                 uint64     `gorm:"primaryKey;autoIncrement;not null;column:id;<-:create" json:"id"`
	Username           string     `gorm:"not null;column:username;index:idx_user_username;type:varchar(255)" json:"username"`
	Email              string     `gorm:"unique;not null;column:email;index:idx_user_email" json:"email"`
	Password           string     `gorm:"not null;column:password" json:"-"`
	OfficeID           uint64     `gorm:"not null;column:office_id;index:idx_user_office_id" json:"office_id"`
	IsActive           bool       `gorm:"not null;default:true;index:idx_users_is_active" json:"is_active"`
	FailedLoginAttempt int        `gorm:"not null;default:0" json:"failed_login_attempts"`
	LockedUntil        *time.Time `gorm:"type:timestamp" json:"locked_until,omitempty"`
	LastLoginAt        *time.Time `gorm:"type:timestamp" json:"last_login_at,omitempty"`
	LastLoginFrom      *string    `gorm:"type:varchar" json:"last_login_from,omitempty"`
}

func (Model) TableName() string {
	return "users"
}
