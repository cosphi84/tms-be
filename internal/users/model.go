package users

import (
	"time"
	"tms-be/internal/offices"

	"gorm.io/gorm"
)

type Model struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement;not null;column:id;<-:create" json:"id"`
	Username    string         `gorm:"not null;column:username;index:idx_user_username;type:varchar(255)" json:"username"`
	Email       string         `gorm:"unique;not null;column:email;index:idx_user_email" json:"email"`
	Password    string         `gorm:"not null;column:password" json:"-"`
	FotoProfile *string        `gorm:"column:foto_profile" json:"foto_profile,omitempty"`
	OfficeID    uint64         `gorm:"not null;column:office_id;index:idx_user_office_id" json:"office_id"`
	Office      *offices.Model `gorm:"foreignKey:OfficeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"office,omitempty"`

	IsActive           bool       `gorm:"not null;default:true;index:idx_users_is_active" json:"is_active"`
	FailedLoginAttempt int        `gorm:"not null;default:0" json:"failed_login_attempts"`
	LockedUntil        *time.Time `gorm:"type:timestamp" json:"locked_until,omitempty"`
	LastLoginAt        *time.Time `gorm:"type:timestamp" json:"last_login_at,omitempty"`
	LastLoginFrom      *string    `gorm:"type:varchar" json:"last_login_from,omitempty"`

	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp" json:"deleted_at,omitempty"`

	// *uint64, BUKAN *int64 — harus match tipe claims.UserID (uint64) dari
	// auth.JWTClaims, karena diisi dari &loggedUser.UserID di service.go
	CreatedBy *uint64 `gorm:"column:created_by" json:"created_by,omitempty"`
}

func (Model) TableName() string {
	return "users"
}
