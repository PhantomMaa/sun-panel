package config

import (
	"errors"
	"sun-panel/internal/util"
)

func Init(configPath string) (*IniConfig, error) {
	exists, err := util.PathExists(configPath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("配置文件不存在: " + configPath)
	}

	config := NewIniConfig(configPath)
	return config, nil
}
