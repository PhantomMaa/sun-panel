package router

import "github.com/gin-gonic/gin"

type IRouter interface {
	InitRouter(Router *gin.RouterGroup)
}
