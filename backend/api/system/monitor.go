package system

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"sun-panel/api/common/apiData/systemApiStructs"
	"sun-panel/api/common/apiReturn"
	"sun-panel/internal/global"
	"sun-panel/internal/monitor"
)

type MonitorApi struct {
}

func (a *MonitorApi) GetCpuState(c *gin.Context) {
	cpuInfo, err := global.Monitor.GetCpuState()
	if err != nil {
		apiReturn.Error(c, "failed")
		return
	}

	apiReturn.SuccessData(c, cpuInfo)
}

func (a *MonitorApi) GetMemonyState(c *gin.Context) {
	memoryInfo, err := global.Monitor.GetMemonyState()
	if err != nil {
		apiReturn.Error(c, "failed")
		return
	}

	apiReturn.SuccessData(c, memoryInfo)
}

func (a *MonitorApi) GetDiskStateByPath(c *gin.Context) {
	req := systemApiStructs.MonitorGetDiskStateByPathReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	diskState, err := global.Monitor.GetDiskStateByPath(req.Path)
	if err != nil {
		apiReturn.Error(c, "failed")
		return
	}

	apiReturn.SuccessData(c, diskState)
}

func (a *MonitorApi) GetDiskMountpoints(c *gin.Context) {
	if list, err := monitor.GetDiskMountpoints(); err != nil {
		apiReturn.Error(c, err.Error())
	} else {
		apiReturn.SuccessData(c, list)
	}
}
