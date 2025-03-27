package config

import (
	"errors"
	"os"
	"sun-panel/internal/util"

	"gopkg.in/yaml.v3"
)

// AppConfig is the global application configuration
var AppConfig *Config

// Config represents the application configuration structure
type Config struct {
	Base   BaseConfig   `yaml:"base"`
	SQLite SQLiteConfig `yaml:"sqlite"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Rclone RcloneConfig `yaml:"rclone"`
	JWT    JWTConfig    `yaml:"jwt"`
	OAuth  OAuthConfig  `yaml:"oauth"`
}

// BaseConfig represents the base section configuration
type BaseConfig struct {
	HTTPPort           string `yaml:"http_port"`
	RootURL            string `yaml:"root_url"`
	DatabaseDrive      string `yaml:"database_drive"`
	EnableStaticServer bool   `yaml:"enable_static_server"`
	EnableMonitor      bool   `yaml:"enable_monitor"`
	URLPrefix          string `yaml:"url_prefix"`
}

// SQLiteConfig represents the sqlite section configuration
type SQLiteConfig struct {
	FilePath string `yaml:"file_path"`
}

// MySQLConfig represents the mysql section configuration
type MySQLConfig struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	DBName      string `yaml:"db_name"`
	WaitTimeout int    `yaml:"wait_timeout"`
}

// RcloneConfig represents the rclone section configuration
type RcloneConfig struct {
	Type   string `yaml:"type"`
	Bucket string `yaml:"bucket"`
}

// JWTConfig represents the jwt section configuration
type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}

// OAuthConfig represents the oauth section configuration
type OAuthConfig struct {
	Enable    bool                  `yaml:"enable"`
	Providers []OAuthProviderConfig `yaml:"providers"`
}

// OAuthProviderConfig represents the configuration for an OAuth provider
type OAuthProviderConfig struct {
	Name                   string `yaml:"name"`
	ClientID               string `yaml:"client_id"`
	ClientSecret           string `yaml:"client_secret"`
	AuthURL                string `yaml:"auth_url"`
	TokenURL               string `yaml:"token_url"`
	UserInfoURL            string `yaml:"user_info_url"`
	Scopes                 string `yaml:"scopes"`
	FieldMappingIdentifier  string `yaml:"field_mapping_identifier"`
	FieldMappingDisplayName string `yaml:"field_mapping_display_name"`
	FieldMappingEmail       string `yaml:"field_mapping_email"`
}

// Init initializes the configuration
func Init(configPath string) (*Config, error) {
	exists, err := util.PathExists(configPath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("配置文件不存在: " + configPath)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Create config object
	config := new(Config)

	// Parse YAML
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	// Set global config
	AppConfig = config

	return config, nil
}
