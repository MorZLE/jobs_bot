package logger

import (
	"go.uber.org/zap"
	"log"
	"os"
)

var mylog *zap.Logger = zap.NewNop()

// Initialize инициализирует собственный zap logger
func Initialize() {
	lvl := zap.LevelFlag("info", zap.InfoLevel, "log level")
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(*lvl)
	sd, er := cfg.Build()
	if er != nil {
		log.Fatal(er)
	}
	mylog = sd
}

func Info(info string) {
	mylog.Info("INFO", zap.String("info", info))
}

func Error(info string, err error) {
	mylog.Error(info, zap.Error(err))
}

func Fatal(info string, err error) {
	mylog.Error(info, zap.Error(err))
	os.Exit(1)
}
