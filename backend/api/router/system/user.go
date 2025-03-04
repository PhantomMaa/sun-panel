package system

import (
	"github.com/gin-gonic/gin/binding"
	"sun-panel/api/common/apiReturn"
	"sun-panel/api/common/base"
	"sun-panel/api/middleware"
	"sun-panel/internal/common"
	"sun-panel/internal/global"
	"sun-panel/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserApi struct{}

func NewUserRouter() *UserApi {
	return &UserApi{}
}

func (a *UserApi) InitRouter(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(middleware.JWTAuth())
	{
		r.GET("/user/getInfo", a.GetInfo)
		r.POST("/user/updatePassword", a.UpdatePasssword)
		r.POST("/user/updateInfo", a.UpdateInfo)
		r.GET("/user/getAuthInfo", a.GetAuthInfo)
	}
}

func (a *UserApi) GetInfo(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	apiReturn.SuccessData(c, gin.H{
		"userId":    userInfo.ID,
		"id":        userInfo.ID,
		"headImage": userInfo.HeadImage,
		"name":      userInfo.Name,
		"role":      userInfo.Role,
	})
}

func (a *UserApi) GetAuthInfo(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	visitMode := base.GetCurrentVisitMode(c)
	user := repository.User{}
	user.ID = userInfo.ID
	user.HeadImage = userInfo.HeadImage
	user.Name = userInfo.Name
	user.Role = userInfo.Role
	user.Username = userInfo.Username
	apiReturn.SuccessData(c, gin.H{
		"user":      user,
		"visitMode": visitMode,
	})
}

func (a *UserApi) UpdateInfo(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	type UpdateUserInfoStruct struct {
		HeadImage string `json:"headImage"`
		Name      string `json:"name" validate:"max=15,min=3,required"`
	}
	params := UpdateUserInfoStruct{}

	err := c.ShouldBindBodyWith(&params, binding.JSON)
	if err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(&params); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	mUser := repository.User{}
	err = mUser.UpdateUserInfoByUserId(userInfo.ID, map[string]interface{}{
		"head_image": params.HeadImage,
		"name":       params.Name,
	})
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
	}
	apiReturn.Success(c)
}

func (a *UserApi) UpdatePasssword(c *gin.Context) {
	type UpdatePasssStruct struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	params := UpdatePasssStruct{}

	err := c.ShouldBindBodyWith(&params, binding.JSON)
	if err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	mUser := repository.User{}
	vUser, err := mUser.GetUserInfoByUid(userInfo.ID)
	if err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if vUser.Password != common.PasswordEncryption(params.OldPassword) {
		// 旧密码不正确
		apiReturn.ErrorByCode(c, 1007)
		return
	}

	res := global.Db.Model(&repository.User{}).Where("id", userInfo.ID).Updates(map[string]interface{}{
		"password": common.PasswordEncryption(params.NewPassword),
		"token":    "",
	})
	if res.Error != nil {
		apiReturn.ErrorDatabase(c, res.Error.Error())
		return
	}

	apiReturn.Success(c)
}
