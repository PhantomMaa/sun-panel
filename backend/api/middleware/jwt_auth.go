package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"sun-panel/api/common/apiReturn"
	"sun-panel/internal/jwt"
	"sun-panel/internal/repository"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apiReturn.ErrorByCode(c, 401)
			c.Abort()
			return
		}

		// 支持Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			authHeader = parts[1]
		}

		// 解析Token
		claims, err := jwt.ParseToken(authHeader)
		if err != nil {
			apiReturn.Error(c, "无效的访问凭证")
			c.Abort()
			return
		}

		// 获取用户信息
		mUser := repository.User{}
		userInfo, err := mUser.GetUserInfoByUid(claims.UserID)
		if err != nil {
			apiReturn.Error(c, "用户不存在")
			c.Abort()
			return
		}

		// 检查用户状态
		if userInfo.Status != 1 {
			apiReturn.ErrorByCode(c, 1004)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("userInfo", userInfo)
		c.Next()
	}
}
