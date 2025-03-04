package panel

import (
	"sun-panel/api"
	"sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitItemIcon(router *gin.RouterGroup) {
	itemIcon := api.ApiGroupApp.ApiPanel.ItemIcon
	r := router.Group("")
	r.Use(middleware.JWTAuth())
	{
		r.POST("/panel/itemIcon/edit", itemIcon.Edit)
		r.POST("/panel/itemIcon/deletes", itemIcon.Deletes)
		r.POST("/panel/itemIcon/saveSort", itemIcon.SaveSort)
		r.POST("/panel/itemIcon/addMultiple", itemIcon.AddMultiple)
		r.POST("/panel/itemIcon/getSiteFavicon", itemIcon.GetSiteFavicon)
		r.GET("/panel/itemIcon/getListByGroupId", itemIcon.GetListByGroupId)
	}
}
