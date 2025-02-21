package api_v1

import (
	"sun-panel/api/api_v1/panel"
	"sun-panel/api/api_v1/system"
)

type ApiGroup struct {
	ApiSystem system.ApiSystem // 系统功能api
	ApiPanel  panel.ApiPanel
}

var ApiGroupApp = new(ApiGroup)
