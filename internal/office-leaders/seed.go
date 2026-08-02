package officeleaders

import (
	"context"
	"fmt"
	"log"
	"time"
	"tms-be/internal/offices"
	"tms-be/internal/users"

	"gorm.io/gorm"
)

// Seed bikin 1 contoh assignment: office "81001" (Cabang Jakarta, dari
// offices.Seed) dipimpin oleh user role service_head (dari users.Seed).
// Idempotent -- skip kalau office itu udah punya leader aktif.
//
// URUTAN WAJIB dijalankan setelah offices.Seed() DAN users.Seed(), karena
// butuh keduanya udah ada di DB.
func Seed(db *gorm.DB) error {
	repo := NewRepository(db)
	officeRepo := offices.NewRepository(db)
	userRepo := users.NewUserRepository(db)
	svc := NewService(repo, officeRepo, userRepo)
	ctx := context.Background()

	office, err := officeRepo.FindOffice(ctx, "81001")
	if err != nil {
		return fmt.Errorf("office-leaders seed: sample office 81001 not found, run offices.Seed first: %w", err)
	}

	leaderUser, err := userRepo.FindUser(ctx, "service_head@tms.local")
	if err != nil {
		return fmt.Errorf("office-leaders seed: sample leader user not found, run users.Seed first: %w", err)
	}

	existing, err := svc.ActiveByOffice(ctx, office.ID)
	if err != nil {
		return fmt.Errorf("office-leaders seed: failed to check existing leader: %w", err)
	}
	if existing != nil {
		log.Println("office-leaders seed: office 81001 already has an active leader, skipping")
		return nil
	}

	err = svc.AssignLeader(ctx, AssignLeaderDTO{
		OfficeID:  office.ID,
		UserID:    leaderUser.ID,
		StartDate: time.Now().Format(DateLayout),
	})
	if err != nil {
		return fmt.Errorf("office-leaders seed: failed to assign sample leader: %w", err)
	}

	log.Println("office-leaders seed: done")
	return nil
}
