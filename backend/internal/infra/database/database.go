package database

import (
	"fmt"
	"log"
	"os"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/constant"
	"sun-panel/internal/global"
	"sun-panel/internal/util"
	"time"

	_ "gorm.io/driver/mysql"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	MYSQL  = "mysql"
	SQLITE = "sqlite"
)

type DbClient interface {
	Connect() (db *gorm.DB, err error)
}

func DbInit(dbClient DbClient) (db *gorm.DB, dbErr error) {
	db, dbErr = dbClient.Connect()
	if dbErr != nil {
		return
	}

	dbErr = initDatabase(db)
	if dbErr != nil {
		return nil, fmt.Errorf("database CreateDatabase error, %w", dbErr)
	}

	return
}

func initDatabase(db *gorm.DB) (err error) {
	// 创建数据表
	err = db.AutoMigrate(
		&repository.User{},
		&repository.SystemSetting{},
		&repository.ItemIcon{},
		&repository.UserConfig{},
		&repository.File{},
		&repository.ItemIconGroup{},
		&repository.ModuleConfig{},
	)

	return err
}

func CreateDefaultUser() error {
	count, err := global.UserRepo.Count()
	if err != nil {
		return err
	}

	if count == 0 {
		mUser := repository.User{}
		mail := "admin@sun.cc"
		mUser.Mail = mail
		mUser.Username = mail
		mUser.Name = mail
		mUser.Status = 1
		mUser.Role = 1
		mUser.Password = util.PasswordEncryption("12345678")
		mUser.OauthProvider = constant.OAuthProviderBuildin
		if errCreate := global.UserService.CreateUser(&mUser); errCreate != nil {
			return errCreate
		}
	}

	return nil
}

func GetLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Warn, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 彩色打印
		},
	)

}
