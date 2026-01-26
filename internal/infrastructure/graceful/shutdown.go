package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/xonvera-core/internal/infrastructure/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Shutdown handles graceful shutdown of the application
func Shutdown(app *fiber.App, db *gorm.DB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown Fiber server
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	// Close database connections
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				logger.Error("Failed to close database", zap.Error(err))
			} else {
				logger.Info("Database connections closed")
			}
		}
	}

	logger.Info("Server exited")
}
