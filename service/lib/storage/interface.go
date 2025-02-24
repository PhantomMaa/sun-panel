package storage

import (
	"context"
	"io"
)

// Storage defines the interface for file storage operations
type Storage interface {
	// Upload handles file upload and returns the file path/URL
	// ctx provides context for cancellation and timeouts
	// reader provides the file content
	// fileName is the name to save the file as
	Upload(ctx context.Context, reader io.Reader, fileName string) (string, error)

	// Delete removes a file by its path
	// ctx provides context for cancellation and timeouts
	Delete(ctx context.Context, path string) error

	// Get reads a file by its path
	// ctx provides context for cancellation and timeouts
	Get(ctx context.Context, path string) ([]byte, error)
}
