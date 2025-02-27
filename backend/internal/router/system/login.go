package system

import (
	"github.com/gin-gonic/gin"
	"sun-panel/api/middleware"
	api_v1 "sun-panel/api/system"
)

// InitLogin 初始化登录相关路由
func InitLogin(router *gin.RouterGroup) {
	loginApi := api_v1.LoginApi{}

	// 公开接口
	router.POST("/login", loginApi.Login)

	// 需要认证的接口
	authGroup := router.Group("")
	authGroup.Use(middleware.JWTAuth())
	{
		authGroup.POST("/logout", loginApi.Logout)
	}
}
