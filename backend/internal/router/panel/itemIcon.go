package panel

import (
	"sun-panel/api"
	middleware2 "sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitItemIcon(router *gin.RouterGroup) {
	itemIcon := api.ApiGroupApp.ApiPanel.ItemIcon
	r := router.Group("")
	r.Use(middleware2.JWTAuth())
	{
		r.POST("/panel/itemIcon/edit", itemIcon.Edit)
		r.POST("/panel/itemIcon/deletes", itemIcon.Deletes)
		r.POST("/panel/itemIcon/saveSort", itemIcon.SaveSort)
		r.POST("/panel/itemIcon/addMultiple", itemIcon.AddMultiple)
		r.POST("/panel/itemIcon/getSiteFavicon", itemIcon.GetSiteFavicon)
	}

	// 公开模式
	rPublic := router.Group("", middleware2.PublicModeInterceptor)
	{
		rPublic.POST("/panel/itemIcon/getListByGroupId", itemIcon.GetListByGroupId)
	}
}
