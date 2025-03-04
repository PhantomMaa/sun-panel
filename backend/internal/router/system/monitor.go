package system

import (
	"sun-panel/api"
	"sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitMonitorRouter(router *gin.RouterGroup) {
	monitorApi := api.ApiGroupApp.ApiSystem.MonitorApi
	r := router.Group("")
	r.Use(middleware.JWTAuth())

	{
		r.POST("/system/monitor/getDiskMountpoints", monitorApi.GetDiskMountpoints)
		r.POST("/system/monitor/getCpuState", monitorApi.GetCpuState)
		r.POST("/system/monitor/getDiskStateByPath", monitorApi.GetDiskStateByPath)
		r.POST("/system/monitor/getMemonyState", monitorApi.GetMemonyState)
	}
}
