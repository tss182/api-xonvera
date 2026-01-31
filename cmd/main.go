// Package main Xonvera Service API
// @title Xonvera API
// @description Xonvera Service API - Cashflow and Invoice Management
// @version 1.0.0
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
package main

import (
	"flag"

	"app/xonvera-core/internal/adapters/routes"
	"app/xonvera-core/internal/dependencies"
	"app/xonvera-core/internal/infrastructure/database"
	"app/xonvera-core/internal/infrastructure/graceful"
	"app/xonvera-core/internal/infrastructure/logger"
	"app/xonvera-core/internal/utils/validator"

	"go.uber.org/zap"
)

func main() {
	// Parse command line flags
	runMigrations := flag.Bool("migrate", false, "Run database migrations")
	migrateDown := flag.Bool("migrate-down", false, "Rollback database migrations by one step")
	migrateReset := flag.Bool("migrate-reset", false, "Reset database migrations")
	flag.Parse()

	// Initialize application using Wire
	app, err := dependencies.InitializeApplication()
	if err != nil {
		logger.Fatal("Failed to initialize application", zap.Error(err))
	}

	// Initialize logger with environment
	logger.Init(app.Config.App.Env)
	defer logger.Sync()

	// Initialize validator with pagination config
	validator.SetPaginationDefaults(app.Config.Pagination.DefaultLimit, app.Config.Pagination.MaxLimit)

	// Handle migrations if requested
	if *runMigrations || *migrateDown || *migrateReset {
		dsn := database.GetDSN(&app.Config.Database)
		migrator, err := database.NewMigrator(dsn, "file://migrations")
		if err != nil {
			logger.Fatal("Failed to create migrator", zap.Error(err))
		}
		defer migrator.Close()

		if *migrateDown {
			if err := migrator.Steps(-1); err != nil {
				logger.Fatal("Failed to rollback migrations", zap.Error(err))
			}
		} else if *migrateReset {
			if err := migrator.Down(); err != nil {
				logger.Fatal("Failed to reset migrations", zap.Error(err))
			}
		} else {
			if err := migrator.Up(); err != nil {
				logger.Fatal("Failed to run migrations", zap.Error(err))
			}
		}
		return
	}

	// Get Fiber app (middleware already configured in server package)
	fiberApp := app.FiberApp

	// Setup routes
	routes.SetupRoutes(fiberApp, app)

	addr := ":" + app.Config.App.Port
	logger.Info("Starting server",
		zap.String("app", app.Config.App.Name),
		zap.String("addr", addr),
		zap.String("env", app.Config.App.Env),
	)

	// Channel to handle server errors
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		if err := fiberApp.Listen(addr); err != nil {
			serverErrors <- err
		}
	}()

	// Handle graceful shutdown or server errors
	graceful.Shutdown(fiberApp, app.DB, app.Redis, serverErrors)
}
