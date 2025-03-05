package panel

import (
	"github.com/gin-gonic/gin/binding"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

type UserConfig struct {
}

func NewUserConfigRouter() *UserConfig {
	return &UserConfig{}
}

func (a *UserConfig) InitRouter(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(interceptor.JWTAuth)
	{
		r.POST("/panel/userConfig/set", a.Set)
		r.GET("/panel/userConfig/get", a.Get)
	}
}

func (a *UserConfig) Get(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	
	cfg, err := repository.GetUserConfig(userInfo.ID)
	if err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.SuccessData(c, cfg)
}

func (a *UserConfig) Set(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := repository.UserConfig{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	// Set user ID
	req.UserId = userInfo.ID

	// Save to database
	if err := repository.SaveUserConfig(&req); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.Success(c)
}
