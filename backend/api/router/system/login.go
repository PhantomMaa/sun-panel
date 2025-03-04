package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"sun-panel/api/common/apiReturn"
	"sun-panel/api/common/base"
	"sun-panel/api/middleware"
	"sun-panel/internal/common"
	"sun-panel/internal/global"
	"sun-panel/internal/jwt"
	"sun-panel/internal/repository"
)

type LoginApi struct {
}

type LoginLoginVerify struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,max=50"`
	VCode    string `json:"vcode" validate:"max=6"`
	Email    string `json:"email"`
}

func NewLoginRouter() *LoginApi {
	return &LoginApi{}
}

func (l *LoginApi) InitRouter(router *gin.RouterGroup) {
	// 公开接口
	router.POST("/login", l.Login)

	// 需要认证的接口
	authGroup := router.Group("")
	authGroup.Use(middleware.JWTAuth())
	{
		authGroup.POST("/logout", l.Logout)
	}
}

func (l *LoginApi) Login(c *gin.Context) {
	param := LoginLoginVerify{}
	if err := c.ShouldBindJSON(&param); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	mUser := repository.User{}
	var (
		err  error
		info repository.User
	)

	param.Username = strings.TrimSpace(param.Username)
	if info, err = mUser.GetUserInfoByUsernameAndPassword(param.Username, common.PasswordEncryption(param.Password)); err != nil {
		// 未找到记录 账号或密码错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiReturn.ErrorByCode(c, 1003)
			return
		} else {
			// 未知错误
			apiReturn.Error(c, err.Error())
			return
		}
	}

	// 停用或未激活
	if info.Status != 1 {
		apiReturn.ErrorByCode(c, 1004)
		return
	}

	// 生成JWT Token
	tokenString, err := jwt.GenerateToken(info.ID)
	if err != nil {
		global.Logger.Error("JWT生成失败:", err)
		apiReturn.Error(c, "系统错误")
		return
	}

	info.Password = "" // 清除敏感信息
	info.Token = tokenString

	// 设置当前用户信息
	c.Set("userInfo", info)
	apiReturn.SuccessData(c, info)
}

// Logout 安全退出
func (l *LoginApi) Logout(c *gin.Context) {
	apiReturn.Success(c)
}
