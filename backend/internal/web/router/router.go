package router

import (
	"sun-panel/internal/global"
	"sun-panel/internal/infra/config"
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
		panel.NewItemIconRouter(),
		panel.NewItemIconGroupRouter(),
		panel.NewUserConfigRouter(),
		panel.NewUsersRouter(),
	}
}

func InitRouters(addr string) error {
	router := gin.Default()
	rootRouter := router.Group("/")
	routerGroup := rootRouter.Group("api")

	for _, router := range RouterArray() {
		router.InitRouter(routerGroup)
	}

	// WEB文件服务
	if config.AppConfig.Base.EnableStaticServer {
		webPath := "./web"
		router.Static("/assets", webPath+"/assets")
		router.Static("/custom", webPath+"/custom")
		router.StaticFile("/", webPath+"/index.html")
		router.StaticFile("/favicon.ico", webPath+"/favicon.ico")
		router.StaticFile("/favicon.svg", webPath+"/favicon.svg")

		if config.AppConfig.Rclone.Type == "local" {
			// 使用本次存储时，为本次存储设置静态文件服务
			bucket := config.AppConfig.Rclone.Bucket
			urlPrefix := config.AppConfig.Base.URLPrefix
			router.Static(urlPrefix, bucket)
		}

		global.Logger.Info("Static file server is enabled")
	} else {
		global.Logger.Info("Static file server is disabled")
	}

	global.Logger.Info("Sun-Panel is Started.  Listening and serving HTTP on ", addr)
	return router.Run(addr)
}
