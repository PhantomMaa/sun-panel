package system

import (
	"github.com/gin-gonic/gin"
	"sun-panel/internal/global"
	"sun-panel/internal/web/model/response"
)

type About struct {
}

func NewAboutRouter() *About {
	return &About{}
}

func (a *About) Get(c *gin.Context) {
	response.SuccessData(c, gin.H{
		"versionName": global.VERSION,
	})
}

func (a *About) InitRouter(router *gin.RouterGroup) {
	router.GET("about", a.Get)
}
