package main

import (
	"flag"

	"app/xonvera-core/internal/adapters/middleware"
	"app/xonvera-core/internal/adapters/routes"
	"app/xonvera-core/internal/dependencies"
	"app/xonvera-core/internal/infrastructure/database"
	"app/xonvera-core/internal/infrastructure/graceful"
	"app/xonvera-core/internal/infrastructure/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
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

	// Initialize Fiber app
	fiberApp := fiber.New(fiber.Config{
		AppName: app.Config.App.Name,
	})

	// Global middleware
	fiberApp.Use(middleware.RequestID())
	fiberApp.Use(recover.New())
	fiberApp.Use(middleware.BodyLogger(app.Config.App.Env))
	fiberApp.Use(fiberlogger.New(fiberlogger.Config{
		Format: "${time} | ${status} | ${latency} | ${locals:request_id} | ${method} ${path}\n",
	}))
	fiberApp.Use(middleware.APIRateLimiter(app.Redis))
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: app.Config.App.CORSOrigins,
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
	routes.SetupRoutes(fiberApp, app.AuthHandler, app.AuthMiddleware, app.Redis)

	// Start server in a goroutine
	addr := ":" + app.Config.App.Port
	logger.Info("Starting server",
		zap.String("app", app.Config.App.Name),
		zap.String("addr", addr),
		zap.String("env", app.Config.App.Env),
	)

	// Channel to handle server errors
	serverErrors := make(chan error, 1)

	go func() {
		if err := fiberApp.Listen(addr); err != nil {
			serverErrors <- err
		}
	}()

	// Handle graceful shutdown or server errors
	graceful.Shutdown(fiberApp, app.DB, []*redis.Client{app.Redis}, serverErrors)
}
