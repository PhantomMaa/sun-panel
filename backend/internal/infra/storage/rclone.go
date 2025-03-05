package storage

import (
	"context"
	"fmt"
	"io"
	"path"
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

	fmt.Printf("loading rclone config : %+v", section.KeysHash())

	for _, key := range section.Keys() {
		rconfig.FileSetValue("rclone", key.Name(), key.Value())
	}
	return nil
}

// creates a new rclone storage instance
func NewRcloneStorage(ctx context.Context, configPath, bucket string) (*RcloneStorage, error) {
	// 创建 rclone fs 实例
	err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	pathName := fmt.Sprintf("rclone:%s", bucket)
	f, err := fs.NewFs(ctx, pathName)
	if err != nil {
		return nil, fmt.Errorf("failed to create rclone fs: %w", err)
	}

	fmt.Println("Rclone storage initialized successfully")
	return &RcloneStorage{
		fs: f,
	}, nil
}

// implements Storage.Upload for rclone
func (r *RcloneStorage) Upload(ctx context.Context, reader io.Reader, fileName string) (string, error) {
	fmt.Printf("Uploading file via S3: %s", fileName)
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
	fmt.Printf("Successfully uploaded file: %s", url)
	return url, nil
}

// implements Storage.Delete for rclone
func (r *RcloneStorage) Delete(ctx context.Context, path string) error {
	fmt.Printf("Deleting file: %s", path)

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

	fmt.Printf("Successfully deleted file: %s", path)
	return nil
}

// implements file reading from storage
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
	defer func() {
		if err := reader.Close(); err != nil {
			_ = fmt.Errorf("failed to close reader. error : %v", err)
		}
	}()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	fmt.Printf("Successfully read file: %s", path)
	return data, nil
}
