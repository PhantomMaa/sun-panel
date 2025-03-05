package system

import (
	"github.com/gin-gonic/gin"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/global"
	"sun-panel/internal/web/model/param/systemApiStructs"
	"sun-panel/internal/web/model/response"
)

type NoticeApi struct {
}

func NewNoticeRouter() *NoticeApi {
	return &NoticeApi{}
}

func (a *NoticeApi) InitRouter(router *gin.RouterGroup) {
	router.POST("/notice/getListByDisplayType", a.GetListByDisplayType)
}

func (a *NoticeApi) GetListByDisplayType(c *gin.Context) {
	req := systemApiStructs.NoticeGetListByDisplayTypeReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	var noticeList []repository.Notice
	if err := global.Db.Find(&noticeList, "display_type in ?", req.DisplayType).Error; err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.SuccessListData(c, noticeList, 0)
}
