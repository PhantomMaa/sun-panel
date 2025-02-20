package global

import (
	"sun-panel/initialize/database"
	"sun-panel/lib/cache"
	"sun-panel/lib/cmn/systemSetting"
	"sun-panel/lib/iniConfig"
	"sun-panel/lib/language"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	RUNCODE = "debug" // 运行模式：debug | release
	// DB_MYSQL  = "mysql"
	// DB_SQLITE = "sqlite"
	DB_DRIVER = database.SQLITE

	// JWT配置
	JWTConfig struct {
		Secret string
		Expire int // 过期时间（小时）
	}
)

// var Log *cmn.LogStruct

var (
	Lang *language.LangStructObj

	Logger        *zap.SugaredLogger
	LoggerLevel   = zap.NewAtomicLevel() // 支持通过http以及配置文件动态修改日志级别
	Config        *iniConfig.IniConfig
	Db            *gorm.DB
	SystemSetting *systemSetting.SystemSettingCache
	SystemMonitor cache.Cacher[interface{}]
	RateLimit     *RateLimiter
)
