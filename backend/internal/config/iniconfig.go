package config

import (
	"gopkg.in/ini.v1"
)

type IniConfig struct {
	Err      error
	Config   *ini.File
	FileName string
}

// 获取配置
func (t *IniConfig) GetValueString(section string, name string) string {
	return t.Config.Section(section).Key(name).String()
}

// 获取配置
func (t *IniConfig) GetValueInt(section string, name string) int {
	return t.Config.Section(section).Key(name).MustInt()
}

// 创建一个配置对象
func NewIniConfig(filename string) *IniConfig {
	config, err := ini.Load(filename)

	return &IniConfig{
		Err:      err,
		Config:   config,
		FileName: filename,
	}
}
