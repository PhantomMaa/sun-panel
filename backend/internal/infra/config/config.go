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
}

// BaseConfig represents the base section configuration
type BaseConfig struct {
	HTTPPort           string `ini:"http_port"`
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

	// Set global config
	AppConfig = config

	return config, nil
}
