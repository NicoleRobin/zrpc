package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int

const (
	LevelNil Level = iota
	LevelDebug
	LevelWarn
	LevelInfo
	LevelError
)

var levelToZapLevel = map[Level]zapcore.Level{
	LevelDebug: zapcore.DebugLevel,
	LevelInfo:  zapcore.InfoLevel,
	LevelWarn:  zapcore.WarnLevel,
	LevelError: zapcore.ErrorLevel,
}

var defaultLogger *zap.Logger

func newConfig() zap.Config {
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.TimeKey = "time"
	encoderConf.MessageKey = "message"
	encoderConf.FunctionKey = "function"
	encoderConf.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConf.EncodeCaller = zapcore.ShortCallerEncoder

	logConf := zap.NewProductionConfig()
	logConf.EncoderConfig = encoderConf
	logConf.OutputPaths = []string{"stdout"}
	logConf.ErrorOutputPaths = []string{"stderr"}
	return logConf
}

func newLogger(config zap.Config) *zap.Logger {
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	return logger
}

func init() {
	logConfig := newConfig()
	logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	defaultLogger = newLogger(logConfig)
}

func SetLevel(level Level) {
	if defaultLogger == nil {
		panic("defaultLogger is nil")
	}
	if zapLevel, ok := levelToZapLevel[level]; ok {
		logConfig := newConfig()
		logConfig.Level = zap.NewAtomicLevelAt(zapLevel)
		defaultLogger = newLogger(logConfig)
	} else {
		panic("unknown log level")
	}

}

func getDefaultLogger() *zap.Logger {
	return defaultLogger
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Debug(fmt.Sprintf(format, args...))
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Info(fmt.Sprintf(format, args...))
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Warn(fmt.Sprintf(format, args...))
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, fields...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Debug(fmt.Sprintf(format, args...))
}
