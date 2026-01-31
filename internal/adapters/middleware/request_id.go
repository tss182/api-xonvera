package middleware

import (
	"app/xonvera-core/internal/infrastructure/logger"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestID middleware adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request ID already exists in header
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			// Generate new UUID v7
			uuidV7, err := uuid.NewV7()
			if err != nil {
				// Fallback to UUID v4
				uuidV7 = uuid.New()
				logger.Debug("Failed to generate UUIDv7, using UUIDv4", zap.Error(err))
			}
			requestID = uuidV7.String()
		}

		// Set request ID in Fiber context
		c.Locals("request_id", requestID)

		// Store in standard context for service layer access
		ctx := context.WithValue(c.UserContext(), "request_id", requestID)
		c.SetUserContext(ctx)

		// Set response header
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(c *fiber.Ctx) string {
	if id, ok := c.Locals("request_id").(string); ok {
		return id
	}
	return ""
}
