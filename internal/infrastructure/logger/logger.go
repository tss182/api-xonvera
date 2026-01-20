package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the zap logger
func Init(env string) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		// config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	var err error
	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

// Get returns the global logger instance
func Get() *zap.Logger {
	if log == nil {
		Init("development")
	}
	return log
}

// Sync flushes any buffered log entries
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
	os.Exit(1)
}

// With creates a child logger with additional fields
func With(fields ...zap.Field) *zap.Logger {
	return Get().With(fields...)
}
