package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	once sync.Once
	Logger *zap.Logger
)

func init() {
	initLogger()
}

func initLogger() *zap.Logger {
	once.Do(func() {
		core := zapcore.NewCore(
			getEncoder(),
			getWriteSyncer(),
			zapcore.InfoLevel,
		)
		caller := zap.AddCaller()
		dev := zap.Development()
		Logger = zap.New(core, caller, dev)
	})
	Logger.Info("zap loh init success")
	return Logger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getWriteSyncer() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}