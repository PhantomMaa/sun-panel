package repository

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Username      string `gorm:"type:varchar(255);uniqueIndex:idx_username_oauth_provider" json:"username"`     // 账号
	Password      string `gorm:"type:varchar(255)" json:"password"`                                             // 密码
	Name          string `gorm:"type:varchar(20)" json:"name"`                                                  // 名称
	HeadImage     string `gorm:"type:varchar(255)" json:"headImage"`                                            // 头像地址
	Status        int8   `gorm:"type:tinyint" json:"status"`                                                    // 状态 1.启用 2.停用 3.未激活
	Role          int8   `gorm:"type:tinyint" json:"role"`                                                      // 角色 1.管理员 2.普通用户
	Mail          string `gorm:"type:varchar(255)" json:"mail"`                                                 // 邮箱
	Token         string `gorm:"-" json:"token"`                                                                // 仅用于API返回
	OauthProvider string `gorm:"type:varchar(50);uniqueIndex:idx_username_oauth_provider" json:"oauthProvider"` // OAuth来源 (github, google)
	OauthID       string `gorm:"type:varchar(255);index" json:"oauthId"`                                        // OAuth提供商中的用户ID
}

type UserRepo struct {
}

type IUserRepo interface {
	Get(id uint) (User, error)
	Count() (uint, error)
	GetByUsernameAndPassword(username, password, oauthProvider string) (User, error)
	GetByOAuthID(source, oauthID string) (User, error)
	GetList(pagedParam PagedParam) ([]User, uint, error)
	Update(id uint, user *User) error
	UpdateUserInfo(id uint, updateInfo map[string]any) error
	Create(user *User) error
	Delete(userId uint) ([]string, error)
	CheckUsernameExist(username, oauthProvider string) (User, error)
}

func NewUserRepo() IUserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Get(id uint) (User, error) {
	mUser := User{}
	err := Db.Where("id=?", id).First(&mUser).Error
	return mUser, err
}

func (r *UserRepo) Count() (uint, error) {
	var count int64
	err := Db.Model(&User{}).Count(&count).Error
	return uint(count), err
}

func (r *UserRepo) GetByUsernameAndPassword(username, password, oauthProvider string) (User, error) {
	userInfo := User{}
	err := Db.Where("username=?", username).Where("oauth_provider=?", oauthProvider).
		Where("password=?", password).First(&userInfo).Error
	return userInfo, err
}

func (r *UserRepo) GetByOAuthID(oauthProvider, oauthID string) (User, error) {
	userInfo := User{}
	err := Db.Where("oauth_provider=?", oauthProvider).Where("oauth_id=?", oauthID).First(&userInfo).Error
	return userInfo, err
}

func (r *UserRepo) GetList(pagedParam PagedParam) ([]User, uint, error) {
	var count int64
	if err := Db.Model(&User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var list []User
	if err := Db.Omit("Password").Limit(pagedParam.Limit).Offset(CalcOffset(pagedParam)).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, uint(count), nil
}

func (r *UserRepo) Update(id uint, user *User) error {
	return Db.Where("id=?", id).Updates(user).Error
}

func (r *UserRepo) UpdateUserInfo(userId uint, updateInfo map[string]any) error {
	data := map[string]any{}
	if v, ok := updateInfo["name"]; ok {
		data["name"] = v
	}
	if v, ok := updateInfo["head_image"]; ok {
		data["head_image"] = v
	}
	if v, ok := updateInfo["status"]; ok {
		data["status"] = v
	}
	if v, ok := updateInfo["role"]; ok {
		data["role"] = v
	}

	if v, ok := updateInfo["mail"]; ok {
		hasUser := User{}
		count := Db.Where("mail=?", updateInfo["mail"]).First(&hasUser).RowsAffected
		if count != 0 && hasUser.ID != userId {
			return errors.New("the mail already exists")
		}
		data["mail"] = v
	}
	if v, ok := updateInfo["username"]; ok {
		hasUser := User{}
		count := Db.Where("username=?", updateInfo["username"]).First(&hasUser).RowsAffected
		if count != 0 && hasUser.ID != userId {
			return errors.New("the username already exists")
		}
		data["username"] = v
	}
	if v, ok := updateInfo["password"]; ok {
		data["password"] = v
	}

	mUser := User{}
	err := Db.Model(&mUser).Where("id=?", userId).Updates(data).Error
	return err
}

func (r *UserRepo) Create(user *User) error {
	err := Db.Create(user).Error
	return err
}

func (r *UserRepo) Delete(userId uint) ([]string, error) {
	var fileNames []string
	err := Db.Transaction(func(tx *gorm.DB) error {
		// Get all files of the user before deletion
		var files []File
		if err := tx.Where("user_id = ?", userId).Find(&files).Error; err != nil {
			return err
		}

		// Store file names for later deletion from storage
		fileNames = make([]string, 0, len(files))
		for _, file := range files {
			fileNames = append(fileNames, file.FileName)
		}

		// 删除图标
		if err := tx.Delete(&ItemIcon{}, "user_id=?", userId).Error; err != nil {
			return err
		}

		// 删除分组
		if err := tx.Delete(&ItemIconGroup{}, "user_id = ?", userId).Error; err != nil {
			return err
		}

		// 删除模块配置
		if err := tx.Delete(&ModuleConfig{}, "user_id=?", userId).Error; err != nil {
			return err
		}

		// 删除文件记录，并没有删除资源文件
		if err := tx.Delete(&File{}, "user_id=?", userId).Error; err != nil {
			return err
		}

		if err := tx.Delete(&User{}, userId).Error; err != nil {
			return err
		}

		return nil
	})
	
	return fileNames, err
}

func (r *UserRepo) CheckUsernameExist(username, oauthProvider string) (User, error) {
	hasUser := User{}
	count := Db.Where("username=?", username).Where("oauth_provider=?", oauthProvider).First(&hasUser).RowsAffected
	if count != 0 {
		return hasUser, errors.New("该用户名已被注册")
	}

	return hasUser, nil
}
