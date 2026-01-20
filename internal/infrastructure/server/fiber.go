package server

import (
	"time"

	"app/xonvera-core/internal/adapters/middleware"
	"app/xonvera-core/internal/infrastructure/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
)

// NewFiberApp creates and configures a new Fiber application instance with middleware and routes
func NewFiberApp(cfg *config.Config, redisClient *redis.Client) *fiber.App {
	startTime := time.Now()

	app := fiber.New(fiber.Config{
		AppName: cfg.App.Name,
	})

	// Global middleware
	app.Use(middleware.RequestID())
	app.Use(recover.New())
	app.Use(middleware.BodyLogger())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "${time} | ${status} | ${latency} | ${locals:requestID} | ${method} ${path}\n",
	}))
	app.Use(middleware.APIRateLimiter(redisClient))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		uptime := time.Since(startTime)
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"app":     cfg.App.Name,
			"version": cfg.App.Version,
			"uptime":  uptime.String(),
		})
	})

	return app
}
