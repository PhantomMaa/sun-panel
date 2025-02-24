package storage

// S3Config holds S3-specific configuration
type S3Config struct {
	Provider        Provider // aws, aliyun, minio
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	Bucket          string
	Region          string
	// 新增SSL配置选项
	DisableSSL bool
	// 新增自定义超时设置
	TimeoutSeconds int
}
