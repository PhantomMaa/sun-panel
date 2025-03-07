package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sun-panel/internal/biz/cache"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/biz/service"
	"sun-panel/internal/infra/config"
	"sun-panel/internal/infra/storage"
)

// 构建时，通过 --ldflags 注入
var (
	RUNCODE = "debug" // 运行模式：debug | release
	VERSION = "v1.0.0"
)

var (
	Logger             *zap.SugaredLogger
	Config             *config.IniConfig
	Db                 *gorm.DB
	Storage            *storage.RcloneStorage
	CacheSystemSetting *cache.SystemSetting
	CacheMonitor       *cache.Monitor
)

var (
	ItemIconRepo      = repository.NewItemIconRepo()
	ItemIconGroupRepo = repository.NewItemIconGroupRepo()
	UserRepo          = repository.NewUserRepo()
	FileRepo          = repository.NewFileRepo()
	ModuleConfigRepo  = repository.NewModuleConfigRepo()
)

var (
	UserService = service.NewUserService(UserRepo, ItemIconGroupRepo)
)
