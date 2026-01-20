package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestID middleware adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request ID already exists in header
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			// Generate new UUID if not present
			requestID = uuid.New().String()
		}

		// Set request ID in context and response header
		c.Locals("requestID", requestID)
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(c *fiber.Ctx) string {
	if id, ok := c.Locals("requestID").(string); ok {
		return id
	}
	return ""
}
