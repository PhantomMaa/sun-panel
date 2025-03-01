package api

import (
	"sun-panel/api/panel"
	"sun-panel/api/system"
	"sun-panel/internal/storage"
)

type ApiGroup struct {
	ApiSystem *system.ApiSystem // 系统功能api
	ApiPanel  *panel.ApiPanel
}

var ApiGroupApp = &ApiGroup{}

// InitApiGroup initializes the API group with required dependencies
func InitApiGroup(storageInstance *storage.RcloneStorage) {
	ApiGroupApp.ApiSystem = system.InitApiSystem(storageInstance)
	ApiGroupApp.ApiPanel = panel.InitApiPanel(storageInstance)
}
