package router

import (
	"github.com/gin-gonic/gin"
	"sun-panel/internal/global"
	"sun-panel/internal/router/panel"
	"sun-panel/internal/router/system"
)

func InitRouters(addr string) error {
	router := gin.Default()
	rootRouter := router.Group("/")
	routerGroup := rootRouter.Group("api")

	// 接口
	system.Init(routerGroup)
	panel.Init(routerGroup)

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
