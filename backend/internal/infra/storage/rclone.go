package storage

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rclone/rclone/fs"
	rconfig "github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/operations"
	"gopkg.in/yaml.v3"

	_ "github.com/rclone/rclone/backend/local"
	_ "github.com/rclone/rclone/backend/s3"
)

type RcloneStorage struct {
	fs fs.Fs
}

// RcloneConfig represents the rclone section in the YAML config
type RcloneConfig struct {
	Bucket string `yaml:"bucket"`
	Conf   string `yaml:"conf"`
}

// Config represents the root configuration structure
type Config struct {
	Rclone RcloneConfig `yaml:"rclone"`
}

func loadConfig(configPath string) error {
	// Read the YAML file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse the YAML
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	fmt.Printf("loading rclone config from conf field\n")

	// Parse the INI formatted config from the conf field
	if config.Rclone.Conf == "" {
		return fmt.Errorf("rclone.conf is empty")
	}

	// Parse the config directly
	scanner := bufio.NewScanner(strings.NewReader(config.Rclone.Conf))
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		
		// Skip section headers
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			continue
		}
		
		// Parse key=value pairs
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				rconfig.FileSetValue("rclone", key, value)
			}
		}
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning config: %w", err)
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
	//global.Logger.Infof("Uploading file via rclone, file : %s", fileName)
	// Check if bucket exists, if not, create it
	_, err := r.fs.List(ctx, "")
	if err != nil {
		fmt.Printf("Bucket does not exist or cannot be accessed: %v, attempting to create it\n", err)

		// Try to create the bucket/directory
		err = operations.Mkdir(ctx, r.fs, "")
		if err != nil {
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}

		fmt.Println("Successfully created bucket")
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
func (r *RcloneStorage) Delete(ctx context.Context, filepath string) error {
	fmt.Printf("Deleting file: %s", filepath)

	// Find the object
	obj, err := r.fs.NewObject(ctx, filepath)
	if err != nil {
		return fmt.Errorf("failed to find object: %w", err)
	}

	// Delete the object
	err = obj.Remove(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	fmt.Printf("Successfully deleted file: %s", filepath)
	return nil
}

// implements file reading from storage
func (r *RcloneStorage) Get(ctx context.Context, filepath string) ([]byte, error) {
	// Find the object
	obj, err := r.fs.NewObject(ctx, filepath)
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

	fmt.Printf("Successfully read file: %s", filepath)
	return data, nil
}
