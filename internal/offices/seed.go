package offices

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
)

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

	log.Println("offices seed: done")

	return hq.ID, nil
}
