package helpers

import (
	"context"
	"tms-be/internal/auth"
)

func GetLoggedUserID(ctx context.Context) (uint64, error) {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
