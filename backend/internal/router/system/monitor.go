package system

import (
	"sun-panel/api"
	middleware2 "sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitMonitorRouter(router *gin.RouterGroup) {
	monitorApi := api.ApiGroupApp.ApiSystem.MonitorApi
	r := router.Group("")
	r.Use(middleware2.JWTAuth())
	r.POST("/system/monitor/getDiskMountpoints", monitorApi.GetDiskMountpoints)

	// 公开模式
	rPublic := router.Group("", middleware2.PublicModeInterceptor)
	{
		rPublic.POST("/system/monitor/getCpuState", monitorApi.GetCpuState)
		rPublic.POST("/system/monitor/getDiskStateByPath", monitorApi.GetDiskStateByPath)
		rPublic.POST("/system/monitor/getMemonyState", monitorApi.GetMemonyState)
	}
}
