package middleware

import (
	"context"
	"strings"
	"time"

	"app/xonvera-core/internal/adapters/handler/http"
	"app/xonvera-core/internal/core/ports/service"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	authService    portService.AuthService
	requestTimeout time.Duration
}

func NewAuthMiddleware(authService portService.AuthService, requestTimeout time.Duration) *AuthMiddleware {
	return &AuthMiddleware{
		authService:    authService,
		requestTimeout: requestTimeout,
	}
}

func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return http.ErrorResponse(c, fiber.StatusUnauthorized, "Missing authorization header")
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return http.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		token := parts[1]

		// Validate token format and check database
		ctx, cancel := context.WithTimeout(c.UserContext(), m.requestTimeout)
		defer cancel()

		userID, err := m.authService.ValidateAccessToken(ctx, token)
		if err != nil {
			return http.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Set user ID and token in context for downstream handlers
		c.Locals("userID", userID)
		c.Locals("accessToken", token)

		return c.Next()
	}
}
