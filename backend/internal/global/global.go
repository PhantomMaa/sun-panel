package global

import (
	"sun-panel/internal/biz/cache"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/biz/service"
	"sun-panel/internal/infra/config"
	"sun-panel/internal/infra/kvcache"
	"sun-panel/internal/infra/storage"
	"time"
)

// 构建时，通过 --ldflags 注入
var (
	RUNCODE = "debug" // 运行模式：debug | release
	VERSION = "v1.0.0"
)

var (
	Config  *config.Config
	Storage *storage.RcloneStorage
)

// repositories
var (
	ItemIconRepo      = repository.NewItemIconRepo()
	ItemIconGroupRepo = repository.NewItemIconGroupRepo()
	UserRepo          = repository.NewUserRepo()
	FileRepo          = repository.NewFileRepo()
	ModuleConfigRepo  = repository.NewModuleConfigRepo()
	UserConfigRepo    = repository.NewUserConfigRepo()
	SystemSettingRepo = repository.NewSystemSettingRepo()
)

// services
var (
	UserService *service.UserService
)

// caches
var (
	CacheSystemSetting = &cache.SystemSetting{
		Cache:             kvcache.NewLocalCache[any](5*time.Hour, -1),
		SystemSettingRepo: SystemSettingRepo,
	}
	CacheMonitor = &cache.Monitor{
		Cache: kvcache.NewLocalCache[any](5*time.Hour, -1),
	}
)
