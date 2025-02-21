package middleware

import (
	"strings"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/lib/jwt"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

// LoginInterceptor JWT认证中间件
func LoginInterceptor(c *gin.Context) {
	// 获取Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		apiReturn.ErrorByCode(c, 1000) // 未登录
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
	mUser := models.User{}
	userInfo, err := mUser.GetUserInfoByUid(claims.UserID)
	if err != nil {
		apiReturn.Error(c, "用户不存在")
		c.Abort()
		return
	}

	// 检查用户状态
	if userInfo.Status != 1 {
		apiReturn.ErrorByCode(c, 1004) // 用户已禁用
		c.Abort()
		return
	}

	// 将用户信息存储到上下文
	c.Set("userInfo", userInfo)
	c.Next()
}

// 不验证缓存直接验证库省去没有缓存每次都要手动登录的问题
func LoginInterceptorDev(c *gin.Context) {

	// 获得token
	token := c.GetHeader("token")
	mUser := models.User{}

	// 去库中查询是否存在该用户；否则返回错误
	if info, err := mUser.GetUserInfoByToken(token); err != nil || info.ID == 0 {
		apiReturn.ErrorCode(c, 1001, "login.err_token_expire", nil)
		c.Abort()
		return
	} else {
		// 通过
		// 设置当前用户信息
		c.Set("userInfo", info)
	}
}
