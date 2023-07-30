package logger

import "go.uber.org/zap"

func InitLogger() *zap.Logger {
	return zap.Must(zap.NewDevelopment())
}
