package helpers

import (
	"fmt"
	"time"
)

func BuildPath(folder, fileUUID, ext string) string {
	now := time.Now()
	return fmt.Sprintf("%s/%04d/%02d/%s%s", folder, now.Year(), int(now.Month()), fileUUID, ext)
}
