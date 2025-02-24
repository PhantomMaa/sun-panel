package storage

// RcloneConfig holds configuration for rclone storage
type RcloneConfig struct {
	Type      string // storage type (s3, oss, etc)
	Provider  string // specific provider (Alibaba, AWS, etc)
	AccessKey string
	SecretKey string
	Endpoint  string
	Bucket    string
}
