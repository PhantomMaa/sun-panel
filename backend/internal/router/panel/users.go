package panel

import (
	"sun-panel/api"
	"sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitUsersRouter(router *gin.RouterGroup) {
	userApi := api.ApiGroupApp.ApiPanel.UsersApi

	rAdmin := router.Group("")
	rAdmin.Use(middleware.JWTAuth(), middleware.AdminInterceptor)
	{
		rAdmin.POST("panel/users/create", userApi.Create)
		rAdmin.GET("panel/users/getList", userApi.GetList)
		rAdmin.POST("panel/users/update", userApi.Update)
		rAdmin.POST("panel/users/deletes", userApi.Deletes)
	}
}
