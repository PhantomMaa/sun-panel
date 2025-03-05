package router

import (
	"github.com/gin-gonic/gin"
	"sun-panel/internal/global"
	panel2 "sun-panel/internal/web/router/panel"
	system2 "sun-panel/internal/web/router/system"
)

type IRouter interface {
	InitRouter(Router *gin.RouterGroup)
}

func RouterArray() []IRouter {
	return []IRouter{
		system2.NewAboutRouter(),
		system2.NewLoginRouter(),
		system2.NewFileRouter(),
		system2.NewUserRouter(),
		system2.NewModuleConfigRouter(),
		system2.NewMonitorRouter(),
		system2.NewNoticeRouter(),
		panel2.NewItemIconRouter(),
		panel2.NewItemIconGroupRouter(),
		panel2.NewUserConfigRouter(),
		panel2.NewUsersRouter(),
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
	if global.Config.GetValueString("base", "enable_static_server") == "true" {
		webPath := "./web"
		router.Static("/assets", webPath+"/assets")
		router.Static("/custom", webPath+"/custom")
		router.StaticFile("/", webPath+"/index.html")
		router.StaticFile("/favicon.ico", webPath+"/favicon.ico")
		router.StaticFile("/favicon.svg", webPath+"/favicon.svg")

		if global.Config.GetValueString("rclone", "type") == "local" {
			// 使用本次存储时，为本次存储设置静态文件服务
			bucket := global.Config.GetValueString("rclone", "bucket")
			urlPrefix := global.Config.GetValueString("base", "url_prefix")
			router.Static(urlPrefix, bucket)
		}

		global.Logger.Info("Static file server is enabled")
	} else {
		global.Logger.Info("Static file server is disabled")
	}

	global.Logger.Info("Sun-Panel is Started.  Listening and serving HTTP on ", addr)
	return router.Run(addr)
}
