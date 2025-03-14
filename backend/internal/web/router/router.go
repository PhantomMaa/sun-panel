package router

import (
	"github.com/gin-gonic/gin"
	"sun-panel/internal/global"
	"sun-panel/internal/web/router/panel"
	"sun-panel/internal/web/router/system"
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
	if global.Config.GetValueString("base", "enable_static_server") == "true" {
		webPath := "./web"
		
		// 添加缓存头的中间件
		cacheMiddleware := func(c *gin.Context) {
			c.Header("Cache-Control", "public, max-age=86400")
			c.Next()
		}
		
		// 为静态资源路由添加缓存中间件
		router.Group("/assets", cacheMiddleware).Static("", webPath+"/assets")
		router.Group("/custom", cacheMiddleware).Static("", webPath+"/custom")
		
		// 为静态文件添加缓存中间件
		router.GET("/", cacheMiddleware, func(c *gin.Context) {
			c.File(webPath + "/index.html")
		})
		router.GET("/favicon.ico", cacheMiddleware, func(c *gin.Context) {
			c.File(webPath + "/favicon.ico")
		})
		router.GET("/favicon.svg", cacheMiddleware, func(c *gin.Context) {
			c.File(webPath + "/favicon.svg")
		})

		if global.Config.GetValueString("rclone", "type") == "local" {
			// 使用本次存储时，为本次存储设置静态文件服务
			bucket := global.Config.GetValueString("rclone", "bucket")
			urlPrefix := global.Config.GetValueString("base", "url_prefix")
			router.Group(urlPrefix, cacheMiddleware).Static("", bucket)
		}

		global.Logger.Info("Static file server is enabled with 1-day cache")
	} else {
		global.Logger.Info("Static file server is disabled")
	}

	global.Logger.Info("Sun-Panel is Started.  Listening and serving HTTP on ", addr)
	return router.Run(addr)
}
