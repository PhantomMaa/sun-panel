package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sun-panel/global"
)

const (
	// LocalStorageBasePath is the fixed base path for local storage
	LocalStorageBasePath = "uploads"
)

type LocalStorage struct{}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage() *LocalStorage {
	// Ensure base directory exists
	if err := os.MkdirAll(LocalStorageBasePath, os.ModePerm); err != nil {
		global.Logger.Errorf("Failed to create uploads directory: %v", err)
	}
	return &LocalStorage{}
}

// sanitizeFileName removes potentially dangerous characters from the filename
func sanitizeFileName(fileName string) string {
	// 移除路径分隔符和潜在的危险字符
	fileName = filepath.Base(fileName)
	fileName = strings.ReplaceAll(fileName, "..", "")
	return fileName
}

// Upload implements Storage.Upload for local filesystem
func (s *LocalStorage) Upload(ctx context.Context, reader io.Reader, fileName string) (string, error) {
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("upload cancelled: %w", ctx.Err())
	default:
	}

	// 清理文件名
	fileName = sanitizeFileName(fileName)
	if fileName == "" {
		return "", fmt.Errorf("invalid filename")
	}

	// 构造文件路径
	filePath := filepath.Join(LocalStorageBasePath, fileName)

	// 创建目标文件
	file, err := os.Create(filePath)
	if err != nil {
		global.Logger.Errorf("Failed to create file %s: %v", filePath, err)
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 复制文件内容
	if _, err := io.Copy(file, reader); err != nil {
		global.Logger.Errorf("Failed to write file %s: %v", filePath, err)
		os.Remove(filePath) // 清理失败的文件
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	global.Logger.Infof("Successfully saved file to: %s", filePath)
	return path.Join("/", LocalStorageBasePath, fileName), nil
}

// Delete implements Storage.Delete for local filesystem
func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("delete cancelled: %w", ctx.Err())
	default:
	}

	// 提取并清理文件名
	fileName := sanitizeFileName(filepath.Base(path))
	if fileName == "" {
		return fmt.Errorf("invalid path")
	}

	// 构造完整路径
	fullPath := filepath.Join(LocalStorageBasePath, fileName)

	global.Logger.Infof("Attempting to delete file: %s", fullPath)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			global.Logger.Warnf("File not found: %s", fullPath)
			return nil // 不将文件不存在视为错误
		}

		global.Logger.Errorf("Failed to delete file %s: %v", fullPath, err)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	global.Logger.Infof("Successfully deleted file: %s", fullPath)
	return nil
}
