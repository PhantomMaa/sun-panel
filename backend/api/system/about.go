package system

import (
	"sun-panel/api/common/apiReturn"
	"sun-panel/internal/common"

	"github.com/gin-gonic/gin"
)

type About struct {
}

func (a *About) Get(c *gin.Context) {
	version := common.GetVersion()
	apiReturn.SuccessData(c, gin.H{
		"versionName": version,
	})
}
