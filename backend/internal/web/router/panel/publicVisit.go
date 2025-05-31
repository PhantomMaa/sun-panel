package panel

import (
	"sun-panel/internal/constant"
	"sun-panel/internal/global"
	"sun-panel/internal/infra/zaplog"
	"sun-panel/internal/util/publiccode"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

type PublicVisitRouter struct {
}

func NewPublicVisitRouter() *PublicVisitRouter {
	return &PublicVisitRouter{}
}

func (a *PublicVisitRouter) InitRouter(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(interceptor.Auth)
	{
		r.POST("/panel/publicVisit/enable", a.Enable)
		r.POST("/panel/publicVisit/disable", a.Disable)
	}
}

// Enable 启用公开访问代码
func (a *PublicVisitRouter) Enable(c *gin.Context) {
	// 获取当前用户ID
	userInfo, exist := base.GetCurrentUserInfo(c)
	if !exist || userInfo.ID == 0 {
		response.ErrorByCode(c, constant.CodeNotLogin)
		return
	}

	// 生成公开访问代码
	code, err := publiccode.GenerateAndSave(userInfo.ID)
	if err != nil {
		zaplog.Logger.Errorf("生成公开访问代码失败: %v", err)
		response.ErrorByCode(c, constant.CodeStatusError)
		return
	}

	// 返回生成的代码
	response.SuccessData(c, gin.H{
		"code": code,
	})
}

// Disable 禁用公开访问代码
func (a *PublicVisitRouter) Disable(c *gin.Context) {
	// 获取当前用户ID
	userInfo, exist := base.GetCurrentUserInfo(c)
	if !exist || userInfo.ID == 0 {
		response.ErrorByCode(c, constant.CodeNotLogin)
		return
	}

	// 清空用户的 publiccode 字段
	updateInfo := map[string]any{
		"publiccode": "",
	}

	if err := global.UserRepo.UpdateUserInfo(userInfo.ID, updateInfo); err != nil {
		zaplog.Logger.Errorf("移除公开访问代码失败: %v", err)
		response.ErrorByCode(c, constant.CodeStatusError)
		return
	}

	response.Success(c)
}
