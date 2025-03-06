package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/common"
	"sun-panel/internal/global"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"
)

type LoginRouter struct {
}

type LoginLoginVerify struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,max=50"`
	VCode    string `json:"vcode" validate:"max=6"`
	Email    string `json:"email"`
}

func NewLoginRouter() *LoginRouter {
	return &LoginRouter{}
}

func (l *LoginRouter) InitRouter(router *gin.RouterGroup) {
	// 公开接口
	router.POST("/login", l.Login)

	// 需要认证的接口
	authGroup := router.Group("")
	authGroup.Use(interceptor.JWTAuth)
	{
		authGroup.POST("/logout", l.Logout)
	}
}

func (l *LoginRouter) Login(c *gin.Context) {
	param := LoginLoginVerify{}
	if err := c.ShouldBindJSON(&param); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		response.ErrorParamFomat(c, errMsg)
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
			response.ErrorByCode(c, 1003)
			return
		} else {
			// 未知错误
			response.Error(c, err.Error())
			return
		}
	}

	// 停用或未激活
	if info.Status != 1 {
		response.ErrorByCode(c, 1004)
		return
	}

	// 生成JWT Token
	tokenString, err := interceptor.GenerateToken(info.ID)
	if err != nil {
		global.Logger.Error("JWT生成失败:", err)
		response.Error(c, "系统错误")
		return
	}

	info.Password = "" // 清除敏感信息
	info.Token = tokenString

	// 设置当前用户信息
	c.Set("userInfo", info)
	response.SuccessData(c, info)
}

// Logout 安全退出
func (l *LoginRouter) Logout(c *gin.Context) {
	response.Success(c)
}
