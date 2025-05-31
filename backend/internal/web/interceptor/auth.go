package interceptor

import (
	"errors"
	"strings"
	"sun-panel/internal/constant"
	"sun-panel/internal/global"
	"sun-panel/internal/infra/zaplog"
	"sun-panel/internal/util/jwt"
	"sun-panel/internal/util/publiccode"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

// Auth 认证中间件
func Auth(c *gin.Context) {
	publiccode := c.GetHeader("publiccode")
	var userId uint
	var err error
	var logined bool = false
	if publiccode != "" {
		userId, err = ParseUserIdFromPubliccode(publiccode)
	} else {
		userId, err = ParseUserIdFromJwtToken(c.GetHeader("Authorization"))
		logined = true
	}

	if err != nil {
		response.ErrorByCode(c, constant.CodeNotLogin)
		c.Abort()
		return
	}

	// 获取用户信息
	user, err := global.UserRepo.Get(userId)
	if err != nil {
		zaplog.Logger.Infof("user not exist. %v", err)
		response.ErrorByCode(c, constant.CodeNotLogin)
		c.Abort()
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		response.ErrorByCode(c, constant.CodeStatusError)
		c.Abort()
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
		Logined:    logined,
	}
	c.Set("userInfo", userInfo)
	c.Next()
}

// ParseUserIdFromJwtToken 解析JWT Token，获取用户ID
func ParseUserIdFromJwtToken(authHeader string) (uint, error) {
	if authHeader == "" {
		return 0, errors.New("authHeader is empty")
	}

	// 支持Bearer token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" {
		authHeader = parts[1]
	}

	// 解析Token
	claims, err := jwt.ParseToken(authHeader)
	if err != nil {
		zaplog.Logger.Infof("invalid token. %v", err)
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}

// ParseUserIdFromPubliccode 解析公开访问代码，获取用户ID
func ParseUserIdFromPubliccode(code string) (uint, error) {
	// 解析公开访问代码，获取用户ID
	userID, err := publiccode.ParseCode(code)
	if err != nil {
		zaplog.Logger.Infof("invalid public access code. %v", err)
		return 0, errors.New("invalid public access code")
	}

	return userID, nil
}
