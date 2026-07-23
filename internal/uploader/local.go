package uploader

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	basePath string
}

// NewLocalStorage membaca root path dari env LOCAL_STORAGE_PATH.
// Tidak ada default hardcoded — kalau env kosong, gagal cepat (fail fast)
// saat startup, bukan saat runtime pertama kali ada yang upload.
func NewLocalStorage() (Storage, error) {
	base := os.Getenv("LOCAL_STORAGE_PATH")
	if base == "" {
		return nil, fmt.Errorf("storage: LOCAL_STORAGE_PATH env var is not set")
	}

	base, err := filepath.Abs(base)
	if err != nil {
		return nil, fmt.Errorf("storage: invalid LOCAL_STORAGE_PATH: %w", err)
	}

	if err := os.MkdirAll(base, 0o755); err != nil {
		return nil, fmt.Errorf("storage: failed to prepare base directory: %w", err)
	}

	return &LocalStorage{basePath: base}, nil
}

// resolvePath membangun absolute path dari relative path, sekaligus
// mem-block path traversal (mis. "../../etc/passwd"). Defensive check
// ini penting meski path biasanya dibentuk server-side dari UUID —
// jangan pernah percaya input mentah tanpa validasi.
func (s *LocalStorage) resolvePath(relPath string) (string, error) {
	cleaned := filepath.Clean("/" + relPath)
	full := filepath.Join(s.basePath, cleaned)

	if !strings.HasPrefix(full, s.basePath+string(os.PathSeparator)) && full != s.basePath {
		return "", ErrInvalidPath
	}
	return full, nil
}

func (s *LocalStorage) Save(ctx context.Context, path string, reader io.Reader) (int64, error) {
	full, err := s.resolvePath(path)
	if err != nil {
		return 0, err
	}

	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return 0, fmt.Errorf("storage: failed to create directory: %w", err)
	}

	f, err := os.Create(full)
	if err != nil {
		return 0, fmt.Errorf("storage: failed to create file: %w", err)
	}
	defer f.Close()

	written, err := io.Copy(f, reader)
	if err != nil {
		return written, fmt.Errorf("storage: failed to write file: %w", err)
	}

	return written, nil
}

func (s *LocalStorage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	full, err := s.resolvePath(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(full)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, fmt.Errorf("storage: failed to open file: %w", err)
	}
	return f, nil
}

func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	full, err := s.resolvePath(path)
	if err != nil {
		return err
	}

	if err := os.Remove(full); err != nil {
		if os.IsNotExist(err) {
			return nil // idempotent
		}
		return fmt.Errorf("storage: failed to delete file: %w", err)
	}
	return nil
}

func (s *LocalStorage) Move(ctx context.Context, srcPath, destPath string) error {
	src, err := s.resolvePath(srcPath)
	if err != nil {
		return err
	}
	dest, err := s.resolvePath(destPath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("storage: failed to create archive directory: %w", err)
	}

	if err := os.Rename(src, dest); err != nil {
		return fmt.Errorf("storage: failed to move file: %w", err)
	}
	return nil
}

func (s *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
	full, err := s.resolvePath(path)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(full)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("storage: failed to stat file: %w", err)
}
