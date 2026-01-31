package logger

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

const requestIDKey = "request_id"

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

// FromContext extracts request ID from Fiber context and returns logger with request ID field
func FromContext(c *fiber.Ctx) *zap.Logger {
	if c == nil {
		return Get()
	}

	requestID := extractRequestID(c)
	if requestID == "" {
		return Get()
	}

	return Get().With(zap.String(requestIDKey, requestID))
}

// FromStdContext extracts request ID from standard context and returns logger with request ID field
func FromStdContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return Get()
	}

	// Try to extract request ID from context value
	if requestID, ok := ctx.Value(requestIDKey).(string); ok && requestID != "" {
		return Get().With(zap.String(requestIDKey, requestID))
	}

	return Get()
}

// extractRequestID safely extracts request ID from Fiber context
func extractRequestID(c *fiber.Ctx) string {
	if id, ok := c.Locals(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// ContextDebug logs a debug message with request ID from context
func ContextDebug(c *fiber.Ctx, msg string, fields ...zap.Field) {
	FromContext(c).Debug(msg, fields...)
}

// ContextInfo logs an info message with request ID from context
func ContextInfo(c *fiber.Ctx, msg string, fields ...zap.Field) {
	FromContext(c).Info(msg, fields...)
}

// ContextWarn logs a warning message with request ID from context
func ContextWarn(c *fiber.Ctx, msg string, fields ...zap.Field) {
	FromContext(c).Warn(msg, fields...)
}

// ContextError logs an error message with request ID from context
func ContextError(c *fiber.Ctx, msg string, fields ...zap.Field) {
	FromContext(c).Error(msg, fields...)
}

// StdContextDebug logs a debug message with request ID from standard context
func StdContextDebug(ctx context.Context, msg string, fields ...zap.Field) {
	FromStdContext(ctx).Debug(msg, fields...)
}

// StdContextInfo logs an info message with request ID from standard context
func StdContextInfo(ctx context.Context, msg string, fields ...zap.Field) {
	FromStdContext(ctx).Info(msg, fields...)
}

// StdContextWarn logs a warning message with request ID from standard context
func StdContextWarn(ctx context.Context, msg string, fields ...zap.Field) {
	FromStdContext(ctx).Warn(msg, fields...)
}

// StdContextError logs an error message with request ID from standard context
func StdContextError(ctx context.Context, msg string, fields ...zap.Field) {
	FromStdContext(ctx).Error(msg, fields...)
}
