package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"sun-panel/global"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client   *s3.Client
	bucket   string
	endpoint string
}

// NewS3Storage creates a new S3 storage instance
func NewS3Storage(accessKeyID, secretAccessKey, endpoint, bucket, region string) (*S3Storage, error) {
	fmt.Printf("Creating S3 storage with endpoint: %s, bucket: %s\n", endpoint, bucket)

	// Create custom credentials provider
	creds := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")

	// Create custom endpoint resolver
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == "s3" {
			global.Logger.Infof("Creating S3 endpoint with URL: %s", endpoint)
			return aws.Endpoint{
				URL:               endpoint,
				HostnameImmutable: true,
				SigningRegion:     region,
				Source:            aws.EndpointSourceCustom,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	// Create S3 client config
	cfg := aws.Config{
		Credentials:                 creds,
		EndpointResolverWithOptions: customResolver,
		Region:                      region,
		// Add retry options
		RetryMaxAttempts: 3,
		// Force path style for compatibility with non-AWS S3 implementations
		DefaultsMode: aws.DefaultsModeInRegion,
	}

	// Add path style addressing for compatibility
	s3Options := []func(*s3.Options){
		func(o *s3.Options) {
			o.UsePathStyle = true
		},
	}

	// Create S3 client with options
	client := s3.NewFromConfig(cfg, s3Options...)

	storage := &S3Storage{
		client:   client,
		endpoint: endpoint,
		bucket:   bucket,
	}

	// Test connection
	if err := storage.testConnection(); err != nil {
		global.Logger.Errorf("S3 connection test failed: %v", err)
		return nil, fmt.Errorf("S3 connection test failed: %w", err)
	}

	global.Logger.Info("S3 storage initialized successfully")
	return storage, nil
}

// testConnection verifies the S3 connection and bucket access
func (s *S3Storage) testConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		global.Logger.Errorf("Failed to connect to S3 bucket %s: %v", s.bucket, err)
		return fmt.Errorf("failed to connect to S3 bucket: %w", err)
	}
	global.Logger.Infof("Successfully connected to S3 bucket: %s", s.bucket)
	return nil
}

// Upload implements Storage.Upload for S3
func (s *S3Storage) Upload(file *multipart.FileHeader, fileName string) (string, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Create the S3 upload input
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(fileName),
		Body:        src,
		ContentType: aws.String(getContentType(fileName)),
	}

	// Upload the file with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	global.Logger.Infof("Uploading file to S3 bucket %s with key %s", s.bucket, fileName)
	_, err = s.client.PutObject(ctx, input)
	if err != nil {
		global.Logger.Errorf("Failed to upload to S3 (bucket: %s, key: %s): %v", s.bucket, fileName, err)
		return "", fmt.Errorf("failed to upload to S3 (bucket: %s): %w", s.bucket, err)
	}
	global.Logger.Infof("Successfully uploaded file to S3: %s/%s", s.bucket, fileName)

	// Return the S3 URL
	return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, fileName), nil
}

// Delete implements Storage.Delete for S3
func (s *S3Storage) Delete(path string) error {
	// Extract the key from the full path
	key := filepath.Base(path)

	// Create the delete input
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	// Delete the object with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		global.Logger.Errorf("Failed to delete from S3 (bucket: %s, key: %s): %v", s.bucket, key, err)
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

// getContentType returns the content type based on file extension
func getContentType(fileName string) string {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}
