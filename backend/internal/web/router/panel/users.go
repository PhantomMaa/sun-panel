package panel

import (
	"strings"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/constant"
	"sun-panel/internal/global"
	"sun-panel/internal/util"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
)

type UsersRouter struct {
}

func NewUsersRouter() *UsersRouter {
	return &UsersRouter{}
}

func (a UsersRouter) InitRouter(router *gin.RouterGroup) {
	rAdmin := router.Group("")
	rAdmin.Use(interceptor.JWTAuth, interceptor.AdminInterceptor)
	{
		rAdmin.POST("panel/users/create", a.Create)
		rAdmin.GET("panel/users/getList", a.GetList)
		rAdmin.POST("panel/users/update", a.Update)
		rAdmin.POST("panel/users/deletes", a.Deletes)
	}
}

func (a UsersRouter) Create(c *gin.Context) {
	param := repository.User{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		response.ErrorParamFomat(c, errMsg)
		return
	}

	param.Username = strings.TrimSpace(param.Username)
	if len(param.Username) < 5 {
		response.ErrorParamFomat(c, "The account must be no less than 5 characters long")
		return
	}

	// 验证账号是否存在
	if _, err := global.UserRepo.CheckUsernameExist(param.Username, ""); err != nil {
		response.ErrorByCode(c, constant.CodeAccountAlreadyExist)
		return
	}

	mUser := repository.User{
		Username:      strings.TrimSpace(param.Username),
		Password:      util.PasswordEncryption(param.Password),
		Name:          param.Name,
		HeadImage:     param.HeadImage,
		Status:        1,
		Role:          param.Role,
		OauthProvider: constant.OAuthProviderBuildin,
	}
	err := global.UserService.CreateUser(&mUser)
	if err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.SuccessData(c, gin.H{"userId": mUser.ID})
}

func (a UsersRouter) Deletes(c *gin.Context) {
	type Param struct {
		UserIds []uint
	}
	param := Param{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		c.Abort()
		return
	}

	if err := global.UserRepo.Deletes(param.UserIds); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.Success(c)
}

func (a UsersRouter) Update(c *gin.Context) {
	user := repository.User{}
	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		c.Abort()
		return
	}

	if user.Password == "" {
		user.Password = "-" // 修改不允许修改密码，为了验证通过
	}

	if errMsg, err := base.ValidateInputStruct(user); err != nil {
		response.ErrorParamFomat(c, errMsg)
		return
	}

	user.Username = strings.Trim(user.Username, " ")
	if len(user.Username) < 3 {
		// 账号不得少于3个字符
		response.ErrorParamFomat(c, "The account must be no less than 3 characters long")
		return
	}

	allowField := []string{"Username", "Name", "Mail", "Token", "Role"}

	// 密码不为默认“-”空，修改密码
	if user.Password != "-" {
		user.Password = util.PasswordEncryption(user.Password)
		allowField = append(allowField, "Password")
	}

	// 验证账号是否存在
	_, err := global.UserRepo.Get(user.ID)
	if err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	user.Token = "" // 修改资料就重置token
	if err := global.UserRepo.Update(user.ID, &user); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	// 返回token等基本信息
	response.SuccessData(c, user)
}

func (a UsersRouter) GetList(c *gin.Context) {
	type ParamsStruct struct {
		repository.PagedParam
	}

	param := ParamsStruct{}
	if err := c.ShouldBind(&param); err != nil {
		response.ErrorParamFomat(c, err.Error())
		c.Abort()
		return
	}

	list, count, err := global.UserRepo.GetList(param.PagedParam)
	if err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.SuccessListData(c, list, count)
}
