package config

import (
	"errors"
	"sun-panel/lib/cmn"
	"sun-panel/lib/embedfs"
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
		err := embedfs.ExtractEmbeddedFile(confName, "conf/"+targetName)
		if err != nil {
			// 记录错误但继续执行，因为这只是示例配置文件
			cmn.AssetsTakeFileToPath(confName, "conf/"+targetName)
		}
	}
}
