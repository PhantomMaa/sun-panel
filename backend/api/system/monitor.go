package system

import (
	"sun-panel/api/common/apiData/systemApiStructs"
	"sun-panel/api/common/apiReturn"
	"sun-panel/internal/global"
	"sun-panel/internal/monitor"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type MonitorApi struct{}

const cacheSecond = 3

func (a *MonitorApi) GetAll(c *gin.Context) {
	if value, ok := global.SystemMonitor.Get("value"); ok {
		apiReturn.SuccessData(c, value)
		return
	}
	apiReturn.Error(c, "failed")
}

func (a *MonitorApi) GetCpuState(c *gin.Context) {
	if v, ok := global.SystemMonitor.Get(monitor.SystemmonitorCpuInfo); ok {
		global.Logger.Debugln("读取缓存的的CPU信息")
		apiReturn.SuccessData(c, v)
		return
	}
	cpuInfo, err := monitor.GetCPUInfo()

	if err != nil {
		apiReturn.Error(c, "failed")
		return
	}
	// 缓存
	global.SystemMonitor.Set(monitor.SystemmonitorCpuInfo, cpuInfo, cacheSecond*time.Second)
	apiReturn.SuccessData(c, cpuInfo)
}

func (a *MonitorApi) GetMemonyState(c *gin.Context) {
	if v, ok := global.SystemMonitor.Get(monitor.SystemmonitorMemoryInfo); ok {
		global.Logger.Debugln("读取缓存的的RAM信息")
		apiReturn.SuccessData(c, v)
		return
	}
	memoryInfo, err := monitor.GetMemoryInfo()

	if err != nil {
		apiReturn.Error(c, "failed")
		return
	}

	// 缓存
	global.SystemMonitor.Set(monitor.SystemmonitorMemoryInfo, memoryInfo, cacheSecond*time.Second)
	apiReturn.SuccessData(c, memoryInfo)
}

func (a *MonitorApi) GetDiskStateByPath(c *gin.Context) {

	req := systemApiStructs.MonitorGetDiskStateByPathReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	cacheDiskName := monitor.SystemmonitorDiskInfo + req.Path

	if v, ok := global.SystemMonitor.Get(cacheDiskName); ok {
		global.Logger.Debugln("读取缓存的的DISK信息")
		apiReturn.SuccessData(c, v)
		return
	}

	diskState, err := monitor.GetDiskInfoByPath(req.Path)
	if err != nil {
		apiReturn.Error(c, "failed")
		return
	}

	// 缓存
	global.SystemMonitor.Set(cacheDiskName, diskState, cacheSecond*time.Second)
	apiReturn.SuccessData(c, diskState)
}

func (a *MonitorApi) GetDiskMountpoints(c *gin.Context) {
	if list, err := monitor.GetDiskMountpoints(); err != nil {
		apiReturn.Error(c, err.Error())
		return
	} else {
		apiReturn.SuccessData(c, list)
	}
}
