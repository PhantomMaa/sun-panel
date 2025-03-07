package system

import (
	"github.com/gin-gonic/gin/binding"
	"sun-panel/internal/biz/constant"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/global"
	"sun-panel/internal/util"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func NewUserRouter() *UserRouter {
	return &UserRouter{}
}

func (a *UserRouter) InitRouter(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(interceptor.JWTAuth)
	{
		r.GET("/user/getInfo", a.GetInfo)
		r.POST("/user/updatePassword", a.UpdatePasssword)
		r.POST("/user/updateInfo", a.UpdateInfo)
		r.GET("/user/getAuthInfo", a.GetAuthInfo)
	}
}

func (a *UserRouter) GetInfo(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	response.SuccessData(c, gin.H{
		"userId":    userInfo.ID,
		"id":        userInfo.ID,
		"headImage": userInfo.HeadImage,
		"name":      userInfo.Name,
		"role":      userInfo.Role,
	})
}

func (a *UserRouter) GetAuthInfo(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	visitMode := base.GetCurrentVisitMode(c)
	user := repository.User{}
	user.ID = userInfo.ID
	user.HeadImage = userInfo.HeadImage
	user.Name = userInfo.Name
	user.Role = userInfo.Role
	user.Username = userInfo.Username
	response.SuccessData(c, gin.H{
		"user":      user,
		"visitMode": visitMode,
	})
}

func (a *UserRouter) UpdateInfo(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	type UpdateUserInfoStruct struct {
		HeadImage string `json:"headImage"`
		Name      string `json:"name" validate:"max=15,min=3,required"`
	}
	params := UpdateUserInfoStruct{}

	err := c.ShouldBindBodyWith(&params, binding.JSON)
	if err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(&params); err != nil {
		response.ErrorParamFomat(c, errMsg)
		return
	}

	err = global.UserRepo.UpdateUserInfo(userInfo.ID, map[string]interface{}{
		"head_image": params.HeadImage,
		"name":       params.Name,
	})
	if err != nil {
		response.ErrorDatabase(c, err.Error())
	}
	response.Success(c)
}

func (a *UserRouter) UpdatePasssword(c *gin.Context) {
	type UpdatePasssStruct struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	params := UpdatePasssStruct{}

	err := c.ShouldBindBodyWith(&params, binding.JSON)
	if err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	vUser, err := global.UserRepo.Get(userInfo.ID)
	if err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	if vUser.Password != util.PasswordEncryption(params.OldPassword) {
		// 旧密码不正确
		response.ErrorByCode(c, constant.CodeOldPasswordWrong)
		return
	}

	err = global.UserRepo.UpdateUserInfo(userInfo.ID, map[string]interface{}{
		"password": util.PasswordEncryption(params.NewPassword),
		"token":    "",
	})
	if err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.Success(c)
}

func (a *UserRouter) Logout(c *gin.Context) {
	c.SetCookie("cloud_tk", "", 0, "/source/", "", false, true)
}
