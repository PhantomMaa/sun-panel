package storage

import (
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"sun-panel/global"

	"github.com/gin-gonic/gin"
)

const (
	// LocalStorageBasePath is the fixed base path for local storage
	LocalStorageBasePath = "uploads"
)

type LocalStorage struct{}

// Just a dummy gin context to pass to SaveUploadedFile
var _gin = &gin.Context{}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage() *LocalStorage {
	// Ensure base directory exists
	if err := os.MkdirAll(LocalStorageBasePath, os.ModePerm); err != nil {
		global.Logger.Errorf("Failed to create uploads directory: %v", err)
	}
	return &LocalStorage{}
}

// Upload implements Storage.Upload for local filesystem
func (s *LocalStorage) Upload(file *multipart.FileHeader, fileName string) (string, error) {
	// Construct file path
	filePath := filepath.Join(LocalStorageBasePath, fileName)

	// Save the file
	if err := _gin.SaveUploadedFile(file, filePath); err != nil {
		global.Logger.Errorf("Failed to save file %s: %v", filePath, err)
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	global.Logger.Infof("Successfully saved file to: %s", filePath)
	return path.Join("/", LocalStorageBasePath, fileName), nil
}

// Delete implements Storage.Delete for local filesystem
func (s *LocalStorage) Delete(path string) error {
	// Extract filename from path
	fileName := filepath.Base(path)

	// Construct the full path
	fullPath := filepath.Join(LocalStorageBasePath, fileName)

	global.Logger.Infof("Attempting to delete file: %s", fullPath)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			global.Logger.Warnf("File not found: %s", fullPath)
			return nil // Don't treat missing file as an error
		}

		global.Logger.Errorf("Failed to delete file %s: %v", fullPath, err)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	global.Logger.Infof("Successfully deleted file: %s", fullPath)
	return nil
}
