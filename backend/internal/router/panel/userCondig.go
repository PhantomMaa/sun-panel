package panel

import (
	"sun-panel/api"
	"sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserConfig(router *gin.RouterGroup) {
	userConfig := api.ApiGroupApp.ApiPanel.UserConfig
	r := router.Group("")
	r.Use(middleware.JWTAuth())
	{
		r.POST("/panel/userConfig/set", userConfig.Set)
		r.GET("/panel/userConfig/get", userConfig.Get)
	}
}
