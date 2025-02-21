package config

import (
	"errors"
	"sun-panel/lib/cmn"
	"sun-panel/lib/iniConfig"
)

func ConfigInit() (*iniConfig.IniConfig, error) {
	exists, err := cmn.PathExists("conf/conf.ini")
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("conf.ini 配置文件不存在，请参考 conf.example.ini 创建")
	}

	config := iniConfig.NewIniConfig("conf/conf.ini")
	return config, nil
}
