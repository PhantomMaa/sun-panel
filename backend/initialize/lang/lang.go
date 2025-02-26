package lang

import (
	"os"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/language"
)

func LangInit(lang string) {
	filename := "lang/" + lang + ".ini"
	exists, err := cmn.PathExists(filename)
	if err != nil || !exists {
		global.Logger.Errorln("语言文件不存在:", filename)
		os.Exit(1)
	}

	global.Lang = language.NewLang(filename)
}
