package system

import (
	"sun-panel/api"
	middleware2 "sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitMonitorRouter(router *gin.RouterGroup) {
	api := api.ApiGroupApp.ApiSystem.MonitorApi
	r := router.Group("")
	r.Use(middleware2.JWTAuth())
	r.POST("/system/monitor/getDiskMountpoints", api.GetDiskMountpoints)

	// 公开模式
	rPublic := router.Group("", middleware2.PublicModeInterceptor)
	{
		rPublic.POST("/system/monitor/getAll", api.GetAll)
		rPublic.POST("/system/monitor/getCpuState", api.GetCpuState)
		rPublic.POST("/system/monitor/getDiskStateByPath", api.GetDiskStateByPath)
		rPublic.POST("/system/monitor/getMemonyState", api.GetMemonyState)
	}
}
