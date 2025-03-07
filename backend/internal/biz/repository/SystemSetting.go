package repository

import (
	"encoding/json"
	"errors"
	"sun-panel/internal/util"

	"gorm.io/gorm"
)

type SystemSetting struct {
	ID          uint   `gorm:"primaryKey"`
	ConfigName  string `gorm:"type:varchar(50)"`
	ConfigValue string `gorm:"type:text"`
}

type SystemSettingRepo struct{}

type ISystemSettingRepo interface {
}

func NewSystemSettingRepo() *SystemSettingRepo {
	return &SystemSettingRepo{}
}

func (r *SystemSettingRepo) Get(configName string) (result string, err error) {
	var res SystemSetting
	if err := Db.Select("ConfigValue").First(&res, "config_name=?", configName).Error; err != nil {
		return result, err
	}
	result = res.ConfigValue
	return result, nil
}

func (r *SystemSettingRepo) GetValueByInterface(configName string, structValue interface{}) error {
	result, err := r.Get(configName)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(result), structValue)
	if err != nil {
		return err
	}
	return nil
}

func (r *SystemSettingRepo) Set(configName string, configValue interface{}) error {
	findRes := SystemSetting{}
	db := Db.First(&findRes, "config_name=?", configName)
	if err := db.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	value := ""
	if s, ok := configValue.(string); !ok {
		value = util.ToJSONString(configValue)
	} else {
		value = s
	}

	if db.RowsAffected == 0 {
		// 添加
		if err := Db.Create(&SystemSetting{ConfigName: configName, ConfigValue: value}).Error; err != nil {
			return err
		}

	} else {
		// 修改
		if err := Db.Where("id=?", findRes.ID).Update("config_value", value).Error; err != nil {
			return err
		}
	}
	return nil
}
