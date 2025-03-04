package repository

import (
	"errors"
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

func (m *User) CreateOne() (User, error) {
	err := Db.Create(m).Error
	return *m, err
}

func (m *User) CheckUsernameExist(username string) (User, error) {
	hasUser := User{}
	count := Db.Where("username=?", username).First(&hasUser).RowsAffected
	if count != 0 {
		return hasUser, errors.New("该用户名已被注册")
	}
	return hasUser, nil
}
