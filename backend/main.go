package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"sun-panel/api"
	"sun-panel/internal/cache"
	"sun-panel/internal/common"
	"sun-panel/internal/common/systemSetting"
	"sun-panel/internal/database"
	"sun-panel/internal/global"
	"sun-panel/internal/iniConfig"
	"sun-panel/internal/jwt"
	"sun-panel/internal/language"
	"sun-panel/internal/repository"
	"sun-panel/internal/router"
	"sun-panel/internal/storage"
	"sun-panel/internal/zapLog"
	"time"
)

func main() {
	// Parse command line arguments
	configPath := flag.String("c", "conf.ini", "Path to configuration file")
	flag.Parse()

	// Initialize the application with the specified config path
	err := InitApp(*configPath)
	if err != nil {
		log.Panicln("初始化错误:", err)
	}

	httpPort := global.Config.GetValueString("base", "http_port")
	if err := router.InitRouters(":" + httpPort); err != nil {
		panic(err)
	}
}

func InitApp(configPath string) error {
	// 打印 logo
	Logo()
	gin.SetMode(global.RUNCODE) // GIN 运行模式

	// 日志
	logger, err := zapLog.InitLog(global.RUNCODE, "running.log")
	if err != nil {
		return fmt.Errorf("log initialization error, %w", err)
	}

	global.Logger = logger
	// 配置初始化
	global.Config, err = iniConfig.ConfigInit(configPath)
	if err != nil {
		return err
	}

	// 多语言初始化
	language.LangInit("zh-cn") // en-us

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
	global.SystemSetting = &systemSetting.SystemSettingCache{
		Cache: cache.NewGoCache[any](5*time.Hour, -1),
	}
	global.SystemMonitor = cache.NewGoCache[any](5*time.Hour, -1)

	// 初始化JWT
	if err := jwt.InitJWT(); err != nil {
		return fmt.Errorf("JWT initialization error: %w", err)
	}

	// 初始化API组件
	api.InitApiGroup(storageInstance)

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
	repository.Db = global.Db

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
	// 使用带超时的上下文初始化存储
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rcloneStorage, err := storage.NewRcloneStorage(ctx, configPath)
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

	version := common.GetVersion()
	fmt.Println("Version:", version)
	fmt.Println("Welcome to the Sun-Panel.")
	fmt.Println("Project address:", "https://github.com/hslr-s/sun-panel")
}
