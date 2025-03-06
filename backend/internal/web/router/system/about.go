package system

import (
	"github.com/gin-gonic/gin"
	"sun-panel/internal/global"
	"sun-panel/internal/web/model/response"
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
	router.GET("about", a.Get)
}
