package initialize

import (
	"fmt"
	"sun-panel/global"
	"sun-panel/initialize/cUserToken"
	"sun-panel/initialize/config"
	"sun-panel/initialize/database"
	"sun-panel/initialize/lang"
	"sun-panel/initialize/other"
	"sun-panel/initialize/runlog"
	"sun-panel/initialize/systemSettingCache"
	"sun-panel/initialize/userToken"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"log"

	"github.com/gin-gonic/gin"

	"errors"
	"sun-panel/lib/storage"
)

var DB_DRIVER = database.SQLITE

func InitApp() error {
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

	// 初始化用户token
	global.UserToken = userToken.InitUserToken()
	global.CUserToken = cUserToken.InitCUserToken()

	// 其他的初始化
	global.VerifyCodeCachePool = other.InitVerifyCodeCachePool()
	global.SystemSetting = systemSettingCache.InItSystemSettingCache()
	global.SystemMonitor = global.NewCache[interface{}](5*time.Hour, -1, "systemMonitorCache")

	// 初始化存储系统
	if err := InitStorage(); err != nil {
		return fmt.Errorf("storage initialization error: %w", err)
	}

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
func InitStorage() error {
	storageType := global.Config.GetValueString("base", "storage_drive")
	global.Logger.Infof("Initializing storage system with type: %s", storageType)

	var config storage.Config

	switch storageType {
	case "local":
		config = storage.Config{
			Type: storage.LocalStorageType,
		}
		global.Logger.Infof("Initializing local storage with path: %s", storage.LocalStorageBasePath)

	case "s3":
		accessKeyID := global.Config.GetValueString("s3", "access_key_id")
		secretAccessKey := global.Config.GetValueString("s3", "secret_access_key")
		endpoint := global.Config.GetValueString("s3", "endpoint")
		bucket := global.Config.GetValueString("s3", "bucket")
		region := global.Config.GetValueString("s3", "region")

		if accessKeyID == "" || secretAccessKey == "" || endpoint == "" || bucket == "" || region == "" {
			return fmt.Errorf("missing required S3 configuration: accessKeyID=%v, secretAccessKey=%v, endpoint=%v, bucket=%v, region=%v",
				accessKeyID != "", secretAccessKey != "", endpoint != "", bucket != "", region != "")
		}

		global.Logger.Infof("Initializing S3 storage with endpoint: %s, bucket: %s, region: %s", endpoint, bucket, region)

		config = storage.Config{
			Type: storage.S3StorageType,
			S3Config: &storage.S3Config{
				AccessKeyID:     accessKeyID,
				SecretAccessKey: secretAccessKey,
				Endpoint:        endpoint,
				Bucket:          bucket,
				Region:          region,
			},
		}
	default:
		return errors.New("invalid storage type : " + storageType)
	}

	storageInstance, err := storage.NewStorage(config)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	storage.SetStorage(storageInstance)
	global.Logger.Info("Storage system initialized successfully")
	return nil
}

func Logo() {
	fmt.Println("     ____            ___                __")
	fmt.Println("    / __/_ _____    / _ \\___ ____  ___ / /")
	fmt.Println("   _\\ \\/ // / _ \\  / ___/ _ `/ _ \\/ -_) / ")
	fmt.Println("  /___/\\_,_/_//_/ /_/   \\_,_/_//_/\\__/_/  ")
	fmt.Println("")

	versionInfo := cmn.GetSysVersionInfo()
	fmt.Println("Version:", versionInfo.Version)
	fmt.Println("Welcome to the Sun-Panel.")
	fmt.Println("Project address:", "https://github.com/hslr-s/sun-panel")

}
