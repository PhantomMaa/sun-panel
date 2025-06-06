package router

import (
	"sun-panel/internal/infra/config"
	"sun-panel/internal/infra/zaplog"
	"sun-panel/internal/web/router/panel"
	"sun-panel/internal/web/router/system"

	"github.com/gin-gonic/gin"
)

type IRouter interface {
	InitRouter(Router *gin.RouterGroup)
}

func RouterArray() []IRouter {
	return []IRouter{
		system.NewAboutRouter(),
		system.NewLoginRouter(),
		system.NewFileRouter(),
		system.NewUserRouter(),
		system.NewModuleConfigRouter(),
		system.NewMonitorRouter(),
		system.NewOAuthRouter(),
		panel.NewItemIconRouter(),
		panel.NewItemIconGroupRouter(),
		panel.NewUserConfigRouter(),
		panel.NewUsersRouter(),
		panel.NewPublicVisitRouter(),
	}
}

func InitRouters(addr string) error {
	router := gin.Default()
	rootRouter := router.Group("/")

	// 注册标准 API 路由
	routerGroup := rootRouter.Group("api")
	for _, router := range RouterArray() {
		router.InitRouter(routerGroup)
	}

	// WEB文件服务
	if config.AppConfig.Base.EnableStaticServer {
		webPath := "./web"

		// 使用StaticFS处理所有静态资源
		router.StaticFS("/assets", gin.Dir(webPath+"/assets", false))
		router.StaticFS("/custom", gin.Dir(webPath+"/custom", false))

		// 处理根目录下的特定文件
		router.StaticFile("/", webPath+"/index.html")

		if config.AppConfig.Rclone.Type == "local" {
			// 使用本次存储时，为本次存储设置静态文件服务
			bucket := config.AppConfig.Rclone.Bucket
			urlPrefix := config.AppConfig.Base.URLPrefix
			router.Static(urlPrefix, bucket)
		}

		zaplog.Logger.Info("Static file server is enabled")
	} else {
		zaplog.Logger.Info("Static file server is disabled")
	}

	zaplog.Logger.Info("Moon-Box is Started.  Listening and serving HTTP on ", addr)
	return router.Run(addr)
}
