package interceptor

import (
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

func AdminInterceptor(c *gin.Context) {
	currentUser, exist := base.GetCurrentUserInfo(c)
	if !exist || currentUser.Role != 1 {
		response.ErrorNoAccess(c)
		c.Abort()
		return
	}
}
