package i18n

import (
	"os"
	"strings"
	"sun-panel/internal/global"
	"sun-panel/internal/util"
	"gopkg.in/ini.v1"
)

var Obj *LangStructObj

type LangStructObj struct {
	LangContent *ini.File
}

func LangInit(lang string) {
	filename := "lang/" + lang + ".ini"
	exists, err := util.PathExists(filename)
	if err != nil || !exists {
		global.Logger.Errorln("语言文件不存在:", filename)
		os.Exit(1)
	}

	Obj = NewLang(filename)
}

func NewLang(langPath string) *LangStructObj {
	langObj := LangStructObj{}
	iniFile, err := ini.Load(langPath)
	if err != nil {
		global.Logger.Errorln("加载语言文件失败:", err)
		os.Exit(1)
	}
	langObj.LangContent = iniFile
	return &langObj
}

func (l *LangStructObj) Get(key string) string {
	if key == "" {
		return key
	}
	keyArr := strings.Split(key, ".")
	if len(keyArr) < 2 {
		return l.LangContent.Section(keyArr[0]).Key("NOT EMPTY").String()
	} else {
		return l.LangContent.Section(keyArr[0]).Key(keyArr[1]).String()
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
