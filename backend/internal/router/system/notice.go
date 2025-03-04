package system

import (
	"github.com/gin-gonic/gin"
	"sun-panel/api"
)

func InitNoticeRouter(router *gin.RouterGroup) {
	noticeApi := api.ApiGroupApp.ApiSystem.NoticeApi
	router.POST("/notice/getListByDisplayType", noticeApi.GetListByDisplayType)
}
