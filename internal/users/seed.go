package users

import (
	"fmt"
	"strconv"
	"tms-be/internal/auth"
	"tms-be/internal/casbin"
	"tms-be/internal/conf"

	"gorm.io/gorm"
)

const seedDefaultPassword = "Password123!" // WAJIB diganti user pas login pertama di real project

// Seed bikin 1 user per role dari conf.AllRoles, idempotent (aman dijalankan
// berkali-kali). Nambah role baru? Cukup edit conf/roles.go — file ini gak
// perlu disentuh sama sekali.
// officeID WAJIB dikirim dari caller (cmd/seed/main.go) — karena
// users.OfficeID itu NOT NULL + foreign key ke tabel offices, jadi
// office seeder HARUS jalan duluan sebelum users.Seed dipanggil.
func Seed(db *gorm.DB, casbinSvc *casbin.Service, officeID uint64) error {
	hashed, err := auth.HashPassword(seedDefaultPassword)
	if err != nil {
		return fmt.Errorf("users seed: failed to hash default password: %w", err)
	}

	for _, role := range conf.AllRoles {
		email := fmt.Sprintf("%s@tms.local", role)

		user := Model{
			Username: string(role),
			Email:    email,
			Password: hashed,
			OfficeID: officeID,
			IsActive: true,
		}

		// FirstOrCreate by email -> idempotent, gak bikin duplikat kalau
		// seeder dijalankan ulang.
		result := db.Where(Model{Email: email}).FirstOrCreate(&user)
		if result.Error != nil {
			return fmt.Errorf("users seed: failed to seed %s: %w", email, result.Error)
		}

		// Subjek Casbin = user_id (string), konsisten sama seluruh flow.
		sub := strconv.FormatUint(user.ID, 10)
		if _, err := casbinSvc.GrantRole(sub, string(role)); err != nil {
			return fmt.Errorf("users seed: failed to grant role %s to %s: %w", role, email, err)
		}
	}

	return nil
}
