package interceptor

import (
	"sun-panel/internal/constant"
	"sun-panel/internal/global"
	"sun-panel/internal/infra/zaplog"
	"sun-panel/internal/util/publiccode"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

// PublicAccess 公开访问认证中间件
func PublicAccess(c *gin.Context) {
	zaplog.Logger.Infof("public access. %v", c.Request.URL.Path)

	// 直接从路由参数中获取公开访问代码
	code := c.Param("code")
	zaplog.Logger.Infof("public access code: %s", code)
	
	if code == "" {
		zaplog.Logger.Infof("empty public access code")
		response.ErrorByCode(c, constant.CodeNotLogin)
		c.Abort()
		return
	}

	// 解析公开访问代码，获取用户ID
	userID, err := publiccode.ParseCode(code)
	if err != nil {
		zaplog.Logger.Infof("invalid public access code. %v", err)
		response.ErrorByCode(c, constant.CodeNotLogin)
		c.Abort()
		return
	}

	// 获取用户信息
	userInfo, err := global.UserRepo.Get(userID)
	if err != nil {
		zaplog.Logger.Infof("user not exist. %v", err)
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
