package storage

import (
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

// S3Config holds S3-specific configuration
type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	Bucket          string
	Region          string
}

// NewStorage creates a new storage instance based on configuration
func NewStorage(config Config) (Storage, error) {
	global.Logger.Infof("Creating storage instance with type: %s", config.Type)
	switch config.Type {
	case LocalStorageType:
		storage := NewLocalStorage()
		return storage, nil
	case S3StorageType:
		if config.S3Config == nil {
			return nil, fmt.Errorf("S3 configuration is required for S3 storage")
		}
		s3Storage, err := NewS3Storage(
			config.S3Config.AccessKeyID,
			config.S3Config.SecretAccessKey,
			config.S3Config.Endpoint,
			config.S3Config.Bucket,
			config.S3Config.Region,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize S3 storage: %w", err)
		}
		return s3Storage, nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}
}
