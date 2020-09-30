package lib

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	once sync.Once
	Logger *zap.Logger
)