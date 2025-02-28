package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sun-panel/internal/cache"
	"sun-panel/internal/iniConfig"
)

var (
	RUNCODE = "debug" // 运行模式：debug | release
	VERSION = "v1.0.0"
)

var (
	Logger        *zap.SugaredLogger
	LoggerLevel   = zap.NewAtomicLevel() // 支持通过http以及配置文件动态修改日志级别
	Config        *iniConfig.IniConfig
	Db            *gorm.DB
	SystemSetting *cache.SystemSetting
	Monitor       *cache.Monitor
)
