package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sun-panel/internal/cache"
	"sun-panel/internal/iniConfig"
	"sun-panel/internal/storage"
)

// 构建时，通过 --ldflags 注入
var (
	RUNCODE = "debug" // 运行模式：debug | release
	VERSION = "v1.0.0"
)

var (
	Logger        *zap.SugaredLogger
	Config        *iniConfig.IniConfig
	Db            *gorm.DB
	SystemSetting *cache.SystemSetting
	Monitor       *cache.Monitor
	Storage       *storage.RcloneStorage
)
