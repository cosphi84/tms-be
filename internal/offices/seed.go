package offices

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Seed bikin HQ (root, satu-satunya) + beberapa cabang contoh di bawahnya.
// Idempotent -- aman dijalankan berkali-kali, gak bakal duplikat.
//
// Return value: ID office HQ, dipakai caller (cmd/seed/main.go) sebagai
// officeID default buat users.Seed() -- karena users.OfficeID itu NOT NULL
// + FK ke tabel offices, office HARUS ke-seed duluan sebelum users.
func Seed(db *gorm.DB) (uint64, error) {
	repo := NewRepository(db)
	svc := NewService(repo)
	ctx := context.Background()

	log.Println("seeding offices: HQ")
	hq, err := svc.CreateHQ(ctx, HQRequestDTO{
		Code: "81",
		Name: "Head Quarter",
		Type: "hq",
	})
	if err != nil {
		return 0, fmt.Errorf("offices seed: failed to seed HQ: %w", err)
	}
	log.Printf("offices seed: HQ ready (id=%d)", hq.ID)

	log.Println("seeding offices: sample branches")
	sampleBranches := []OfficeRequestDTO{
		{Code: "8101", Name: "TC", Type: "tc", ParentID: &hq.ID},
		{Code: "10", Name: "SWADAYA", Type: "cabang", ParentID: &hq.ID},
	}

	for _, branch := range sampleBranches {
		// Idempotent: skip kalau code udah ada, biar seeder aman di-rerun.
		existing, err := svc.FindOffice(ctx, branch.Code)
		if err != nil {
			return 0, fmt.Errorf("offices seed: failed to check %s: %w", branch.Code, err)
		}
		if existing != nil {
			continue
		}
		if err := svc.CreateOffice(ctx, branch); err != nil {
			return 0, fmt.Errorf("offices seed: failed to seed %s: %w", branch.Code, err)
		}
	}
	log.Println("offices seed: done")

	return hq.ID, nil
}
