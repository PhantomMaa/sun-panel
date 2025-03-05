package zaplog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sun-panel/internal/common"
	"time"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogStruct struct {
	Writer    io.Writer
	File      *os.File
	PrintCfg  bool // 此条打印到控制台
	Separator string
}

type LogFileld map[string]string

// 日志颜色
var colors = map[string]func(a ...interface{}) string{
	"Warning": color.New(color.FgYellow).Add(color.Bold).SprintFunc(),
	"Panic":   color.New(color.BgRed).Add(color.Bold).SprintFunc(),
	"Error":   color.New(color.FgRed).Add(color.Bold).SprintFunc(),
	"Info":    color.New(color.FgCyan).Add(color.Bold).SprintFunc(),
	"Debug":   color.New(color.FgWhite).Add(color.Bold).SprintFunc(),
}

// 不同级别前缀与时间的间隔，保持宽度一致
var spaces = map[string]string{
	"Warning": "",
	"Panic":   "  ",
	"Error":   "  ",
	"Info":    "   ",
	"Debug":   "  ",
}

// 运行日志静态类
var runLogStatic = LogStruct{}

func InitLog(runmode string, filePath string) (*zap.SugaredLogger, error) {
	runtimePath := "./logs"
	if err := os.MkdirAll(runtimePath, 0777); err != nil {
		return nil, err
	}

	var level zap.AtomicLevel
	if runmode == "debug" {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		level = zap.NewAtomicLevel() // 支持通过http以及配置文件动态修改日志级别
	}

	logger := InitLogger(runtimePath+"/"+filePath, level)
	return logger, nil
}

func InitLogger(fileName string, level zapcore.LevelEnabler) *zap.SugaredLogger {
	fileWriteSyncer := getLogWriter(fileName)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(fileWriteSyncer, zapcore.AddSync(os.Stdout)), level)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}

func Pln(prefix string, msg string) {
	fmt.Printf(
		"%s%s %s %s\n",
		colors[prefix]("["+prefix+"]"),
		spaces[prefix],
		time.Now().Format(common.TimeFormatMode1),
		msg,
	)
}

func (t *LogStruct) Write(content string) (n int, err error) {

	return io.WriteString(t.Writer, content)
}

func (t *LogStruct) Format(log_type string, content string) (n int, err error) {
	content = log_type + spaces[log_type] + " " + common.GetTime() + " " + content + "\n"
	return t.Write(content)
}

func (t *LogStruct) Info(content ...string) (n int, err error) {
	str := ""
	for i := 0; i < len(content); i++ {
		if i != 0 {
			str += t.Separator + content[i]
		} else {
			str += content[i]
		}
	}
	n, err = t.Format("Info", str)
	if t.PrintCfg == true {
		Pln("Info", str)
		t.PrintCfg = false
	}
	return
}

func (t *LogStruct) Debug(content string) {
	t.Format("Debug", content)
	if t.PrintCfg == true {
		Pln("Debug", content)
		t.PrintCfg = false
	}
}

func (t *LogStruct) Error(content ...string) {
	contentStr := ""
	for i := 0; i < len(content); i++ {
		if i != 0 {
			contentStr += t.Separator + content[i]
		} else {
			contentStr += content[i]
		}
	}
	t.Format("Error", contentStr)
	if t.PrintCfg == true {
		Pln("Error", contentStr)
		t.PrintCfg = false
	}
}

func getEncoder() zapcore.Encoder {
	logConf := zap.NewProductionEncoderConfig()
	logConf.EncodeTime = zapcore.ISO8601TimeEncoder
	logConf.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(logConf)
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("failed to create zaplog file", fileName)
	}
	return zapcore.AddSync(file)
}
