package storage

import (
	"context"
	"fmt"
	"sun-panel/global"
)

// StorageType represents the type of storage
type StorageType string

const (
	LocalStorageType StorageType = "local"
	S3StorageType    StorageType = "s3"
)

// Config holds storage configuration
type Config struct {
	Type     StorageType
	S3Config *S3Config
}

// NewStorage creates a new storage instance based on configuration
func NewStorage(ctx context.Context, config Config) (Storage, error) {
	global.Logger.Infof("Creating storage instance with type: %s", config.Type)

	switch config.Type {
	case LocalStorageType:
		storage := NewLocalStorage()
		return storage, nil

	case S3StorageType:
		if config.S3Config == nil {
			return nil, fmt.Errorf("S3 configuration is required for S3 storage")
		}

		// 设置默认超时
		if config.S3Config.TimeoutSeconds == 0 {
			config.S3Config.TimeoutSeconds = 30
		}

		s3Storage, err := NewS3Storage(ctx, config.S3Config)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize S3 storage: %w", err)
		}
		return s3Storage, nil

	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}
}
