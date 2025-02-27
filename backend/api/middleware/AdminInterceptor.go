package middleware

import (
	"sun-panel/api/common/apiReturn"
	"sun-panel/api/common/base"

	"github.com/gin-gonic/gin"
)

func AdminInterceptor(c *gin.Context) {
	currentUser, _ := base.GetCurrentUserInfo(c)
	if currentUser.Role != 1 {
		apiReturn.ErrorNoAccess(c)
		c.Abort()
		return
	}
}
