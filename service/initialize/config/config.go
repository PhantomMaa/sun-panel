package config

import (
	"errors"
	"fmt"
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
		if err := createConfExample("conf.example.ini", "conf.example.ini"); err != nil {
			return nil, fmt.Errorf("创建配置示例文件失败: %v", err)
		}
		return nil, errors.New("conf.ini 配置文件不存在，请参考 conf.example.ini 创建")
	}

	config := iniConfig.NewIniConfig("conf/conf.ini") // 读取配置
	return config, nil
}

// 生成示例配置文件
func createConfExample(confName string, targetName string) error {
	targetPath := "conf/" + targetName
	exists, err := cmn.PathExists(targetPath)
	if err != nil {
		return fmt.Errorf("检查配置文件路径失败: %v", err)
	}

	if !exists {
		// ExtractEmbeddedFile 会自动在内部添加 "assets/" 前缀并创建必要的目录
		if err := embedfs.ExtractEmbeddedFile(confName, targetPath); err != nil {
			return fmt.Errorf("提取配置示例文件失败: %v", err)
		}
	}
	return nil
}
