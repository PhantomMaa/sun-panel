package system

import (
	"sun-panel/api"
	"sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	userApi := api.ApiGroupApp.ApiSystem.UserApi
	r := router.Group("")
	r.Use(middleware.JWTAuth())
	{
		r.GET("/user/getInfo", userApi.GetInfo)
		r.POST("/user/updatePassword", userApi.UpdatePasssword)
		r.POST("/user/updateInfo", userApi.UpdateInfo)
		r.GET("/user/getAuthInfo", userApi.GetAuthInfo)
	}
}
