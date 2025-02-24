package storage

import (
	"context"
	"fmt"
	"io"
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
	config   *S3Config
	provider ProviderConfigurator
}

// NewS3Storage creates a new S3 storage instance
func NewS3Storage(ctx context.Context, config *S3Config) (*S3Storage, error) {
	global.Logger.Infof("Creating S3 storage with provider: %s, endpoint: %s, bucket: %s",
		config.Provider, config.Endpoint, config.Bucket)

	// 获取对应的provider配置器
	provider := GetProvider(config.Provider)

	// 创建凭证提供者
	creds := credentials.NewStaticCredentialsProvider(config.AccessKeyID, config.SecretAccessKey, "")

	// 调整区域格式
	region := provider.AdjustRegion(config.Region)

	// 创建基础S3客户端配置
	cfg := aws.Config{
		Credentials:      creds,
		Region:           region,
		RetryMaxAttempts: 3,
		RetryMode:        aws.RetryModeAdaptive,
	}

	// 添加基础S3客户端选项
	s3Options := []func(*s3.Options){
		func(o *s3.Options) {
			if config.DisableSSL {
				o.UseAccelerate = false
				o.UseDualstack = false
			}
		},
	}

	// 应用provider特定的配置
	cfg, s3Options = provider.ConfigureClient(&cfg, s3Options)

	// 创建S3客户端
	client := s3.NewFromConfig(cfg, s3Options...)

	storage := &S3Storage{
		client:   client,
		endpoint: config.Endpoint,
		bucket:   config.Bucket,
		config:   config,
		provider: provider,
	}

	// 测试连接
	if err := storage.testConnection(ctx); err != nil {
		global.Logger.Errorf("S3 connection test failed: %v", err)
		return nil, fmt.Errorf("S3 connection test failed: %w", err)
	}

	global.Logger.Info("S3 storage initialized successfully")
	return storage, nil
}

// testConnection verifies the S3 connection and bucket access
func (s *S3Storage) testConnection(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 使用HeadBucket测试连接
	_, err := s.client.HeadBucket(timeoutCtx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		global.Logger.Errorf("Failed to connect to bucket %s: %v", s.bucket, err)
		return fmt.Errorf("failed to connect to bucket: %w", err)
	}

	global.Logger.Infof("Successfully connected to bucket: %s", s.bucket)
	return nil
}

// Upload implements Storage.Upload for S3
func (s *S3Storage) Upload(ctx context.Context, reader io.Reader, fileName string) (string, error) {
	// 创建上传超时上下文
	timeoutDuration := time.Duration(s.config.TimeoutSeconds) * time.Second
	if timeoutDuration == 0 {
		timeoutDuration = 30 * time.Second
	}
	uploadCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	// 创建S3上传输入
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(fileName),
		Body:        reader,
		ContentType: aws.String(getContentType(fileName)),
	}

	global.Logger.Infof("Uploading file to S3 bucket %s with key %s", s.bucket, fileName)
	_, err := s.client.PutObject(uploadCtx, input)
	if err != nil {
		global.Logger.Errorf("Failed to upload to S3 (bucket: %s, key: %s, endpoint: %s): %v", s.bucket, fileName, s.endpoint, err)
		return "", fmt.Errorf("failed to upload to S3 (bucket: %s, endpoint: %s): %w", s.bucket, s.endpoint, err)
	}

	global.Logger.Infof("Successfully uploaded file to S3: %s/%s", s.bucket, fileName)

	// 使用provider生成URL
	return s.provider.GenerateURL(s.endpoint, s.bucket, fileName), nil
}

// Delete implements Storage.Delete for S3
func (s *S3Storage) Delete(ctx context.Context, path string) error {
	// 创建删除超时上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 从路径中提取键
	key := filepath.Base(path)

	// 创建删除输入
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(timeoutCtx, input)
	if err != nil {
		global.Logger.Errorf("Failed to delete from S3 (bucket: %s, key: %s): %v", s.bucket, key, err)
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	global.Logger.Infof("Successfully deleted file from S3: %s/%s", s.bucket, key)
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
