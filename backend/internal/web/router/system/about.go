package system

import (
	"sun-panel/internal/global"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

type AboutRouter struct {
}

func NewAboutRouter() *AboutRouter {
	return &AboutRouter{}
}

func (a *AboutRouter) Get(c *gin.Context) {
	response.SuccessData(c, gin.H{
		"versionName": global.VERSION,
	})
}

func (a *AboutRouter) InitRouter(router *gin.RouterGroup) {
	// 公开接口
	router.GET("about", a.Get)
}
