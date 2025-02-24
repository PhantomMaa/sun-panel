package initialize

import (
	"context"
	"errors"
	"fmt"
	"log"
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

var DB_DRIVER = database.SQLITE

// 全局存储实例
var globalStorage storage.Storage

func InitApp() error {
	// 打印 logo
	Logo()
	gin.SetMode(global.RUNCODE) // GIN 运行模式

	// 日志
	if logger, err := runlog.InitRunlog(global.RUNCODE, "running.log"); err != nil {
		log.Panicln("Log initialization error", err)
		panic(err)
	} else {
		global.Logger = logger
	}

	// 配置初始化
	config, err := config.ConfigInit()
	if err != nil {
		return err
	}

	global.Config = config

	// 多语言初始化
	lang.LangInit("zh-cn") // en-us

	DatabaseConnect()

	// 其他的初始化
	global.SystemSetting = systemSettingCache.InItSystemSettingCache()
	global.SystemMonitor = global.NewCache[interface{}](5*time.Hour, -1, "systemMonitorCache")

	// 初始化JWT
	if err := jwt.InitJWT(); err != nil {
		return fmt.Errorf("JWT initialization error: %w", err)
	}

	// 初始化存储系统
	storageInstance, err := InitStorage()
	if err != nil {
		return fmt.Errorf("storage initialization error: %w", err)
	}
	globalStorage = storageInstance

	// 初始化API组件
	api_v1.InitApiGroup(storageInstance)

	return nil
}

func DatabaseConnect() {
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

	if db, err := database.DbInit(dbClientInfo); err != nil {
		log.Panicln("Database initialization error", err)
		panic(err)
	} else {
		global.Db = db
		models.Db = global.Db
	}

	database.CreateDatabase(databaseDrive, global.Db)

	database.NotFoundAndCreateUser(global.Db)
}

// InitStorage initializes the storage system based on configuration
func InitStorage() (storage.Storage, error) {
	storageType := global.Config.GetValueString("base", "storage_drive")
	global.Logger.Infof("Initializing storage system with type: %s", storageType)

	var config storage.Config

	switch storageType {
	case "local":
		config = storage.Config{
			Type: storage.LocalStorageType,
		}
		global.Logger.Infof("Initializing local storage with path: %s", storage.LocalStorageBasePath)

	case "rclone":
		storageType := global.Config.GetValueString("rclone", "type")
		provider := global.Config.GetValueString("rclone", "provider")
		accessKey := global.Config.GetValueString("rclone", "access_key_id")
		secretKey := global.Config.GetValueString("rclone", "secret_access_key")
		endpoint := global.Config.GetValueString("rclone", "endpoint")
		bucket := global.Config.GetValueString("rclone", "bucket")

		config = storage.Config{
			Type: storage.RcloneStorageType,
			RcloneConfig: &storage.RcloneConfig{
				Type:      storageType,
				Provider:  provider,
				AccessKey: accessKey,
				SecretKey: secretKey,
				Endpoint:  endpoint,
				Bucket:    bucket,
			},
		}
	default:
		return nil, errors.New("invalid storage type : " + storageType)
	}

	// 使用带超时的上下文初始化存储
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	storageInstance, err := storage.NewStorage(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	global.Logger.Info("Storage system initialized successfully")
	return storageInstance, nil
}

// GetStorage returns the global storage instance
func GetStorage() storage.Storage {
	return globalStorage
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
