package initialize

import (
	"context"
	"fmt"
	"sun-panel/api/api_v1"
	"sun-panel/global"
	"sun-panel/initialize/config"
	"sun-panel/initialize/database"
	"sun-panel/initialize/jwt"
	"sun-panel/initialize/lang"
	"sun-panel/initialize/runlog"
	"sun-panel/initialize/systemSettingCache"
	"sun-panel/lib/cmn"
	"sun-panel/lib/storage"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
)

func InitApp(configPath string) error {
	// 打印 logo
	Logo()
	gin.SetMode(global.RUNCODE) // GIN 运行模式

	// 日志
	logger, err := runlog.InitRunlog(global.RUNCODE, "running.log")
	if err != nil {
		return fmt.Errorf("log initialization error, %w", err)
	}

	global.Logger = logger
	// 配置初始化
	iniConfig, err := config.ConfigInit(configPath)
	if err != nil {
		return err
	}

	global.Config = iniConfig

	// 多语言初始化
	lang.LangInit("zh-cn") // en-us

	// 初始化数据库
	err = DatabaseConnect()
	if err != nil {
		return err
	}

	// 初始化存储系统
	storageInstance, err := InitStorage(configPath)
	if err != nil {
		return fmt.Errorf("storage initialization error: %w", err)
	}

	// 其他的初始化
	global.SystemSetting = systemSettingCache.InItSystemSettingCache()
	global.SystemMonitor = global.NewCache[any](5*time.Hour, -1, "systemMonitorCache")

	// 初始化JWT
	if err := jwt.InitJWT(); err != nil {
		return fmt.Errorf("JWT initialization error: %w", err)
	}

	// 初始化API组件
	api_v1.InitApiGroup(storageInstance)

	return nil
}

func DatabaseConnect() error {
	// 数据库连接 - 开始
	var dbClientInfo database.DbClient
	databaseDrive := global.Config.GetValueString("base", "database_drive")
	if databaseDrive == database.MYSQL {
		dbClientInfo = &database.MySQLConfig{
			Username:    global.Config.GetValueString("mysql", "username"),
			Password:    global.Config.GetValueString("mysql", "password"),
			Host:        global.Config.GetValueString("mysql", "host"),
			Port:        global.Config.GetValueString("mysql", "port"),
			Database:    global.Config.GetValueString("mysql", "db_name"),
			WaitTimeout: global.Config.GetValueInt("mysql", "wait_timeout"),
		}
	} else {
		dbClientInfo = &database.SQLiteConfig{
			Filename: global.Config.GetValueString("sqlite", "file_path"),
		}
	}
	db, err := database.DbInit(dbClientInfo)
	if err != nil {
		return fmt.Errorf("database DbInit error, %w", err)
	}

	global.Db = db
	models.Db = global.Db

	err = database.CreateDatabase(databaseDrive, global.Db)
	if err != nil {
		return fmt.Errorf("database CreateDatabase error, %w", err)
	}

	err = database.NotFoundAndCreateUser(global.Db)
	if err != nil {
		return fmt.Errorf("database NotFoundAndCreateUser error, %w", err)
	}

	return nil
}

// InitStorage initializes the storage system based on configuration
func InitStorage(configPath string) (*storage.RcloneStorage, error) {
	storageType := global.Config.GetValueString("rclone", "type")
	provider := global.Config.GetValueString("rclone", "provider")
	accessKey := global.Config.GetValueString("rclone", "access_key_id")
	secretKey := global.Config.GetValueString("rclone", "secret_access_key")
	endpoint := global.Config.GetValueString("rclone", "endpoint")
	bucket := global.Config.GetValueString("rclone", "bucket")
	rcloneConfig := &storage.RcloneConfig{
		Type:      storageType,
		Provider:  provider,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Endpoint:  endpoint,
		Bucket:    bucket,
	}
	// 使用带超时的上下文初始化存储
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rcloneStorage, err := storage.NewRcloneStorage(ctx, rcloneConfig, configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rclone storage: %w", err)
	}

	global.Logger.Info("Storage system initialized successfully")
	return rcloneStorage, nil
}

func Logo() {
	fmt.Println("     ____            ___                __")
	fmt.Println("    / __/_ _____    / _ \\___ ____  ___ / /")
	fmt.Println("   _\\ \\/ // / _ \\  / ___/ _ `/ _ \\/ -_) / ")
	fmt.Println("  /___/\\_,_/_//_/ /_/   \\_,_/_//_/\\__/_/  ")
	fmt.Println("")

	version := cmn.GetVersion()
	fmt.Println("Version:", version)
	fmt.Println("Welcome to the Sun-Panel.")
	fmt.Println("Project address:", "https://github.com/hslr-s/sun-panel")
}
