package cache

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"sun-panel/internal/repository"
)

const (
	PanelPublicUserId = "panel_public_user_id" // 公开访问模式用户id *uint|null
)

type SystemSetting struct {
	Cache Cacher[interface{}]
}

type Register struct {
	EmailSuffix  string `json:"emailSuffix"`  // 注册邮箱后缀
	OpenRegister bool   `json:"openRegister"` // 开放注册
}

type Login struct {
	LoginCaptcha bool `json:"loginCaptcha"` // 登录验证码
}

type ApplicationSetting struct {
	Register
	Login
	WebSiteUrl string `json:"webSiteUrl"` // 站点地址
}

var (
	ErrorNoExists = errors.New("no exists")
)

func (s *SystemSetting) GetValueString(configName string) (result string, err error) {
	if v, ok := s.Cache.Get(configName); ok {
		if v1, ok1 := v.(string); ok1 {
			// fmt.Println("读取缓存")
			return v1, nil
		}
	}

	mSetting := repository.SystemSetting{}
	result, err = mSetting.Get(configName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrorNoExists
	}
	// 查询出来，缓存起来
	s.Cache.SetDefault(configName, result)
	return
}

func (s *SystemSetting) GetValueByInterface(configName string, value interface{}) error {
	if v, ok := s.Cache.Get(configName); ok {
		if s, sok := v.(string); sok {
			if err := json.Unmarshal([]byte(s), value); err != nil {
				return err
			}
			return nil
		}
	}

	mSetting := repository.SystemSetting{}
	result, err := mSetting.Get(configName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrorNoExists
		return err
	}

	err = json.Unmarshal([]byte(result), value)
	if err != nil {
		return err
	}
	s.Cache.SetDefault(configName, result)
	return nil
}

func (s *SystemSetting) Set(configName string, configValue interface{}) error {
	s.Cache.Delete(configName)
	mSetting := repository.SystemSetting{}
	err := mSetting.Set(configName, configValue)
	return err
}
