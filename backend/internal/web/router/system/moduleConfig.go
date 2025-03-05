package system

import (
	"github.com/gin-gonic/gin/binding"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/global"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

type ModuleConfigApi struct{}

func NewModuleConfigRouter() *ModuleConfigApi {
	return &ModuleConfigApi{}
}

func (a *ModuleConfigApi) InitRouter(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(interceptor.JWTAuth)
	{
		r.GET("/system/moduleConfig/getByName", a.GetByName)
		r.POST("/system/moduleConfig/save", a.Save)
	}
}

func (a *ModuleConfigApi) GetByName(c *gin.Context) {
	req := repository.ModuleConfig{}

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	mCfg := repository.ModuleConfig{}
	if cfg, err := mCfg.GetConfigByUserIdAndName(global.Db, userInfo.ID, req.Name); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	} else {
		response.SuccessData(c, cfg)
		return
	}

}

func (a *ModuleConfigApi) Save(c *gin.Context) {
	req := repository.ModuleConfig{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}
	userInfo, _ := base.GetCurrentUserInfo(c)
	mCfg := repository.ModuleConfig{}
	mCfg.UserId = userInfo.ID
	mCfg.Value = req.Value
	mCfg.Name = req.Name

	if err := mCfg.Save(global.Db); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}
	response.Success(c)
}
