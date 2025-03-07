package interceptor

import (
	"strings"
	"sun-panel/internal/biz/constant"
	"sun-panel/internal/global"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.ErrorByCode(c, constant.CodeNotLogin)
		c.Abort()
		return
	}

	// 支持Bearer token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" {
		authHeader = parts[1]
	}

	// 解析Token
	claims, err := ParseToken(authHeader)
	if err != nil {
		global.Logger.Infof("invalid token. %v", err)
		response.ErrorByCode(c, constant.CodeNotLogin)
		c.Abort()
		return
	}

	// 获取用户信息
	userInfo, err := global.UserRepo.Get(claims.UserID)
	if err != nil {
		global.Logger.Infof("user not exist. %v", err)
		response.ErrorByCode(c, constant.CodeNotLogin)
		c.Abort()
		return
	}

	// 检查用户状态
	if userInfo.Status != 1 {
		response.ErrorByCode(c, constant.CodeStatusError)
		c.Abort()
		return
	}

	// 将用户信息存储到上下文
	c.Set("userInfo", userInfo)
	c.Next()
}
