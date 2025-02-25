package storage

// RcloneConfig holds configuration for rclone storage
type RcloneConfig struct {
	Type      string // local or s3
	Bucket    string
	Provider  string // s3 only
	AccessKey string // s3 only
	SecretKey string // s3 only
	Endpoint  string // s3 only
}
