package panel

import (
	"sun-panel/api"
	middleware2 "sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitUsersRouter(router *gin.RouterGroup) {
	userApi := api.ApiGroupApp.ApiPanel.UsersApi

	rAdmin := router.Group("")
	rAdmin.Use(middleware2.JWTAuth(), middleware2.AdminInterceptor)
	{
		rAdmin.POST("panel/users/create", userApi.Create)
		rAdmin.POST("panel/users/update", userApi.Update)
		rAdmin.POST("panel/users/getList", userApi.GetList)
		rAdmin.POST("panel/users/deletes", userApi.Deletes)
		rAdmin.POST("panel/users/getPublicVisitUser", userApi.GetPublicVisitUser)
		rAdmin.POST("panel/users/setPublicVisitUser", userApi.SetPublicVisitUser)
	}
}
