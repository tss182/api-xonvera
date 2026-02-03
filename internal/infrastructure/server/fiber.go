package server

import (
	"time"

	"app/xonvera-core/internal/adapters/middleware"
	"app/xonvera-core/internal/infrastructure/config"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewFiberApp creates and configures a new Fiber application with middleware
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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	app.Get("/health", createHealthCheckHandler(startTime, db, redisClient, cfg))

	return app
}

// createHealthCheckHandler returns a health check handler
func createHealthCheckHandler(startTime time.Time, db *gorm.DB, redisClient *redis.Client, cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		uptime := time.Since(startTime)

		// Check database connectivity
		dbHealth := checkDatabaseHealth(db)

		// Check Redis connectivity
		redisHealth := checkRedisHealth(c, redisClient)

		// Determine overall status
		overallStatus := "healthy"
		statusCode := fiber.StatusOK

		if dbHealth != "healthy" || redisHealth != "healthy" {
			overallStatus = "degraded"
			statusCode = fiber.StatusServiceUnavailable
		}

		response := fiber.Map{
			"status":  overallStatus,
			"app":     cfg.App.Name,
			"version": cfg.App.Version,
			"uptime":  uptime.String(),
			"checks": fiber.Map{
				"database": dbHealth,
				"redis":    redisHealth,
			},
		}

		return c.Status(statusCode).JSON(response)
	}
}

// checkDatabaseHealth verifies database connectivity
func checkDatabaseHealth(db *gorm.DB) string {
	if db == nil {
		return "unavailable"
	}

	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		return "unhealthy"
	}

	return "healthy"
}

// checkRedisHealth verifies Redis connectivity
func checkRedisHealth(c fiber.Ctx, redisClient *redis.Client) string {
	if redisClient == nil {
		return "unavailable"
	}

	if err := redisClient.Ping(c.Context()).Err(); err != nil {
		return "unhealthy"
	}

	return "healthy"
}

// fiber:context-methods migrated
