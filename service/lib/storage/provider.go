package storage

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Provider represents the S3 storage provider type
type Provider string

const (
	ProviderAWS    Provider = "aws"
	ProviderAliyun Provider = "aliyun"
	ProviderMinIO  Provider = "minio"
)

// ProviderConfigurator defines the interface for provider-specific configurations
type ProviderConfigurator interface {
	// ConfigureClient configures the AWS config and S3 options for the specific provider
	ConfigureClient(cfg *aws.Config, opts []func(*s3.Options)) (aws.Config, []func(*s3.Options))

	// GenerateURL generates the appropriate URL format for the provider
	GenerateURL(endpoint, bucket, key string) string

	// AdjustRegion adjusts the region format if needed
	AdjustRegion(region string) string
}

// GetProvider returns the appropriate provider configurator
func GetProvider(provider Provider) ProviderConfigurator {
	switch provider {
	case ProviderAWS:
		return &AWSProvider{}
	case ProviderAliyun:
		return &AliyunProvider{}
	case ProviderMinIO:
		return &MinIOProvider{}
	default:
		return &AWSProvider{} // Default to AWS provider
	}
}

// AWSProvider implements ProviderConfigurator for AWS S3
type AWSProvider struct{}

func (p *AWSProvider) ConfigureClient(cfg *aws.Config, opts []func(*s3.Options)) (aws.Config, []func(*s3.Options)) {
	return *cfg, opts
}

func (p *AWSProvider) GenerateURL(endpoint, bucket, key string) string {
	// Extract region from endpoint
	region := "us-east-1" // default region
	if parts := strings.Split(endpoint, "."); len(parts) > 2 {
		region = parts[1]
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
}

func (p *AWSProvider) AdjustRegion(region string) string {
	return region
}

// AliyunProvider implements ProviderConfigurator for Aliyun OSS
type AliyunProvider struct{}

func (p *AliyunProvider) ConfigureClient(cfg *aws.Config, opts []func(*s3.Options)) (aws.Config, []func(*s3.Options)) {
	// 阿里云OSS特定配置
	opts = append(opts, func(o *s3.Options) {
		o.UsePathStyle = false
		o.UseAccelerate = false
		o.UseDualstack = false
	})

	// 设置阿里云OSS端点
	endpoint := fmt.Sprintf("https://oss-%s.aliyuncs.com", strings.TrimPrefix(cfg.Region, "oss-"))

	// 添加自定义解析器
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == "s3" {
			return aws.Endpoint{
				URL:               endpoint,
				HostnameImmutable: true,
				SigningRegion:     region,
				Source:            aws.EndpointSourceCustom,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
	cfg.EndpointResolverWithOptions = customResolver

	// 设置基础端点
	cfg.BaseEndpoint = aws.String(endpoint)

	// 设置区域为原始区域（不带oss-前缀）
	cfg.Region = strings.TrimPrefix(cfg.Region, "oss-")

	return *cfg, opts
}

func (p *AliyunProvider) GenerateURL(endpoint, bucket, key string) string {
	// 移除https://前缀
	cleanEndpoint := strings.TrimPrefix(endpoint, "https://")
	return fmt.Sprintf("https://%s.%s/%s", bucket, cleanEndpoint, key)
}

func (p *AliyunProvider) AdjustRegion(region string) string {
	if !strings.HasPrefix(region, "oss-") {
		return "oss-" + region
	}
	return region
}

// MinIOProvider implements ProviderConfigurator for MinIO
type MinIOProvider struct{}

func (p *MinIOProvider) ConfigureClient(cfg *aws.Config, opts []func(*s3.Options)) (aws.Config, []func(*s3.Options)) {
	// Force path style for MinIO
	opts = append(opts, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	// Add MinIO specific options
	cfg.DefaultsMode = aws.DefaultsModeLegacy

	return *cfg, opts
}

func (p *MinIOProvider) GenerateURL(endpoint, bucket, key string) string {
	return fmt.Sprintf("%s/%s/%s", endpoint, bucket, key)
}

func (p *MinIOProvider) AdjustRegion(region string) string {
	return region
}
