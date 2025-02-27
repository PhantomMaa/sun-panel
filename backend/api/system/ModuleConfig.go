package system

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"sun-panel/api/common/apiReturn"
	"sun-panel/api/common/base"
	"sun-panel/internal/global"
	"sun-panel/internal/repository"
)

type ModuleConfigApi struct{}

func (a *ModuleConfigApi) GetByName(c *gin.Context) {
	req := repository.ModuleConfig{}

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	mCfg := repository.ModuleConfig{}
	if cfg, err := mCfg.GetConfigByUserIdAndName(global.Db, userInfo.ID, req.Name); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	} else {
		apiReturn.SuccessData(c, cfg)
		return
	}

}

func (a *ModuleConfigApi) Save(c *gin.Context) {
	req := repository.ModuleConfig{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	userInfo, _ := base.GetCurrentUserInfo(c)
	mCfg := repository.ModuleConfig{}
	mCfg.UserId = userInfo.ID
	mCfg.Value = req.Value
	mCfg.Name = req.Name

	if err := mCfg.Save(global.Db); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}
