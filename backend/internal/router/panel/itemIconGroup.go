package panel

import (
	"sun-panel/api"
	"sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitItemIconGroup(router *gin.RouterGroup) {
	itemIconGroup := api.ApiGroupApp.ApiPanel.ItemIconGroup
	r := router.Group("")
	r.Use(middleware.JWTAuth())
	{
		r.POST("/panel/itemIconGroup/edit", itemIconGroup.Edit)
		r.POST("/panel/itemIconGroup/deletes", itemIconGroup.Deletes)
		r.POST("/panel/itemIconGroup/saveSort", itemIconGroup.SaveSort)
		r.GET("/panel/itemIconGroup/getList", itemIconGroup.GetList)
	}
}
