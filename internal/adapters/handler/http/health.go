package http

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status   string         `json:"status"`
	Version  string         `json:"version"`
	Database DatabaseHealth `json:"database"`
	Uptime   string         `json:"uptime"`
}

type DatabaseHealth struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

var startTime = time.Now()

// Health performs a comprehensive health check
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	dbHealth := h.checkDatabase(ctx)

	status := "healthy"
	if dbHealth.Status != "healthy" {
		status = "degraded"
	}

	response := HealthResponse{
		Status:   status,
		Version:  "1.0.0",
		Database: dbHealth,
		Uptime:   time.Since(startTime).String(),
	}

	statusCode := fiber.StatusOK
	if status == "degraded" {
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(response)
}

// checkDatabase checks database connectivity
func (h *HealthHandler) checkDatabase(ctx context.Context) DatabaseHealth {
	sqlDB, err := h.db.DB()
	if err != nil {
		return DatabaseHealth{
			Status: "unhealthy",
			Error:  err.Error(),
		}
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return DatabaseHealth{
			Status: "unhealthy",
			Error:  err.Error(),
		}
	}

	return DatabaseHealth{
		Status: "healthy",
	}
}
