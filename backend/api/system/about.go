package system

import (
	"sun-panel/api/common/apiReturn"
	"sun-panel/internal/global"

	"github.com/gin-gonic/gin"
)

type About struct {
}

func (a *About) Get(c *gin.Context) {
	apiReturn.SuccessData(c, gin.H{
		"versionName": global.VERSION,
	})
}
