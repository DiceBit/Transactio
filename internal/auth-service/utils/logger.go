package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger(filePath string) *zap.Logger {

	var file *os.File
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		file, _ = os.Create(filePath)
	}

	stdout := zapcore.AddSync(os.Stdout)
	logFile := zapcore.AddSync(file)

	logFile = zapcore.Lock(logFile)

	var encoderCfg zapcore.EncoderConfig
	if AppEnv == "dev" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	lvl := zapcore.InfoLevel

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, stdout, lvl),
		zapcore.NewCore(encoder, logFile, lvl),
	)

	return zap.New(core)
}
