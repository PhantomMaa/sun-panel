package system

import (
	"sun-panel/api"
	"sun-panel/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitFileRouter(router *gin.RouterGroup) {
	FileApi := api.ApiGroupApp.ApiSystem.FileApi

	// 公共访问组，不需要 JWT 认证
	public := router.Group("")
	{
		// S3 文件访问路由
		public.GET("/file/s3/*filepath", FileApi.GetS3File)
	}

	// 需要 JWT 认证的私有访问组
	private := router.Group("")
	private.Use(middleware.JWTAuth())
	{
		private.POST("/file/uploadImg", FileApi.UploadImg)
		private.POST("/file/deletes", FileApi.Deletes)
		private.GET("/file/getList", FileApi.GetList)
	}
}
