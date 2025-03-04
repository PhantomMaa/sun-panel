package system

import (
	"github.com/gin-gonic/gin"
	"sun-panel/api/common/apiData/systemApiStructs"
	"sun-panel/api/common/apiReturn"
	"sun-panel/internal/global"
	"sun-panel/internal/repository"
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
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	var noticeList []repository.Notice
	if err := global.Db.Find(&noticeList, "display_type in ?", req.DisplayType).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, noticeList, 0)
}
