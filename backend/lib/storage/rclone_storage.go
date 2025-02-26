package storage

import (
	"context"
	"fmt"
	"io"
	"path"
	"sun-panel/global"
	"time"

	"github.com/rclone/rclone/fs"
	rconfig "github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/operations"
	"gopkg.in/ini.v1"

	_ "github.com/rclone/rclone/backend/local"
	_ "github.com/rclone/rclone/backend/s3"
)

type RcloneStorage struct {
	fs fs.Fs
}

func loadConfig(configPath string) error {
	configFile, err := ini.Load(configPath)
	if err != nil {
		return fmt.Errorf("ini.Load error : %w", err)
	}

	section := configFile.Section("rclone")
	for _, key := range section.Keys() {
		rconfig.FileSetValue("rclone", key.Name(), key.Value())
	}
	return nil
}

// NewRcloneStorage creates a new rclone storage instance
func NewRcloneStorage(ctx context.Context, config *RcloneConfig, configPath string) (*RcloneStorage, error) {
	global.Logger.Infof("Creating rclone storage with config: %+v", config)

	// 创建 rclone fs 实例
	err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	path := fmt.Sprintf("rclone:%s", config.Bucket)
	f, err := fs.NewFs(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to create rclone fs: %w", err)
	}

	global.Logger.Info("Rclone storage initialized successfully")
	return &RcloneStorage{
		fs: f,
	}, nil
}

// Upload implements Storage.Upload for rclone
func (r *RcloneStorage) Upload(ctx context.Context, reader io.Reader, fileName string) (string, error) {
	global.Logger.Infof("Uploading file via S3: %s", fileName)
	// 一定要检查 bucket 是否存在，避免 BucketAlreadyExists
	_, err := r.fs.List(ctx, "")
	if err != nil {
		return "", fmt.Errorf("bucket check failed: %w", err)
	}

	// 将 io.Reader 转换为 io.ReadCloser
	readCloser := io.NopCloser(reader)

	// Create a new object in the bucket
	obj, err := operations.RcatSize(ctx, r.fs, fileName, readCloser, -1, time.Now(), fs.Metadata{})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate URL
	url := path.Join(obj.Remote())
	global.Logger.Infof("Successfully uploaded file: %s", url)
	return url, nil
}

// Delete implements Storage.Delete for rclone
func (r *RcloneStorage) Delete(ctx context.Context, path string) error {
	global.Logger.Infof("Deleting file: %s", path)

	// Find the object
	obj, err := r.fs.NewObject(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to find object: %w", err)
	}

	// Delete the object
	err = obj.Remove(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	global.Logger.Infof("Successfully deleted file: %s", path)
	return nil
}

// Get implements file reading from storage
func (r *RcloneStorage) Get(ctx context.Context, path string) ([]byte, error) {
	// Find the object
	obj, err := r.fs.NewObject(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to find object: %w", err)
	}

	// Read the object
	reader, err := obj.Open(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open object: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	global.Logger.Infof("Successfully read file: %s", path)
	return data, nil
}
