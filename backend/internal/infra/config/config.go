package config

import (
	"errors"
	"sun-panel/internal/util"

	"gopkg.in/ini.v1"
)

// AppConfig is the global application configuration
var AppConfig *Config

// Config represents the application configuration structure
type Config struct {
	Base   BaseConfig   `ini:"base"`
	SQLite SQLiteConfig `ini:"sqlite"`
	MySQL  MySQLConfig  `ini:"mysql"`
	Rclone RcloneConfig `ini:"rclone"`
	JWT    JWTConfig    `ini:"jwt"`
	OAuth  OAuthConfig  `ini:"oauth"`
}

// BaseConfig represents the base section configuration
type BaseConfig struct {
	HTTPPort           string `ini:"http_port"`
	RootURL            string `ini:"root_url"`
	DatabaseDrive      string `ini:"database_drive"`
	EnableStaticServer bool   `ini:"enable_static_server"`
	EnableMonitor      bool   `ini:"enable_monitor"`
	URLPrefix          string `ini:"url_prefix"`
}

// SQLiteConfig represents the sqlite section configuration
type SQLiteConfig struct {
	FilePath string `ini:"file_path"`
}

// MySQLConfig represents the mysql section configuration
type MySQLConfig struct {
	Host        string `ini:"host"`
	Port        string `ini:"port"`
	Username    string `ini:"username"`
	Password    string `ini:"password"`
	DBName      string `ini:"db_name"`
	WaitTimeout int    `ini:"wait_timeout"`
}

// RcloneConfig represents the rclone section configuration
type RcloneConfig struct {
	Type   string `ini:"type"`
	Bucket string `ini:"bucket"`
}

// JWTConfig represents the jwt section configuration
type JWTConfig struct {
	Secret string `ini:"secret"`
	Expire int    `ini:"expire"`
}

// OAuthConfig represents the oauth section configuration
type OAuthConfig struct {
	Enable bool                `ini:"enable"`
	GitHub OAuthProviderConfig `ini:"-"`
	Google OAuthProviderConfig `ini:"-"`
}

// OAuthProviderConfig represents the configuration for an OAuth provider
type OAuthProviderConfig struct {
	ClientID                string `ini:"client_id"`
	ClientSecret            string `ini:"client_secret"`
	AuthURL                 string `ini:"auth_url"`
	TokenURL                string `ini:"token_url"`
	UserInfoURL             string `ini:"user_info_url"`
	Scopes                  string `ini:"scopes"`
	FieldMappingIdentifier  string `ini:"field_mapping_identifier"`
	FieldMappingDisplayName string `ini:"field_mapping_display_name"`
	FieldMappingEmail       string `ini:"field_mapping_email"`
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

	// Load INI file
	iniFile, err := ini.Load(configPath)
	if err != nil {
		return nil, err
	}

	// Create config object
	config := new(Config)

	// Map INI to struct
	err = iniFile.MapTo(config)
	if err != nil {
		return nil, err
	}

	// 手动处理嵌套节点
	githubSection := iniFile.Section("oauth.github")
	if githubSection != nil {
		err = githubSection.MapTo(&config.OAuth.GitHub)
		if err != nil {
			return nil, err
		}
	}

	googleSection := iniFile.Section("oauth.google")
	if googleSection != nil {
		err = googleSection.MapTo(&config.OAuth.Google)
		if err != nil {
			return nil, err
		}
	}

	// Set global config
	AppConfig = config

	return config, nil
}
