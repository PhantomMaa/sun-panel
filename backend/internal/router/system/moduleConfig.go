package system

import (
	"sun-panel/api"
	"sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitModuleConfigRouter(router *gin.RouterGroup) {
	moduleConfigApi := api.ApiGroupApp.ApiSystem.ModuleConfigApi
	r := router.Group("")
	r.Use(middleware.JWTAuth())
	{
		r.GET("/system/moduleConfig/getByName", moduleConfigApi.GetByName)
		r.POST("/system/moduleConfig/save", moduleConfigApi.Save)
	}
}
