package system

import (
	"errors"
	"strings"
	"sun-panel/internal/constant"
	"sun-panel/internal/global"
	"sun-panel/internal/infra/zaplog"
	"sun-panel/internal/util"
	"sun-panel/internal/util/jwt"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	r := router.Group("")
	r.Use(interceptor.Auth)
	{
		r.POST("/logout", l.Logout)
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

	param.Username = strings.TrimSpace(param.Username)
	user, err := global.UserRepo.GetByUsernameAndPassword(param.Username, util.PasswordEncryption(param.Password), constant.OAuthProviderBuildin)
	if err != nil {
		// 未找到记录 账号或密码错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorByCode(c, constant.CodePasswordWrong)
			return
		} else {
			// 未知错误
			response.Error(c, err.Error())
			return
		}
	}

	// 停用或未激活
	if user.Status != 1 {
		response.ErrorByCode(c, constant.CodeStatusError)
		return
	}

	// 生成JWT Token
	user.Token, err = jwt.GenerateToken(user.ID)
	if err != nil {
		zaplog.Logger.Error("JWT生成失败:", err)
		response.Error(c, "系统错误")
		return
	}

	// 将用户信息存储到上下文
	userInfo := base.UserInfo{
		ID:         user.ID,
		Name:       user.Name,
		Role:       user.Role,
		Username:   user.Username,
		Publiccode: user.Publiccode,
		Token:      user.Token,
	}
	response.SuccessData(c, userInfo)
}

// Logout 安全退出
func (l *LoginRouter) Logout(c *gin.Context) {
	response.Success(c)
}
