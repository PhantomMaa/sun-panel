package panel

import (
	"sun-panel/api"
	middleware2 "sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitItemIconGroup(router *gin.RouterGroup) {
	itemIconGroup := api.ApiGroupApp.ApiPanel.ItemIconGroup
	r := router.Group("")
	r.Use(middleware2.JWTAuth())
	{
		r.POST("/panel/itemIconGroup/edit", itemIconGroup.Edit)
		r.POST("/panel/itemIconGroup/deletes", itemIconGroup.Deletes)
		r.POST("/panel/itemIconGroup/saveSort", itemIconGroup.SaveSort)
	}

	// 公开模式
	rPublic := router.Group("", middleware2.PublicModeInterceptor)
	{
		rPublic.GET("/panel/itemIconGroup/getList", itemIconGroup.GetList)
	}
}
