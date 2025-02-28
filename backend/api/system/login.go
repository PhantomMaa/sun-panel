package system

import (
	"errors"
	"strings"
	"sun-panel/api/common/apiReturn"
	"sun-panel/api/common/base"
	"sun-panel/internal/common"
	"sun-panel/internal/global"
	"sun-panel/internal/jwt" // 新增jwt包导入
	"sun-panel/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginApi struct {
}

// 登录输入验证
type LoginLoginVerify struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,max=50"`
	VCode    string `json:"vcode" validate:"max=6"`
	Email    string `json:"email"`
}

// @Summary 登录账号
// @Accept application/json
// @Produce application/json
// @Param LoginLoginVerify body LoginLoginVerify true "登陆验证信息"
// @Tags user
// @Router /login [post]
func (l LoginApi) Login(c *gin.Context) {
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
