package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.UserApi
	r := router.Group("")
	r.Use(middleware.JWTAuth())
	r.POST("/user/getInfo", api.GetInfo)
	r.POST("/user/updatePassword", api.UpdatePasssword)
	r.POST("/user/updateInfo", api.UpdateInfo)

	// 公开模式
	rPublic := router.Group("", middleware.PublicModeInterceptor)
	{
		rPublic.POST("/user/getAuthInfo", api.GetAuthInfo)
	}
}
