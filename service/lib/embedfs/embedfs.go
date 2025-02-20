// Package embedfs provides embedded file system functionality
package embedfs

import (
	"embed"
	"os"
	"path/filepath"
)

// FS represents an embedded file system
type FS struct {
	fs embed.FS
}

// New creates a new embedded file system instance
func New(fs embed.FS) *FS {
	return &FS{fs: fs}
}

// ReadFile reads a file from the embedded file system
func (e *FS) ReadFile(path string) ([]byte, error) {
	return e.fs.ReadFile("assets/" + path)
}

// ExtractFile extracts a file from the embedded file system to the specified path
func (e *FS) ExtractFile(embedPath, targetPath string) error {
	data, err := e.ReadFile(embedPath)
	if err != nil {
		return err
	}

	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(targetPath, data, 0644)
}

// Default is the default embedded file system instance
var Default *FS

// Init initializes the default embedded file system
func Init(fs embed.FS) {
	Default = New(fs)
}

// ReadEmbeddedFile reads a file using the default embedded file system
func ReadEmbeddedFile(path string) ([]byte, error) {
	return Default.ReadFile(path)
}

// ExtractEmbeddedFile extracts a file using the default embedded file system
func ExtractEmbeddedFile(embedPath, targetPath string) error {
	return Default.ExtractFile(embedPath, targetPath)
}
