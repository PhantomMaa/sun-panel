package cache

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/infra/kvcache"
)

type SystemSetting struct {
	Cache             kvcache.Cacher[interface{}]
	SystemSettingRepo *repository.SystemSettingRepo
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

	result, err = s.SystemSettingRepo.Get(configName)
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

	result, err := s.SystemSettingRepo.Get(configName)
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
	err := s.SystemSettingRepo.Set(configName, configValue)
	return err
}
