package uploader

import (
	"context"
	"errors"
	"io"
)

// Sentinel errors — dicek dengan errors.Is() di layer atas (service/handler),
// mengikuti pola error handling project (lihat sloc.repository.go: gorm.ErrRecordNotFound).
var (
	ErrFileNotFound = errors.New("storage: file not found")
	ErrInvalidPath  = errors.New("storage: invalid or unsafe path")
)

// Storage adalah kontrak I/O fisik untuk File Manager Module.
// Implementasi konkret (LocalStorage, MinIOStorage, S3Storage) di-inject
// ke FileService lewat constructor — tidak ada logic bisnis di sini.
//
// Semua path yang diterima/dikembalikan adalah RELATIVE path
// (contoh: "tools/2026/07/uuid.jpg"), bukan absolute path di disk.
// Base path (LOCAL_STORAGE_PATH, bucket name, dst) adalah detail
// implementasi masing-masing backend.
type Storage interface {
	// Save menulis isi reader ke path yang ditentukan.
	// Directory akan dibuat otomatis jika belum ada.
	// Mengembalikan jumlah byte yang berhasil ditulis.
	Save(ctx context.Context, path string, reader io.Reader) (int64, error)

	// Open membuka file untuk dibaca (streaming download).
	// Caller WAJIB memanggil Close() pada hasilnya.
	Open(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete menghapus file fisik. Idempotent — tidak error jika file
	// sudah tidak ada (untuk mendukung retry pada retention job).
	Delete(ctx context.Context, path string) error

	// Move memindahkan file dari srcPath ke destPath.
	// Dipakai untuk proses Archive (path asli → path archive).
	Move(ctx context.Context, srcPath, destPath string) error

	// Exists mengecek apakah file ada di path tersebut.
	Exists(ctx context.Context, path string) (bool, error)
}
