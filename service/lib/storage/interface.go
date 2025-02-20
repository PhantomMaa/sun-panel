package storage

import (
	"mime/multipart"
)

// Storage defines the interface for file storage operations
type Storage interface {
	// Upload handles file upload and returns the file path/URL
	Upload(file *multipart.FileHeader, fileName string) (string, error)

	// Delete removes a file by its path
	Delete(path string) error
}

var storageInstance Storage

// SetStorage sets the global storage instance
func SetStorage(s Storage) {
	storageInstance = s
}

// GetStorage returns the global storage instance
func GetStorage() Storage {
	return storageInstance
}
