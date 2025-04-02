package zaplog

import (
	"io"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogStruct struct {
	Writer    io.Writer
	File      *os.File
	PrintCfg  bool // 此条打印到控制台
	Separator string
}

var Logger *zap.SugaredLogger

type LogFileld map[string]string

func InitLog(runmode string, filePath string) error {
	runtimePath := "./logs"
	if err := os.MkdirAll(runtimePath, 0777); err != nil {
		return err
	}

	var level zap.AtomicLevel
	if runmode == "debug" {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		level = zap.NewAtomicLevel() // 支持通过http以及配置文件动态修改日志级别
	}

	fileWriteSyncer := getLogWriter(runtimePath + "/" + filePath)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(fileWriteSyncer, zapcore.AddSync(os.Stdout)), level)
	Logger = zap.New(core, zap.AddCaller()).Sugar()
	return nil
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
