package config

import (
	"errors"
	"sun-panel/lib/cmn"
	"sun-panel/lib/iniConfig"
)

func ConfigInit(configPath string) (*iniConfig.IniConfig, error) {
	exists, err := cmn.PathExists(configPath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("配置文件不存在: " + configPath)
	}

	config := iniConfig.NewIniConfig(configPath)
	return config, nil
}
