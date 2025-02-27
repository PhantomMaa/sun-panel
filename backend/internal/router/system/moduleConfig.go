package system

import (
	"sun-panel/api"
	middleware2 "sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitModuleConfigRouter(router *gin.RouterGroup) {
	api := api.ApiGroupApp.ApiSystem.ModuleConfigApi
	r := router.Group("")
	r.Use(middleware2.JWTAuth())
	r.POST("/system/moduleConfig/save", api.Save)

	// 公开模式
	rPublic := router.Group("", middleware2.PublicModeInterceptor)
	{
		rPublic.POST("/system/moduleConfig/getByName", api.GetByName)
	}
}
