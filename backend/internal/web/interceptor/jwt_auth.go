package interceptor

import (
	"strings"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.ErrorByCode(c, 1001)
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
		response.Error(c, "无效的访问凭证")
		c.Abort()
		return
	}

	// 获取用户信息
	mUser := repository.User{}
	userInfo, err := mUser.GetUserInfoByUid(claims.UserID)
	if err != nil {
		response.Error(c, "用户不存在")
		c.Abort()
		return
	}

	// 检查用户状态
	if userInfo.Status != 1 {
		response.ErrorByCode(c, 1004)
		c.Abort()
		return
	}

	// 将用户信息存储到上下文
	c.Set("userInfo", userInfo)
	c.Next()
}
