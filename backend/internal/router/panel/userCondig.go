package panel

import (
	"sun-panel/api"
	middleware2 "sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserConfig(router *gin.RouterGroup) {
	api := api.ApiGroupApp.ApiPanel.UserConfig
	r := router.Group("")
	r.Use(middleware2.JWTAuth())
	{
		r.POST("/panel/userConfig/set", api.Set)
	}

	// 公开模式
	rPublic := router.Group("", middleware2.PublicModeInterceptor)
	{
		rPublic.POST("/panel/userConfig/get", api.Get)
	}
}
