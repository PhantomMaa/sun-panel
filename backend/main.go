package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/biz/service"
	"sun-panel/internal/global"
	"sun-panel/internal/infra/config"
	"sun-panel/internal/infra/database"
	"sun-panel/internal/infra/storage"
	"sun-panel/internal/infra/zaplog"
	"sun-panel/internal/util/i18n"
	"sun-panel/internal/util/jwt"
	"sun-panel/internal/web/router"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command line arguments
	configPath := flag.String("c", "conf.yaml", "Path to configuration file")
	flag.Parse()

	// Initialize the application with the specified config path
	err := InitApp(*configPath)
	if err != nil {
		log.Panicln("初始化错误:", err)
	}
}

func InitApp(configPath string) error {
	// 打印 logo
	Logo()
	gin.SetMode(global.RUNCODE) // GIN 运行模式

	// 日志
	err := zaplog.InitLog(global.RUNCODE, "running.zaplog")
	if err != nil {
		return fmt.Errorf("zaplog initialization error, %w", err)
	}

	// 配置初始化
	global.Config, err = config.Init(configPath)
	if err != nil {
		return err
	}

	// 多语言初始化
	i18n.LangInit("zh-cn") // en-us

	// 初始化 UserService
	global.UserService = service.NewUserService(global.UserRepo, global.ItemIconGroupRepo)

	// 初始化数据库
	err = DatabaseConnect()
	if err != nil {
		return err
	}

	// 初始化存储系统
	global.Storage, err = InitStorage(configPath)
	if err != nil {
		return fmt.Errorf("storage initialization error: %w", err)
	}

	// 初始化JWT
	if err := jwt.InitJWT(); err != nil {
		return fmt.Errorf("JWT initialization error: %w", err)
	}

	// 初始化路由
	httpPort := config.AppConfig.Base.HTTPPort
	if err := router.InitRouters(":" + httpPort); err != nil {
		panic(err)
	}

	return nil
}

func DatabaseConnect() error {
	var dbClientInfo database.DbClient
	databaseDrive := config.AppConfig.Base.DatabaseDrive
	switch databaseDrive {
	case database.MYSQL:
		dbClientInfo = &database.MySQLConfig{
			Username:    config.AppConfig.MySQL.Username,
			Password:    config.AppConfig.MySQL.Password,
			Host:        config.AppConfig.MySQL.Host,
			Port:        config.AppConfig.MySQL.Port,
			Database:    config.AppConfig.MySQL.DBName,
			WaitTimeout: config.AppConfig.MySQL.WaitTimeout,
		}
	case database.SQLITE:
		dbClientInfo = &database.SQLiteConfig{
			Filename: config.AppConfig.SQLite.FilePath,
		}
	default:
		return fmt.Errorf("unsupported database drive: %s", databaseDrive)
	}

	db, err := database.DbInit(dbClientInfo)
	if err != nil {
		return fmt.Errorf("database DbInit error, %w", err)
	}

	repository.Db = db
	err = database.CreateDefaultUser()
	if err != nil {
		return fmt.Errorf("database CreateDefaultUser error, %w", err)
	}

	return nil
}

// initializes the storage system based on configuration
func InitStorage(configPath string) (*storage.RcloneStorage, error) {
	// 使用带超时的上下文初始化存储
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bucket := config.AppConfig.Rclone.Bucket
	rcloneStorage, err := storage.NewRcloneStorage(ctx, configPath, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rclone storage: %w", err)
	}

	zaplog.Logger.Info("Storage system initialized successfully")
	return rcloneStorage, nil
}

func Logo() {
	fmt.Println("     ____            ___                __")
	fmt.Println("    / __/_ _____    / _ \\___ ____  ___ / /")
	fmt.Println("   _\\ \\/ // / _ \\  / ___/ _ `/ _ \\/ -_) / ")
	fmt.Println("  /___/\\_,_/_//_/ /_/   \\_,_/_//_/\\__/_/  ")
	fmt.Println("")

	fmt.Println("Version:", global.VERSION)
	fmt.Println("Welcome to the Moon-Box.")
	fmt.Println("Project address:", "https://github.com/hslr-s/sun-panel")
}
