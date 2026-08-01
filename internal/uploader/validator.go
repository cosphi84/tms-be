package uploader

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

var (
	ErrUnsupportedMimeType = errors.New("validator: unsupported file type")
	ErrFileTooLarge        = errors.New("validator: file exceeds maximum allowed size")
	ErrEmptyFile           = errors.New("validator: file is empty")
	ErrMimeMismatch        = errors.New("validator: declared content-type does not match actual file content")
)

const MaxFileSizeBytes int64 = 5 * 1024 * 1024 // 5 MB

// sniffBufferSize — 512 bytes cukup untuk http.DetectContentType
// mengenali magic number JPEG/PNG/WEBP secara akurat.
const sniffBufferSize = 512

var allowedMimeTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type FileValidator interface {
	// ValidateHeader mengecek Content-Type klaim client & size, fail fast
	// tanpa I/O. Dipanggil pertama, sebelum file dibuka.
	ValidateHeader(fh *multipart.FileHeader) (extension string, err error)

	// ValidateContent membaca magic bytes ASLI dari file, memverifikasi
	// bahwa isi file benar-benar sesuai tipe yang diklaim di header.
	// PENTING: fungsi ini melakukan Seek(0) di awal DAN di akhir,
	// sehingga aman dipanggil sebelum ComputeChecksum/Storage.Save
	// tanpa perlu caller mengurus posisi reader secara manual.
	ValidateContent(file multipart.File, declaredExt string) error

	// ComputeChecksum menghitung SHA256 sambil streaming.
	// PENTING: caller wajib Seek(0, io.SeekStart) SETELAH memanggil ini,
	// karena io.Copy menghabiskan reader hingga EOF.
	ComputeChecksum(file multipart.File) (checksum string, size uint64, err error)
}

type fileValidator struct {
	maxSize      int64
	allowedTypes map[string]string
}

func NewFileValidator() FileValidator {
	return &fileValidator{
		maxSize:      MaxFileSizeBytes,
		allowedTypes: allowedMimeTypes,
	}
}

func (v *fileValidator) ValidateHeader(fh *multipart.FileHeader) (string, error) {
	if fh.Size <= 0 {
		return "", ErrEmptyFile
	}
	if fh.Size > v.maxSize {
		return "", fmt.Errorf("%w: %d bytes (max %d bytes)", ErrFileTooLarge, fh.Size, v.maxSize)
	}

	mimeType := fh.Header.Get("Content-Type")
	ext, ok := v.allowedTypes[mimeType]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedMimeType, mimeType)
	}

	return ext, nil
}

// ValidateContent adalah lapisan pertahanan kedua: header Content-Type
// bisa dipalsukan klien (curl -H "Content-Type: image/jpeg" --data-binary @malware.exe),
// tapi magic bytes di awal file tidak bisa dibohongi tanpa merusak filenya sendiri.
func (v *fileValidator) ValidateContent(file multipart.File, declaredExt string) error {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("validator: failed to seek file: %w", err)
	}

	buf := make([]byte, sniffBufferSize)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("validator: failed to read file for sniffing: %w", err)
	}
	buf = buf[:n]

	detectedMime := http.DetectContentType(buf)
	detectedExt, isKnown := v.allowedTypes[detectedMime]
	if !isKnown {
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("validator: failed to reset file position: %w", err)
		}
		return fmt.Errorf("%w: actual content is %s", ErrUnsupportedMimeType, detectedMime)
	}

	if detectedExt != declaredExt {
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("validator: failed to reset file position: %w", err)
		}
		return fmt.Errorf("%w: declared %s, detected %s", ErrMimeMismatch, declaredExt, detectedExt)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("validator: failed to reset file position: %w", err)
	}

	return nil
}

func (v *fileValidator) ComputeChecksum(file multipart.File) (string, uint64, error) {
	hasher := sha256.New()
	size, err := io.Copy(hasher, file)
	if err != nil {
		return "", 0, fmt.Errorf("validator: failed to compute checksum: %w", err)
	}
	return hex.EncodeToString(hasher.Sum(nil)), uint64(size), nil
}

// bytesReader kecil sebagai helper internal jika suatu saat butuh
// re-validate dari buffer tanpa reader fisik — disiapkan untuk symmetry,
// belum dipakai di alur utama saat ini.
var _ = bytes.NewReader
