package system

import (
	"github.com/gin-gonic/gin"
	"sun-panel/api"
)

func InitNoticeRouter(router *gin.RouterGroup) {
	api := api.ApiGroupApp.ApiSystem.NoticeApi

	router.POST("/notice/getListByDisplayType", api.GetListByDisplayType)
}
