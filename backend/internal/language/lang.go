package language

import (
	"os"
	"strings"
	"sun-panel/internal/common"
	"sun-panel/internal/config"
	"sun-panel/internal/global"
)

var Obj *LangStructObj

type LangStructObj struct {
	LangContet *config.IniConfig
}

func LangInit(lang string) {
	filename := "lang/" + lang + ".ini"
	exists, err := common.PathExists(filename)
	if err != nil || !exists {
		global.Logger.Errorln("语言文件不存在:", filename)
		os.Exit(1)
	}

	Obj = NewLang(filename)
}

func NewLang(langPath string) *LangStructObj {
	langObj := LangStructObj{}
	langObj.LangContet = config.NewIniConfig(langPath) // 读取配置
	return &langObj
}

func (l *LangStructObj) Get(key string) string {
	if key == "" {
		return key
	}
	keyArr := strings.Split(key, ".")
	if len(keyArr) < 2 {
		return l.LangContet.GetValueString(keyArr[0], "NOT EMPTY")
	} else {
		return l.LangContet.GetValueString(keyArr[0], keyArr[1])
	}
}

// 获取并替换字段
func (l *LangStructObj) GetWithFields(key string, fields map[string]string) string {
	c := l.Get(key)
	for k, v := range fields {
		c = strings.ReplaceAll(c, `{`+k+`}`, v)
	}
	return c
}

// 获取值并向后追加
func (l *LangStructObj) GetAndInsert(key string, insertContent ...string) string {
	content := l.Get(key) + " "
	for _, v := range insertContent {
		content += v
	}
	return content
}
