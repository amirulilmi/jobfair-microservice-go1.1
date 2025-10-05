package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AllowedImageExtensions adalah ekstensi file yang diperbolehkan
var AllowedImageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

// EnsureDir membuat direktori jika belum ada
func EnsureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

// SaveUploadedFile menyimpan file upload ke direktori tujuan
func SaveUploadedFile(file *multipart.FileHeader, destPath string) error {
	fmt.Printf("[UTILS] Opening uploaded file: %s\n", file.Filename)
	// Buka file upload
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Pastikan direktori ada
	dir := filepath.Dir(destPath)
	fmt.Printf("[UTILS] Ensuring directory exists: %s\n", dir)
	if err := EnsureDir(dir); err != nil {
		return err
	}

	// Buat file tujuan
	fmt.Printf("[UTILS] Creating destination file: %s\n", destPath)
	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file
	fmt.Printf("[UTILS] Copying file content...\n")
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	fmt.Printf("[UTILS] ✅ File saved successfully to: %s\n", destPath)
	return nil
}

// ValidateImageFile memvalidasi file image
func ValidateImageFile(file *multipart.FileHeader, maxSize int64) error {
	// Check file size
	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxSize)
	}

	// Check extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExt := false
	for _, allowedExt := range AllowedImageExtensions {
		if ext == allowedExt {
			validExt = true
			break
		}
	}

	if !validExt {
		return fmt.Errorf("invalid file type. Allowed: %v", AllowedImageExtensions)
	}

	return nil
}

// GenerateFileName menggenerate nama file yang unik
func GenerateFileName(userID uint, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d_%d%s", userID, timestamp, ext)
}

// GetUploadBasePath mendapatkan base path untuk uploads
func GetUploadBasePath() string {
	// Bisa dikonfigurasi via environment variable
	basePath := os.Getenv("UPLOAD_BASE_PATH")
	if basePath == "" {
		basePath = "./uploads"
	}
	return basePath
}

// DeleteFile menghapus file dari disk
func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	// Pastikan path adalah absolute atau relative dari project root
	if !filepath.IsAbs(filePath) {
		// Jika path dimulai dengan /, hapus untuk membuat relative path
		filePath = strings.TrimPrefix(filePath, "/")
	}

	// Check apakah file ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File sudah tidak ada, tidak perlu error
		return nil
	}

	// Hapus file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// FileExists mengecek apakah file ada
func FileExists(filePath string) bool {
	if !filepath.IsAbs(filePath) {
		filePath = strings.TrimPrefix(filePath, "/")
	}
	_, err := os.Stat(filePath)
	return err == nil
}
