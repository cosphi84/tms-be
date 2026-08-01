package uploader

import "time"

type UploadRequestDTO struct {
	Folder string `form:"folder" binding:"required"`
	Access bool   `form:"access" binding:"required"`
}

type UploadResponseDTO struct {
	ID           uint64    `json:"id"`
	UUID         string    `json:"uuid"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	Extension    string    `json:"extension"`
	Size         uint64    `json:"size"`
	IsPublic     bool      `json:"is_public"`
	Checksum     string    `json:"checksum"`
	URL          string    `json:"url"`
	IsArchived   bool      `json:"is_archived"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ReplaceRequestDTO struct {
	Folder string `form:"folder" binding:"required"`
	Access bool   `form:"access" binding:"required"`
}
