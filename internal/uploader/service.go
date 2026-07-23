package uploader

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

var (
	ErrFileMetaNotFound = errors.New("file: metadata not found")
)

type UploadService interface {
	UploadFile(ctx context.Context, folder string, fileHeader *multipart.FileHeader) (*UploadResponseDto, error)
	GetMetadata(uuid string) (*UploadResponseDto, error)
	OpenForDownload(ctx context.Context, uuid string) (file interface{ Read([]byte) (int, error) }, meta *UploaderModel, err error)
	DeleteFile(ctx context.Context, uuid string) error
	ReplaceWithTx(ctx context.Context, tx *gorm.DB, folder string, fh *multipart.FileHeader, oldFileUUID string) (newFile *UploadResponseDto, err error)
	FinalizeReplace(ctx context.Context, oldFileUUID string) error
}

type fileService struct {
	repo    UploaderRepository
	storage Storage
}

func NewUploadService(repo UploaderRepository, storage Storage) UploadService {
	return &fileService{
		repo:    repo,
		storage: storage,
	}
}

func buildPath(folder, fileUUID, ext string) string {
	now := time.Now()
	return fmt.Sprintf("%s/%04d/%02d/%s%s", folder, now.Year(), int(now.Month()), fileUUID, ext)
}

func toUploadResponse(m *UploaderModel) *UploadResponseDto {
	return &UploadResponseDto{
		ID:           m.ID,
		UUID:         m.UUID,
		OriginalName: m.OriginalName,
		MimeType:     m.MimeType,
		Extension:    m.Extension,
		Size:         m.Size,
		Checksum:     m.Checksum,
		URL:          fmt.Sprintf("/files/%s", m.UUID),
		IsArchived:   m.IsArchived,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (s *fileService) UploadFile(ctx context.Context, folder string, fileHeader *multipart.FileHeader) (*UploadResponseDto, error) {
	usr, err := s.repo.GetCurrentUser(ctx)
	return nil, nil
}
