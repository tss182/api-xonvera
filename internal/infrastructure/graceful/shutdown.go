package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"app/xonvera-core/internal/infrastructure/logger"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	shutdownTimeout = 30 * time.Second
)

// Shutdown handles graceful shutdown of the application
func Shutdown(app *fiber.App, db *gorm.DB, redisClient *redis.Client, serverErrors <-chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	// Wait for either a signal or server error
	select {
	case sig := <-quit:
		logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	case err := <-serverErrors:
		logger.Error("Server error, initiating shutdown", zap.Error(err))
	}

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Use WaitGroup to coordinate concurrent shutdown operations
	var wg sync.WaitGroup
	shutdownComplete := make(chan struct{})

	// Shutdown Fiber server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Shutting down HTTP server...")
		if err := app.ShutdownWithContext(ctx); err != nil {
			logger.Error("HTTP server forced shutdown", zap.Error(err))
		} else {
			logger.Info("HTTP server shutdown complete")
		}
	}()

	// Close database connections
	if db != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			logger.Info("Closing database connections...")
			if err := closeDatabase(ctx, db); err != nil {
				logger.Error("Failed to close database", zap.Error(err))
			} else {
				logger.Info("Database connections closed")
			}
		}()
	}

	// Close Redis connections
	if redisClient != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			logger.Info("Closing Redis connection...")
			if err := redisClient.Close(); err != nil {
				logger.Error("Failed to close Redis connection", zap.Error(err))
			} else {
				logger.Info("Redis connection closed")
			}
		}()
	}

	// Wait for all shutdown operations to complete
	go func() {
		wg.Wait()
		close(shutdownComplete)
	}()

	// Wait for completion or timeout
	select {
	case <-shutdownComplete:
		logger.Info("Graceful shutdown completed successfully")
	case <-ctx.Done():
		logger.Warn("Shutdown timeout exceeded, forcing exit", zap.Duration("timeout", shutdownTimeout))
	}
}

// closeDatabase safely closes the database connection
func closeDatabase(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.WithContext(ctx).DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
