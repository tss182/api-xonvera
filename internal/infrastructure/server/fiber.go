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
	"gorm.io/gorm"
)

// NewFiberApp creates and configures a new Fiber application instance with middleware and routes
func NewFiberApp(cfg *config.Config, redisClient *redis.Client, db *gorm.DB) *fiber.App {
	startTime := time.Now()

	app := fiber.New(fiber.Config{
		AppName:       cfg.App.Name,
		BodyLimit:     4 * 1024 * 1024, // 4MB max body size
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  10 * time.Second,
		StrictRouting: true,
		CaseSensitive: true,
	})

	// Global middleware
	app.Use(middleware.RequestID())
	app.Use(recover.New())
	app.Use(middleware.BodyLogger(cfg.App.Env))
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "${time} | ${status} | ${latency} | ${locals:requestID} | ${method} ${path}\n",
	}))
	app.Use(middleware.APIRateLimiter(redisClient))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.App.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		uptime := time.Since(startTime)

		// Check database connectivity
		dbHealth := "healthy"
		if db != nil {
			sqlDB, err := db.DB()
			if err != nil || sqlDB.Ping() != nil {
				dbHealth = "unhealthy"
			}
		}

		// Check Redis connectivity
		redisHealth := "healthy"
		if err := redisClient.Ping(c.Context()).Err(); err != nil {
			redisHealth = "unhealthy"
		}

		overallStatus := "healthy"
		if dbHealth != "healthy" || redisHealth != "healthy" {
			overallStatus = "degraded"
			// Return 503 if dependencies are unhealthy
			c.Status(fiber.StatusServiceUnavailable)
		}

		return c.JSON(fiber.Map{
			"status":  overallStatus,
			"app":     cfg.App.Name,
			"version": cfg.App.Version,
			"uptime":  uptime.String(),
			"checks": fiber.Map{
				"database": dbHealth,
				"redis":    redisHealth,
			},
		})
	})

	return app
}
