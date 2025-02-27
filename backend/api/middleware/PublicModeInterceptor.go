package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"sun-panel/api/common/apiReturn"
	"sun-panel/api/common/base"
	"sun-panel/internal/common/systemSetting"
	"sun-panel/internal/global"
	"sun-panel/internal/jwt"
	"sun-panel/internal/language"
	"sun-panel/internal/repository"
)

// 公开访问模式（访客模式）
// [有token将自动登录，无token/过期将使用公开账号，不可以与LoginInterceptor一起使用]
func PublicModeInterceptor(c *gin.Context) {
	// 获取Authorization header
	authHeader := c.GetHeader("Authorization")
	var userInfo repository.User

	// 尝试JWT验证
	if authHeader != "" {
		// 支持Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			authHeader = parts[1]
		}

		// 解析Token
		if claims, err := jwt.ParseToken(authHeader); err == nil {
			// Token有效，获取用户信息
			mUser := repository.User{}
			if info, err := mUser.GetUserInfoByUid(claims.UserID); err == nil {
				userInfo = info
				c.Set("userInfo", userInfo)
				return
			}
		}
	}

	// 如果JWT验证失败或没有token，使用公开账号
	var userId *uint
	if err := global.SystemSetting.GetValueByInterface(systemSetting.PanelPublicUserId, &userId); err == nil && userId != nil {
		if err := global.Db.First(&userInfo, "id=?", userId).Error; err != nil {
			apiReturn.ErrorCode(c, 1001, language.Obj.Get("login.err_token_expire"), nil)
			c.Abort()
			return
		}
		global.Logger.Debug("访客用户ID:", userInfo.ID)
		c.Set("userInfo", userInfo)
		c.Set(base.GIN_GET_VISIT_MODE, base.VISIT_MODE_PUBLIC)
		return
	}

	global.Logger.Debug("访客用户不存在:", userId)
	apiReturn.ErrorCode(c, 1001, language.Obj.Get("login.err_token_expire"), nil)
	c.Abort()
	return
}
