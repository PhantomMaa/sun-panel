package system

import (
	"sun-panel/internal/global"
	"sun-panel/internal/util/monitor"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/param/systemApi"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

type MonitorRouter struct {
}

func NewMonitorRouter() *MonitorRouter {
	return &MonitorRouter{}
}

func (a *MonitorRouter) InitRouter(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(interceptor.JWTAuth)

	{
		r.POST("/system/monitor/getDiskMountpoints", a.GetDiskMountpoints)
		r.POST("/system/monitor/getCpuState", a.GetCpuState)
		r.POST("/system/monitor/getDiskStateByPath", a.GetDiskStateByPath)
		r.POST("/system/monitor/getMemonyState", a.GetMemonyState)
		r.POST("/system/monitor/getEnableStatus", a.GetEnableStatus)
	}
}

func (a *MonitorRouter) GetCpuState(c *gin.Context) {
	cpuInfo, err := global.CacheMonitor.GetCpuState()
	if err != nil {
		response.Error(c, "failed")
		return
	}

	response.SuccessData(c, cpuInfo)
}

func (a *MonitorRouter) GetMemonyState(c *gin.Context) {
	memoryInfo, err := global.CacheMonitor.GetMemonyState()
	if err != nil {
		response.Error(c, "failed")
		return
	}

	response.SuccessData(c, memoryInfo)
}

func (a *MonitorRouter) GetDiskStateByPath(c *gin.Context) {
	req := systemApi.MonitorGetDiskStateByPathReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	diskState, err := global.CacheMonitor.GetDiskStateByPath(req.Path)
	if err != nil {
		response.Error(c, "failed")
		return
	}

	response.SuccessData(c, diskState)
}

func (a *MonitorRouter) GetDiskMountpoints(c *gin.Context) {
	if list, err := monitor.GetDiskMountpoints(); err != nil {
		response.Error(c, err.Error())
	} else {
		response.SuccessData(c, list)
	}
}

// GetEnableStatus returns the enableMonitor configuration from conf.ini
func (a *MonitorRouter) GetEnableStatus(c *gin.Context) {
	response.SuccessData(c, gin.H{
		"enabled": global.Config.Base.EnableMonitor,
	})
}
