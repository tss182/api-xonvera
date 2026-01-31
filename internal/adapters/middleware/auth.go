package middleware

import (
	"app/xonvera-core/internal/adapters/handler/http"
	portService "app/xonvera-core/internal/core/ports/service"
	"app/xonvera-core/internal/infrastructure/logger"
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	authorizationHeader = "Authorization"
	bearerScheme        = "bearer"
	userIDContextKey    = "userID"
	accessTokenKey      = "accessToken"
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

// Authenticate validates the access token from the Authorization header
func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(authorizationHeader)
		if authHeader == "" {
			logger.ContextDebug(c, "Missing authorization header")
			return http.NoAuth(c)
		}

		// Extract Bearer token
		token, err := extractBearerToken(authHeader)
		if err != nil {
			logger.ContextDebug(c, "Invalid authorization header format")
			return http.NoAuth(c)
		}

		// Validate token with timeout
		ctx, cancel := context.WithTimeout(c.UserContext(), m.requestTimeout)
		defer cancel()

		userID, err := m.authService.ValidateAccessToken(ctx, token)
		if err != nil {
			logger.ContextDebug(c, "Token validation failed",
				zap.Error(err),
			)
			return http.NoAuth(c)
		}

		// Set user ID and token in context for downstream handlers
		c.Locals(userIDContextKey, userID)
		c.Locals(accessTokenKey, token)

		return c.Next()
	}
}

// extractBearerToken extracts the token from a Bearer authorization header
func extractBearerToken(authHeader string) (string, error) {
	parts := strings.Fields(authHeader)
	if len(parts) != 2 {
		return "", fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header")
	}

	if strings.ToLower(parts[0]) != bearerScheme {
		return "", fiber.NewError(fiber.StatusUnauthorized, "invalid authorization scheme")
	}

	return parts[1], nil
}
