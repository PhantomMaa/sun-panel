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
		createConfExample("conf.example.ini", "conf.example.ini")
		return nil, errors.New("conf.ini 配置文件不存在，请参考 conf.example.ini 创建")
	}

	config := iniConfig.NewIniConfig("conf/conf.ini") // 读取配置
	return config, nil
}

// 生成示例配置文件
func createConfExample(confName string, targetName string) {
	exists, _ := cmn.PathExists("conf/" + targetName)
	if !exists {
		cmn.AssetsTakeFileToPath(confName, "conf/"+targetName)
	}
}
