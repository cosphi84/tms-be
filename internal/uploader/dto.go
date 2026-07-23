package uploader

import "time"

type UploadRequestDto struct {
	Folder string `form:"folder" binding:"required"`
}

type UploadResponseDto struct {
	ID           int64     `json:"id"`
	UUID         string    `json:"uuid"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	Extension    string    `json:"extension"`
	Size         int64     `json:"size"`
	Checksum     string    `json:"checksum"`
	URL          string    `json:"url"`
	IsArchived   bool      `json:"is_archived"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ReplaceRequestDto struct {
	Folder string `form:"folder" binding:"required"`
}
