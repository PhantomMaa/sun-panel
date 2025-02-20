package global

import (
	"sun-panel/initialize/database"
	"sun-panel/lib/cache"
	"sun-panel/lib/iniConfig"
	"sun-panel/lib/language"
	"sun-panel/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SystemSettingCache represents the system settings cache interface
type SystemSettingCache interface {
	GetValue(key string) interface{}
	GetValueString(key string) (string, error)
	GetValueByInterface(key string, value interface{}) error
	Set(key string, value interface{}) error
	SetValue(key string, value interface{}) error
}

var (
	RUNCODE = "debug" // 运行模式：debug | release
	// DB_MYSQL  = "mysql"
	// DB_SQLITE = "sqlite"
	DB_DRIVER = database.SQLITE
)

// var Log *cmn.LogStruct

var (
	Lang *language.LangStructObj

	UserToken           cache.Cacher[models.User]
	CUserToken          cache.Cacher[string] // 用户token
	Logger              *zap.SugaredLogger
	LoggerLevel         = zap.NewAtomicLevel() // 支持通过http以及配置文件动态修改日志级别
	VerifyCodeCachePool cache.Cacher[string]
	Config              *iniConfig.IniConfig
	Db                  *gorm.DB
	SystemSetting       SystemSettingCache
	SystemMonitor       cache.Cacher[interface{}]
	RateLimit           *RateLimiter
)
