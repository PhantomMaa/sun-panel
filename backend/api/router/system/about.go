package system

import (
	"github.com/gin-gonic/gin"
	"sun-panel/api/common/apiReturn"
	"sun-panel/internal/global"
)

type About struct {
}

func NewAboutRouter() *About {
	return &About{}
}

func (a *About) Get(c *gin.Context) {
	apiReturn.SuccessData(c, gin.H{
		"versionName": global.VERSION,
	})
}

func (a *About) InitRouter(router *gin.RouterGroup) {
	{
		router.GET("about", a.Get)
	}
}
