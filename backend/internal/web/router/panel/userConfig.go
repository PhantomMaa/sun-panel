package panel

import (
	"errors"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/constant"
	"sun-panel/internal/global"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type UserConfigRouter struct {
}

func NewUserConfigRouter() *UserConfigRouter {
	return &UserConfigRouter{}
}

func (a *UserConfigRouter) InitRouter(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(interceptor.Auth)
	{
		r.POST("/panel/userConfig/setConfig", a.SetConfig)
		r.GET("/panel/userConfig/getConfig", a.GetConfig)
	}
}

func (a *UserConfigRouter) GetConfig(c *gin.Context) {
	userInfo, exist := base.GetCurrentUserInfo(c)
	if !exist || userInfo.ID == 0 {
		response.ErrorByCode(c, constant.CodeNotLogin)
		return
	}

	cfg, err := global.UserConfigRepo.GetUserConfig(userInfo.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorDataNotFound(c)
			return
		}

		response.ErrorDatabase(c, err.Error())
		return
	}

	response.SuccessData(c, cfg)
}

func (a *UserConfigRouter) SetConfig(c *gin.Context) {
	userInfo, exist := base.GetCurrentUserInfo(c)
	if !exist || userInfo.ID == 0 {
		response.ErrorByCode(c, constant.CodeNotLogin)
		return
	}

	req := repository.UserConfig{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	// Set user ID
	req.UserId = userInfo.ID

	// Save to database
	if err := global.UserConfigRepo.SaveUserConfig(&req); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.Success(c)
}
