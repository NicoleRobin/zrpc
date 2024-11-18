package log

import (
	"context"
	"go.uber.org/zap"
)

var defaultLogger *zap.SugaredLogger

func getDefaultLogger() *zap.SugaredLogger {
	return defaultLogger
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
}

func DebugF(ctx context.Context, format string, args ...interface{}) {
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
}

func InfoF(ctx context.Context, format string, args ...interface{}) {
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
}

func WarnF(ctx context.Context, format string, args ...interface{}) {
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
}

func ErrorF(ctx context.Context, format string, args ...interface{}) {
}
