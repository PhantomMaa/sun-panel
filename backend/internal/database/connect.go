package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sun-panel/internal/common"
	repository2 "sun-panel/internal/repository"
	"time"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	MYSQL  = "mysql"
	SQLITE = "sqlite"
)

type DbClient interface {
	Connect() (db *gorm.DB, err error)
	InitDatabase(db *gorm.DB) (err error)
}

type MySQLConfig struct {
	Username    string
	Password    string
	Host        string
	Port        string
	Database    string
	WaitTimeout int
}

type SQLiteConfig struct {
	Filename string
}

func DbInit(dbClient DbClient) (db *gorm.DB, dbErr error) {
	db, dbErr = dbClient.Connect()
	if dbErr != nil {
		return
	}

	dbErr = dbClient.InitDatabase(db)
	if dbErr != nil {
		return nil, fmt.Errorf("database CreateDatabase error, %w", dbErr)
	}

	return
}

func (d *MySQLConfig) Connect() (db *gorm.DB, err error) {
	dsn := d.Username + ":" + d.Password + "@tcp(" + d.Host + ":" + d.Port + ")/" + d.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: GetLogger(),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(10)                                                 // SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxOpenConns(100)                                                // SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDb.SetConnMaxLifetime(time.Duration(d.WaitTimeout * int(time.Second))) // SetConnMaxLifetime 设置了连接可复用的最大时间。
	return
}

func (d *SQLiteConfig) Connect() (db *gorm.DB, err error) {
	filePath := d.Filename
	exists := false
	if exists, err = common.PathExists(path.Dir(filePath)); err != nil {
		return
	} else {
		// 创建文件夹
		if !exists {
			if err = os.MkdirAll(path.Dir(filePath), 0700); err != nil {
				return
			}
		}

		db, err = gorm.Open(sqlite.Open(filePath), &gorm.Config{
			Logger: GetLogger(),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	}

	return
}

func (d *MySQLConfig) InitDatabase(db *gorm.DB) (err error) {
	db = db.Set("gorm:table_options", "ENGINE=InnoDB")

	// 创建数据表
	err = db.AutoMigrate(
		&repository2.User{},
		&repository2.SystemSetting{},
		&repository2.ItemIcon{},
		&repository2.UserConfig{},
		&repository2.File{},
		&repository2.ItemIconGroup{},
		&repository2.ModuleConfig{},
	)

	return err
}

func (d *SQLiteConfig) InitDatabase(db *gorm.DB) (err error) {
	return
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

func NotFoundAndCreateUser(db *gorm.DB) error {
	fUser := repository2.User{}
	if err := db.First(&fUser).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		username := "admin@sun.cc"
		fUser.Mail = username
		fUser.Username = username
		fUser.Name = username
		fUser.Status = 1
		fUser.Role = 1
		fUser.Password = common.PasswordEncryption("12345678")

		if errCreate := db.Create(&fUser).Error; errCreate != nil {
			return errCreate
		}
	}

	return nil
}
