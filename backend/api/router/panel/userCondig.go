package panel

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"sun-panel/api/common/apiReturn"
	"sun-panel/api/common/base"
	"sun-panel/api/middleware"
	"sun-panel/internal/common"
	"sun-panel/internal/global"
	"sun-panel/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserConfig struct {
}

func NewUserConfigRouter() *UserConfig {
	return &UserConfig{}
}

func (a *UserConfig) InitRouter(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(middleware.JWTAuth())
	{
		r.POST("/panel/userConfig/set", a.Set)
		r.GET("/panel/userConfig/get", a.Get)
	}
}

func (a *UserConfig) Get(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	cfg := repository.UserConfig{}
	if err := global.Db.First(&cfg, "user_id=?", userInfo.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiReturn.ErrorDataNotFound(c)
			return
		} else {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
	}

	// 处理字段
	if err := json.Unmarshal([]byte(cfg.PanelJson), &cfg.Panel); err != nil {
		cfg.Panel = nil
	}
	if err := json.Unmarshal([]byte(cfg.SearchEngineJson), &cfg.SearchEngine); err != nil {
		cfg.SearchEngine = nil
	}
	apiReturn.SuccessData(c, cfg)
}

func (a *UserConfig) Set(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := repository.UserConfig{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 处理字段
	req.PanelJson = common.ToJSONString(req.Panel)
	req.SearchEngineJson = common.ToJSONString(req.SearchEngine)

	// 保存操作
	if err := global.Db.First(&repository.UserConfig{}, "user_id=?", userInfo.ID).Error; err != nil {
		req.UserId = userInfo.ID
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 新增
			if err := global.Db.Create(&req).Error; err != nil {
				apiReturn.ErrorDatabase(c, err.Error())
				return
			}
		} else {
			// 报错
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
	} else {
		// 修改
		if err := global.Db.Where("user_id=?", userInfo.ID).Updates(&req).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
	}

	apiReturn.Success(c)
}
