package system

import (
	"sun-panel/api"
	middleware2 "sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	api := api.ApiGroupApp.ApiSystem.UserApi
	r := router.Group("")
	r.Use(middleware2.JWTAuth())
	r.POST("/user/getInfo", api.GetInfo)
	r.POST("/user/updatePassword", api.UpdatePasssword)
	r.POST("/user/updateInfo", api.UpdateInfo)

	// 公开模式
	rPublic := router.Group("", middleware2.PublicModeInterceptor)
	{
		rPublic.POST("/user/getAuthInfo", api.GetAuthInfo)
	}
}
