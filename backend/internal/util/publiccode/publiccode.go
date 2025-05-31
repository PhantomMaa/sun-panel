package publiccode

import (
	"errors"
	"sun-panel/internal/global"
	"sun-panel/internal/infra/zaplog"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	LENGTH  = 10
	CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	// ErrInvalidCode 无效的公开访问代码
	ErrInvalidCode = errors.New("invalid public access code")
)

// GenerateAndSave 生成 public visit 代码
func GenerateAndSave(userID uint) (string, error) {
	code, err := gonanoid.Generate(CHARSET, LENGTH)
	if err != nil {
		zaplog.Logger.Errorf("generate public visit code failed: %v", err)
		return "", err
	}

	zaplog.Logger.Infof("generate public visit code success: %s", code)
	// 更新用户的 publiccode 字段
	updateInfo := map[string]any{
		"publiccode": code,
	}

	// 更新用户信息
	if err := global.UserRepo.UpdateUserInfo(userID, updateInfo); err != nil {
		zaplog.Logger.Errorf("update user publiccode failed: %v", err)
		return "", err
	}

	return code, nil
}

// ParseCode 解析 public visit 代码，获取用户ID
func ParseCode(code string) (uint, error) {
	if code == "" {
		return 0, ErrInvalidCode
	}

	// 从数据库中查询具有指定 publiccode 的用户
	userInfo, err := global.UserRepo.GetByPubliccode(code)
	if err != nil {
		zaplog.Logger.Errorf("get user by publiccode failed: %v", err)
		return 0, err
	}

	return userInfo.ID, nil
}

// RemoveCode 移除 public visit 代码
func RemoveCode(userID uint) error {
	// 更新用户的 publiccode 字段为空
	updateInfo := map[string]any{
		"publiccode": "",
	}

	// 更新用户信息
	if err := global.UserRepo.UpdateUserInfo(userID, updateInfo); err != nil {
		zaplog.Logger.Errorf("remove user publiccode failed: %v", err)
		return err
	}

	return nil
}
