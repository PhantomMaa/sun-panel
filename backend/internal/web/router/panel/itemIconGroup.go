package panel

import (
	"math"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/constant"
	"sun-panel/internal/global"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/param/commonApi"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
)

type ItemIconGroupRouter struct {
}

func NewItemIconGroupRouter() *ItemIconGroupRouter {
	return &ItemIconGroupRouter{}
}

func (a *ItemIconGroupRouter) InitRouter(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(interceptor.JWTAuth)
	{
		r.POST("/panel/itemIconGroup/edit", a.Edit)
		r.POST("/panel/itemIconGroup/deletes", a.Deletes)
		r.POST("/panel/itemIconGroup/saveSort", a.SaveSort)
		r.GET("/panel/itemIconGroup/getGroups", a.GetGroups)
	}

	// public visit 路由组
	publicR := router.Group(":code")
	publicR.Use(interceptor.PublicAccess)
	{
		publicR.GET("/panel/itemIconGroup/getGroups", a.GetGroups)
	}
}

func (a *ItemIconGroupRouter) Edit(c *gin.Context) {
	userInfo, exist := base.GetCurrentUserInfo(c)
	if !exist || userInfo.ID == 0 {
		response.ErrorByCode(c, constant.CodeNotLogin)
		return
	}

	itemIconGroup := &repository.ItemIconGroup{}

	if err := c.ShouldBindBodyWith(itemIconGroup, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	if itemIconGroup.ID != 0 && itemIconGroup.UserId != userInfo.ID {
		response.ErrorCode(c, 1203, "You do not have permission to edit this item", nil)
		return
	}

	itemIconGroup.UserId = userInfo.ID
	if err := global.ItemIconGroupRepo.Save(itemIconGroup); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.SuccessData(c, itemIconGroup)
}

func (a *ItemIconGroupRouter) GetGroups(c *gin.Context) {
	userInfo, exist := base.GetCurrentUserInfo(c)
	if !exist || userInfo.ID == 0 {
		response.ErrorByCode(c, constant.CodeNotLogin)
		return
	}

	groups, err := global.ItemIconGroupRepo.GetList(userInfo.ID)
	if err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.SuccessListData(c, groups, 0)
}

func (a *ItemIconGroupRouter) Deletes(c *gin.Context) {
	req := commonApi.RequestDeleteIds[uint]{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, exist := base.GetCurrentUserInfo(c)
	if !exist || userInfo.ID == 0 {
		response.ErrorByCode(c, constant.CodeNotLogin)
		return
	}

	if count, err := global.ItemIconGroupRepo.Count(userInfo.ID); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	} else {
		if math.Abs(float64(len(req.Ids))-float64(count)) < 1 {
			response.ErrorCode(c, 1201, "At least one must be retained", nil)
			return
		}
	}

	if err := global.ItemIconGroupRepo.Deletes(userInfo.ID, req.Ids); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.Success(c)
}

func (a *ItemIconGroupRouter) SaveSort(c *gin.Context) {
	req := commonApi.SortRequest{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, exist := base.GetCurrentUserInfo(c)
	if !exist || userInfo.ID == 0 {
		response.ErrorByCode(c, constant.CodeNotLogin)
		return
	}

	err := global.ItemIconGroupRepo.BatchSaveSort(userInfo.ID, req.SortItems)
	if err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.Success(c)
}
