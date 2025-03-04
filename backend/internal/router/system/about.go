package system

import (
	"github.com/gin-gonic/gin"
	"sun-panel/api"
)

func InitAbout(router *gin.RouterGroup) {
	about := api.ApiGroupApp.ApiSystem.About
	{
		router.GET("about", about.Get)
	}
}
