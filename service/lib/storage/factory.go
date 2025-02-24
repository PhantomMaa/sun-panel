package storage

import (
	"context"
	"fmt"
	"sun-panel/global"
)

// StorageType represents the type of storage
type StorageType string

const (
	LocalStorageType  StorageType = "local"
	RcloneStorageType StorageType = "rclone"
)

// Config holds storage configuration
type Config struct {
	Type         StorageType
	RcloneConfig *RcloneConfig
}

// NewStorage creates a new storage instance based on configuration
func NewStorage(ctx context.Context, config Config) (Storage, error) {
	global.Logger.Infof("Creating storage instance with type: %s", config.Type)

	switch config.Type {
	case LocalStorageType:
		storage := NewLocalStorage()
		return storage, nil

	case RcloneStorageType:
		rcloneStorage, err := NewRcloneStorage(ctx, config.RcloneConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize rclone storage: %w", err)
		}
		return rcloneStorage, nil

	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}
}
