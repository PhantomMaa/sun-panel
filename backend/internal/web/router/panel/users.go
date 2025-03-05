package panel

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"strings"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/common"
	"sun-panel/internal/global"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

type UsersApi struct {
}

var (
	ErrUsersApiAtLeastKeepOne = errors.New("at least keep one")
)

func NewUsersRouter() *UsersApi {
	return &UsersApi{}
}

func (a UsersApi) InitRouter(router *gin.RouterGroup) {
	rAdmin := router.Group("")
	rAdmin.Use(interceptor.JWTAuth, interceptor.AdminInterceptor)
	{
		rAdmin.POST("panel/users/create", a.Create)
		rAdmin.GET("panel/users/getList", a.GetList)
		rAdmin.POST("panel/users/update", a.Update)
		rAdmin.POST("panel/users/deletes", a.Deletes)
	}
}

func (a UsersApi) Create(c *gin.Context) {
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

	userInfo, err := mUser.CreateOne()

	if err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.SuccessData(c, gin.H{"userId": userInfo.ID})
}

func (a UsersApi) Deletes(c *gin.Context) {
	type UserIds struct {
		UserIds []uint
	}
	param := UserIds{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		c.Abort()
		return
	}

	txErr := global.Db.Transaction(func(tx *gorm.DB) error {
		mitemIconGroup := repository.ItemIconGroup{}

		for _, v := range param.UserIds {
			// 删除图标
			if err := tx.Delete(&repository.ItemIcon{}, "user_id=?", v).Error; err != nil {
				return err
			}

			// 删除分组
			if err := mitemIconGroup.DeleteByUserId(tx, v); err != nil {
				return err
			}

			// 删除模块配置
			if err := tx.Delete(&repository.ModuleConfig{}, "user_id=?", v).Error; err != nil {
				return err
			}

			// 删除用户配置
			if err := tx.Delete(&repository.ModuleConfig{}, "user_id=?", v).Error; err != nil {
				return err
			}

			// 删除文件记录（不删除资源文件）
			if err := tx.Delete(&repository.File{}, "user_id=?", v).Error; err != nil {
				return err
			}
		}

		if err := tx.Delete(&repository.User{}, &param.UserIds).Error; err != nil {
			response.ErrorDatabase(c, err.Error())
			return err
		}

		// 验证是否还存在管理员
		var count int64
		if err := tx.Model(&repository.User{}).Where("role=?", 1).Count(&count).Error; err != nil {
			return err
		} else if count == 0 {
			return ErrUsersApiAtLeastKeepOne
		}

		return nil
	})

	if errors.Is(txErr, ErrUsersApiAtLeastKeepOne) {
		response.ErrorByCode(c, 1201)
		return
	} else if txErr != nil {
		response.ErrorDatabase(c, txErr.Error())
		return
	}

	response.Success(c)
}

func (a UsersApi) Update(c *gin.Context) {
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

func (a UsersApi) GetList(c *gin.Context) {
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

	response.SuccessListData(c, list, count)
}
