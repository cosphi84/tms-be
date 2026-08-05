package users

import (
	"fmt"
	"os"
	"strconv"
	"tms-be/internal/auth"
	"tms-be/internal/casbin"
	"tms-be/internal/conf"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB, casbinSvc *casbin.Service, office uint64) error {
	su_name := os.Getenv("SU_USERNAME")
	su_pass := os.Getenv("SU_PASSWORD")
	su_email := os.Getenv("SU_EMAIL")
	su_role := conf.RoleSuperadmin

	if su_name == "" || su_pass == "" || su_email == "" {
		return fmt.Errorf("env SU_USERNAME or SU_PASSWORD or SU_EMAIL are empty")
	}

	hashed, err := auth.HashPassword(su_pass)
	if err != nil {
		return fmt.Errorf("users seed: failed to hash default password: %w", err)
	}

	if office == 0 {
		return fmt.Errorf("users seed: office is zero")
	}
	user := Model{
		Username: su_name,
		Email:    su_email,
		Password: hashed,
		OfficeID: office,
		IsActive: true,
	}

	result := db.Where(Model{Email: su_email}).FirstOrCreate(&user)
	if result.Error != nil {
		return fmt.Errorf("users seed: failed to seed %s: %w", su_email, result.Error)
	}

	// Subjek Casbin = user_id (string), konsisten sama seluruh flow.
	sub := strconv.FormatUint(user.ID, 10)
	if _, err := casbinSvc.GrantRole(sub, string(su_role)); err != nil {
		return fmt.Errorf("users seed: failed to grant role %s to %s: %w", su_role, su_email, err)
	}

	return nil
}
