package panel

import (
	"github.com/gin-gonic/gin/binding"
	"strings"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/common"
	"sun-panel/internal/global"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"

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

	mUser := repository.User{
		Username:  strings.TrimSpace(param.Username),
		Password:  common.PasswordEncryption(param.Password),
		Name:      param.Name,
		HeadImage: param.HeadImage,
		Status:    1,
		Role:      param.Role,
	}

	// 验证账号是否存在
	if _, err := mUser.CheckUsernameExist(param.Username); err != nil {
		response.ErrorByCode(c, 1008)
		return
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
	param := repository.User{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		c.Abort()
		return
	}

	if param.Password == "" {
		param.Password = "-" // 修改不允许修改密码，为了验证通过
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		response.ErrorParamFomat(c, errMsg)
		return
	}

	param.Username = strings.Trim(param.Username, " ")
	if len(param.Username) < 3 {
		// 账号不得少于3个字符
		response.ErrorParamFomat(c, "The account must be no less than 3 characters long")
		return
	}

	allowField := []string{"Username", "Name", "Mail", "Token", "Role"}

	// 密码不为默认“-”空，修改密码
	if param.Password != "-" {
		param.Password = common.PasswordEncryption(param.Password)
		allowField = append(allowField, "Password")
	}

	// 验证账号是否存在
	mUser := repository.User{}
	_, err := mUser.GetUserInfoByUid(param.ID)
	if err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	param.Token = "" // 修改资料就重置token
	if err := global.Db.Select(allowField).Where("id=?", param.ID).Updates(&param).Error; err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	// 返回token等基本信息
	response.SuccessData(c, param)
}

func (a UsersRouter) GetList(c *gin.Context) {
	type ParamsStruct struct {
		repository.User
		Limit   int    `form:"limit" json:"limit"`
		Page    int    `form:"page" json:"page"`
		Keyword string `form:"keyword" json:"keyword"`
	}

	var (
		list  []repository.User
		count int64
	)

	db := global.Db

	// 查询条件
	param := ParamsStruct{}
	if err := c.ShouldBind(&param); err != nil {
		response.ErrorParamFomat(c, err.Error())
		c.Abort()
		return
	}

	if param.Keyword != "" {
		db = db.Where("name LIKE ? OR username LIKE ?", "%"+param.Keyword+"%", "%"+param.Keyword+"%")
	}

	if err := db.Omit("Password").Limit(param.Limit).Offset((param.Page - 1) * param.Limit).Find(&list).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.SuccessListData(c, list, uint(count))
}
