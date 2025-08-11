package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	MaxFileSize = 5 << 20 // 5MB
	UploadDir   = "uploads"
)

var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

type FileUploadResult struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
}

func InitUploadDirectories() error {
	directories := []string{
		filepath.Join(UploadDir, "news"),
		filepath.Join(UploadDir, "members"),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	return nil
}

func ValidateImageFile(file *multipart.FileHeader) error {

	if file.Size > MaxFileSize {
		return fmt.Errorf("file size too large. Maximum allowed size is %d bytes", MaxFileSize)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png"}
	isValidExt := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			isValidExt = true
			break
		}
	}

	if !isValidExt {
		return fmt.Errorf("invalid file extension. Allowed extensions: %v", allowedExts)
	}

	return nil
}

func SaveUploadedFile(c *gin.Context, file *multipart.FileHeader, subDir string) (*FileUploadResult, error) {

	if err := ValidateImageFile(file); err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("20060102150405")
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%s%s", timestamp, generateRandomString(8), ext)

	dirPath := filepath.Join(UploadDir, subDir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	fullPath := filepath.Join(dirPath, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	url := fmt.Sprintf("/uploads/%s/%s", subDir, filename)

	return &FileUploadResult{
		Filename: filename,
		Path:     fullPath,
		Size:     file.Size,
		URL:      url,
	}, nil
}

func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	if strings.HasPrefix(filePath, "/uploads/") {
		filePath = strings.TrimPrefix(filePath, "/")
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(filePath)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
