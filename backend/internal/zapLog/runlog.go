package zapLog

import (
	"os"
	"sun-panel/internal/common"
	"sun-panel/internal/global"

	"go.uber.org/zap"
)

func InitLog(runmode string, filePath string) (*zap.SugaredLogger, error) {
	runtimePath := "./logs"
	if err := os.MkdirAll(runtimePath, 0777); err != nil {
		return nil, err
	}

	var level zap.AtomicLevel
	if runmode == "debug" {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		level = global.LoggerLevel
	}

	logger := common.InitLogger(runtimePath+"/"+filePath, level)
	return logger, nil
}
