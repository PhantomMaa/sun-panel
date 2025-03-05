package database

import (
	"time"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type MySQLConfig struct {
	Username    string
	Password    string
	Host        string
	Port        string
	Database    string
	WaitTimeout int
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
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	sqlDb.SetMaxIdleConns(10)                                                 // SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxOpenConns(100)                                                // SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDb.SetConnMaxLifetime(time.Duration(d.WaitTimeout * int(time.Second))) // SetConnMaxLifetime 设置了连接可复用的最大时间。
	return
}
