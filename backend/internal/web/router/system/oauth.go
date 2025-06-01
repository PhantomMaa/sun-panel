package system

import (
	"strings"
	"sun-panel/internal/global"
	"sun-panel/internal/infra/config"
	"sun-panel/internal/infra/zaplog"
	"sun-panel/internal/util"
	"sun-panel/internal/util/jwt"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

type OAuthRouter struct {
}

func NewOAuthRouter() *OAuthRouter {
	return &OAuthRouter{}
}

func (r *OAuthRouter) InitRouter(router *gin.RouterGroup) {
	// 公开接口
	router.GET("/oauth/config", r.GetConfig)

	// OAuth相关接口
	if config.AppConfig.OAuth.Enable {
		router.GET("/oauth/:provider", r.OAuthLogin)
		router.GET("/oauth/:provider/callback", r.OAuthCallback)
	}
}

// GetConfig 获取OAuth配置信息
func (r *OAuthRouter) GetConfig(c *gin.Context) {
	// 构建返回数据
	data := map[string]any{
		"enabled":   config.AppConfig.OAuth.Enable,
		"providers": []string{},
	}

	// 检查哪些提供商已配置
	var providers []string
	for _, provider := range config.AppConfig.OAuth.Providers {
		providers = append(providers, provider.Name)
	}
	data["providers"] = providers

	response.SuccessData(c, data)
}

// OAuthLogin 处理OAuth登录请求
func (r *OAuthRouter) OAuthLogin(c *gin.Context) {
	provider := c.Param("provider")

	// 检查是否支持该OAuth提供商
	if !r.isProviderSupported(provider) {
		response.Error(c, "不支持的OAuth提供商")
		return
	}

	// 获取OAuth授权URL
	authURL, err := global.UserService.GetOAuthLoginURL(provider, util.RedirectURL(config.AppConfig.Base.RootURL, provider))
	if err != nil {
		zaplog.Logger.Error("获取OAuth授权URL失败:", err)
		response.Error(c, "获取OAuth授权URL失败")
		return
	}

	// 重定向到OAuth授权页面
	c.Redirect(302, authURL)
}

// OAuthCallback 处理OAuth回调
func (r *OAuthRouter) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")

	// 检查是否支持该OAuth提供商
	if !r.isProviderSupported(provider) {
		response.Error(c, "不支持的OAuth提供商")
		return
	}

	// 处理OAuth回调
	user, err := global.UserService.HandleOAuthCallback(provider, code, util.RedirectURL(config.AppConfig.Base.RootURL, provider))
	if err != nil {
		zaplog.Logger.Error("处理OAuth回调失败:", err)
		response.Error(c, "处理OAuth回调失败")
		return
	}

	// 生成JWT Token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		zaplog.Logger.Error("生成token失败:", err)
		response.Error(c, "生成token失败")
		return
	}

	redirectUrl := config.AppConfig.Base.RootURL + "/login?token=" + token
	c.Redirect(302, redirectUrl)
}

// 检查是否支持该OAuth提供商
func (r *OAuthRouter) isProviderSupported(provider string) bool {
	for _, p := range config.AppConfig.OAuth.Providers {
		if strings.EqualFold(p.Name, provider) {
			return true
		}
	}
	return false
}
