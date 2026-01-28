package database

import (
	"context"
	"fmt"
	"time"

	"app/xonvera-core/internal/infrastructure/config"
	"app/xonvera-core/internal/infrastructure/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// zapGormLogger adapts zap logger to gorm logger interface
type zapGormLogger struct {
	zap *zap.Logger
}

func newZapGormLogger() gormlogger.Interface {
	return &zapGormLogger{zap: logger.Get()}
}

func (l *zapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *zapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.zap.Sugar().Infof(msg, data...)
}

func (l *zapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.zap.Sugar().Warnf(msg, data...)
}

func (l *zapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.zap.Sugar().Errorf(msg, data...)
}

func (l *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
	}

	// Log errors with error level
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.zap.Error("database query error", fields...)
		return
	}

	// Log slow queries with warning level
	const slowQueryThreshold = 200 * time.Millisecond
	if elapsed > slowQueryThreshold {
		l.zap.Warn("slow database query", fields...)
	} else {
		l.zap.Debug("database query", fields...)
	}
}

func NewConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newZapGormLogger(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifeTime)

	logger.Info("Database connection established",
		zap.String("host", cfg.Host),
		zap.String("database", cfg.Name),
		zap.Int("max_idle_conn", cfg.MaxIdleConn),
		zap.Int("max_open_conn", cfg.MaxOpenConn),
		zap.Duration("conn_max_idle_time", cfg.ConnMaxIdleTime),
		zap.Duration("conn_max_life_time", cfg.ConnMaxLifeTime),
	)

	return db, nil
}

// GetDSN returns the database connection string for migrations
func GetDSN(cfg *config.DatabaseConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)
}
