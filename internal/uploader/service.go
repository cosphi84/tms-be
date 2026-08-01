package uploader

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"
	"tms-be/internal/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrFileMetaNotFound = errors.New("file: metadata not found")
)

type Service interface {
	UploadFile(ctx context.Context, dto UploadRequestDTO, fileHeader *multipart.FileHeader) (*UploadResponseDTO, error)
	GetMetadata(ctx context.Context, uuid string) (*UploadResponseDTO, error)
	OpenForDownload(ctx context.Context, uuid string) (file interface{ Read([]byte) (int, error) }, meta *Model, err error)
	DeleteFile(ctx context.Context, uuid string) error
	ReplaceWithTx(ctx context.Context, tx *gorm.DB, folder string, fh *multipart.FileHeader, oldFileUUID string) (newFile *UploadResponseDTO, err error)
	FinalizeReplace(ctx context.Context, oldFileUUID string) error
}

type fileService struct {
	repo      Repository
	storage   Storage
	validator FileValidator
}

func NewUploadService(repo Repository, storage Storage, validator FileValidator) Service {

	return &fileService{
		repo:      repo,
		storage:   storage,
		validator: validator,
	}
}

func toUploadResponse(m *Model) *UploadResponseDTO {
	return &UploadResponseDTO{
		ID:           m.ID,
		UUID:         m.UUID,
		OriginalName: m.OriginalName,
		MimeType:     m.MimeType,
		Extension:    m.Extension,
		Size:         m.Size,
		IsPublic:     m.IsPublic,
		Checksum:     m.Checksum,
		URL:          fmt.Sprintf("/files/%s", m.UUID),
		IsArchived:   m.IsArchived,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (s *fileService) UploadFile(ctx context.Context, dto UploadRequestDTO, fileHeader *multipart.FileHeader) (*UploadResponseDTO, error) {
	// Ambil siapa yang login
	loggedID, err := helpers.GetLoggedUserID(ctx)
	if err != nil {
		return nil, err
	}
	if loggedID == 0 {
		return nil, errors.New("no logged user")
	}

	// detect header
	ext, err := s.validator.ValidateHeader(fileHeader)
	if err != nil {
		return nil, err
	}

	// buka file
	f, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("file: failed to open uploaded file: %w", err)
	}
	defer f.Close()

	// validate content file
	if err := s.validator.ValidateContent(f, ext); err != nil {
		return nil, err
	}

	// check cheksum
	cks, size, err := s.validator.ComputeChecksum(f)
	if err != nil {
		return nil, err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("file: failed to reset file position: %w", err)
	}

	// compute uuid
	fileUUID := uuid.NewString()

	// calculate path
	realPath := helpers.BuildPath(dto.Folder, fileUUID, ext)

	// save file
	if _, err := s.storage.Save(ctx, realPath, f); err != nil {
		return nil, fmt.Errorf("file: failed to persist file: %w", err)
	}

	record := &Model{
		UUID:         fileUUID,
		DiskName:     fileUUID + ext,
		OriginalName: fileHeader.Filename,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		Extension:    ext,
		Size:         size,
		Checksum:     cks,
		Path:         realPath,
		Storage:      "local",
		CreatedAt:    time.Now(),
		CreatedBy:    &loggedID,
		IsPublic:     dto.Access,
	}

	if err := s.repo.Create(ctx, record); err != nil {
		_ = s.storage.Delete(ctx, realPath)
		return nil, fmt.Errorf("file: failed to save metadata: %w", err)
	}

	return toUploadResponse(record), nil
}

func (s *fileService) GetMetadata(ctx context.Context, fileUUID string) (*UploadResponseDTO, error) {
	record, err := s.repo.FindByUUID(ctx, fileUUID)
	if err != nil {
		return nil, err
	}

	return toUploadResponse(&record), nil
}

func (s *fileService) OpenForDownload(ctx context.Context, fileUUID string) (file interface{ Read([]byte) (int, error) }, meta *Model, err error) {
	record, err := s.repo.FindByUUID(ctx, fileUUID)
	if err != nil {
		return nil, nil, err
	}

	path := record.Path
	if record.IsArchived && record.ArchivedPath != nil {
		path = *record.ArchivedPath
	}

	reader, err := s.storage.Open(ctx, path)
	if err != nil {
		return nil, nil, err
	}

	return reader, &record, nil
}

func (s *fileService) DeleteFile(ctx context.Context, fileUUID string) error {
	usr, err := helpers.GetLoggedUserID(ctx)
	if err != nil {
		return errors.New("invalid claims")
	}

	record, err := s.repo.FindByUUID(ctx, fileUUID)
	if err != nil {
		return err
	}

	record.DeletedBy = &usr

	// Sesuai requirement: soft delete metadata dulu, penghapusan fisik
	// dilakukan belakangan oleh retention job — bukan di sini.
	num, err := s.repo.SoftDelete(ctx, record)
	if err != nil {
		return err
	}

	if num == 0 {
		return errors.New("file not found")
	}

	return nil
}

func (s *fileService) ReplaceWithTx(ctx context.Context, tx *gorm.DB, folder string, fh *multipart.FileHeader, oldFileUUID string) (newFile *UploadResponseDTO, err error) {
	usr, err := helpers.GetLoggedUserID(ctx)
	if err != nil {
		return nil, errors.New("invalid claims")
	}

	_, err = s.repo.FindByUUID(ctx, oldFileUUID)
	if err != nil {
		return nil, err
	}

	ext, err := s.validator.ValidateHeader(fh)
	if err != nil {
		return nil, err
	}

	f, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("file: failed to open uploaded file: %w", err)
	}
	defer f.Close()

	if err := s.validator.ValidateContent(f, ext); err != nil {
		return nil, err
	}

	checksum, size, err := s.validator.ComputeChecksum(f)
	if err != nil {
		return nil, err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("file: failed to reset file position: %w", err)
	}

	newUUID := uuid.NewString()
	relPath := helpers.BuildPath(folder, newUUID, ext)

	if _, err := s.storage.Save(ctx, relPath, f); err != nil {
		return nil, fmt.Errorf("file: failed to persist replacement file: %w", err)
	}

	record := &Model{
		UUID:         newUUID,
		DiskName:     newUUID + ext,
		OriginalName: fh.Filename,
		MimeType:     fh.Header.Get("Content-Type"),
		Extension:    ext,
		Size:         size,
		Checksum:     checksum,
		Path:         relPath,
		Storage:      "local",
		CreatedAt:    time.Now(),
		CreatedBy:    &usr,
	}

	// Insert metadata baru DI DALAM transaction caller — jika caller
	// rollback (mis. update FK gagal), record ini ikut hilang otomatis.
	txRepo := s.repo.WithTx(tx)
	if err := txRepo.Create(ctx, record); err != nil {
		_ = s.storage.Delete(ctx, relPath)
		return nil, fmt.Errorf("file: failed to save replacement metadata: %w", err)
	}

	return toUploadResponse(record), nil
}

func (s *fileService) FinalizeReplace(ctx context.Context, oldFileUUID string) error {
	oldFile, err := s.repo.FindByUUID(ctx, oldFileUUID)
	if err != nil {
		return err
	}

	if err := s.storage.Delete(ctx, oldFile.Path); err != nil {
		return fmt.Errorf("file: failed to delete old physical file: %w", err)
	}

	_, err = s.repo.SoftDelete(ctx, oldFile)
	if err != nil {
		return errors.New("file: failed to soft delete old physical file")
	}
	return nil
}
