package repository

import (
	"errors"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Username  string `gorm:"type:varchar(255);uniqueIndex" json:"username"` // 账号
	Password  string `gorm:"type:varchar(255)" json:"password"`             // 密码
	Name      string `gorm:"type:varchar(20)" json:"name"`                  // 名称
	HeadImage string `gorm:"type:varchar(255)" json:"headImage"`            // 头像地址
	Status    int8   `gorm:"type:tinyint" json:"status"`                    // 状态 1.启用 2.停用 3.未激活
	Role      int8   `gorm:"type:tinyint" json:"role"`                      // 角色 1.管理员 2.普通用户
	Mail      string `gorm:"type:varchar(255)" json:"mail"`                 // 邮箱
	Token     string `gorm:"-" json:"token"`                                // 仅用于API返回
	UserId    uint   `gorm:"-"  json:"userId"`
}

type UserRepo struct {
}

type IUserRepo interface {
	Create(user *User) error
	Deletes(userIds []uint) error
}

func NewUserRepo() IUserRepo {
	return &UserRepo{}
}

func (m *User) GetUserInfoByUid(uid uint) (User, error) {
	mUser := User{}
	err := Db.Where("id=?", uid).First(&mUser).Error
	return mUser, err
}

func (m *User) GetUserInfoByUsernameAndPassword(username, password string) (User, error) {
	userInfo := User{}
	err := Db.Where("username=?", username).Where("password=?", password).First(&userInfo).Error
	return userInfo, err
}

func (m *User) GetUserInfoByUsername(username string) (User, error) {
	mUser := User{}
	err := Db.Where("username=?", username).First(&mUser).Error
	return mUser, err
}

func (m *User) GetUserInfoByMail() *User {
	mUser := User{}
	if Db.Where("mail=?", m.Mail).First(&mUser).Error != nil {
		return nil
	}
	return &mUser
}

func (m *User) UpdateUserInfoByUserId(userId uint, updateInfo map[string]interface{}) error {
	mUser := User{}

	data := map[string]interface{}{}
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
	if v, ok := updateInfo["gender"]; ok {
		data["gender"] = v
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

	err := Db.Model(&mUser).Where("id=?", userId).Updates(data).Error
	return err
}

func (r *UserRepo) Create(user *User) error {
	err := Db.Create(user).Error
	return err
}

func (m *User) CheckUsernameExist(username string) (User, error) {
	hasUser := User{}
	count := Db.Where("username=?", username).First(&hasUser).RowsAffected
	if count != 0 {
		return hasUser, errors.New("该用户名已被注册")
	}
	return hasUser, nil
}

func (r *UserRepo) Deletes(userIds []uint) error {
	return Db.Transaction(func(tx *gorm.DB) error {
		for _, v := range userIds {
			// 删除图标
			if err := tx.Delete(&ItemIcon{}, "user_id=?", v).Error; err != nil {
				return err
			}

			// 删除分组
			if err := tx.Delete(&ItemIconGroup{}, "user_id = ?", v).Error; err != nil {
				return err
			}

			// 删除模块配置
			if err := tx.Delete(&ModuleConfig{}, "user_id=?", v).Error; err != nil {
				return err
			}

			// 删除文件记录，并没有删除资源文件
			if err := tx.Delete(&File{}, "user_id=?", v).Error; err != nil {
				return err
			}
		}

		if err := tx.Delete(&User{}, userIds).Error; err != nil {
			return err
		}

		return nil
	})
}
