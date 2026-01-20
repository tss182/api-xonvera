package main

import (
	"flag"

	"app/xonvera-core/internal/adapters/routes"
	"app/xonvera-core/internal/dependencies"
	"app/xonvera-core/internal/infrastructure/database"
	"app/xonvera-core/internal/infrastructure/graceful"
	"app/xonvera-core/internal/infrastructure/logger"
	"app/xonvera-core/internal/infrastructure/redis"

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

	// Handle migrations if requested
	if *runMigrations || *migrateDown || *migrateReset {
		dsn := database.GetDSN(&app.Config.Database)
		migrator, err := database.NewMigrator(dsn, "file://internal/infrastructure/database/migrations")
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

	// Get Fiber app and Redis client from Wire
	defer redis.CloseRedis(app.Redis)

	// Global middleware
	fiberApp.Use(middleware.RequestID())
	fiberApp.Use(recover.New())
	fiberApp.Use(middleware.BodyLogger(app.Config.App.Env))
	fiberApp.Use(fiberlogger.New(fiberlogger.Config{
		Format: "${time} | ${status} | ${latency} | ${locals:request_id} | ${method} ${path}\n",
	}))
	fiberApp.Use(middleware.APIRateLimiter(app.Redis))
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check endpoint
	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"app":    app.Config.App.Name,
		})
	})

	// Setup routes
	routes.SetupRoutes(app.FiberApp, app)

	// Start server in goroutine
	addr := ":" + app.Config.App.Port
	go func() {
		logger.Info("Starting server",
			zap.String("app", app.Config.App.Name),
			zap.String("addr", addr),
			zap.String("env", app.Config.App.Env),
		)
		if err := app.FiberApp.Listen(addr); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Handle graceful shutdown
	graceful.Shutdown(app.FiberApp, app.DB)
}
