package middleware

import (
	"app/xonvera-core/internal/adapters/handler/http"
	"app/xonvera-core/internal/infrastructure/logger"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// Rate limiter constants
	AuthRateLimitRequests = 10
	AuthRateLimitDuration = 15 * time.Minute
	APIRateLimitRequests  = 100
	APIRateLimitDuration  = 1 * time.Minute
	RedisContextTimeout   = 1 * time.Second
)

// Lua script for atomic INCR + EXPIRE operation
const rateLimitLuaScript = `
	local current = redis.call('INCR', KEYS[1])
	if current == 1 then
		redis.call('EXPIRE', KEYS[1], ARGV[1])
	end
	return current
`

// NewRateLimiter creates a rate limiting middleware with Redis storage
func NewRateLimiter(max int, duration time.Duration, redisClient *redis.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Resolve rate limit key from userID or IP
		keyID := resolveRateLimitKey(c)
		key := fmt.Sprintf("ratelimit:%s", keyID)

		ctx, cancel := context.WithTimeout(context.Background(), RedisContextTimeout)
		defer cancel()

		// Execute atomic increment with expiration
		count, err := redisClient.Eval(ctx, rateLimitLuaScript, []string{key}, int(duration.Seconds())).Int64()
		if err != nil {
			// If Redis fails, allow the request to proceed (fail-open)
			logger.ContextWarn(c, "Rate limiter Redis error, allowing request",
				zap.String("key", key),
				zap.Error(err),
			)
			return c.Next()
		}

		// Check if limit exceeded
		if count > int64(max) {
			logger.ContextDebug(c, "Rate limit exceeded",
				zap.String("key", keyID),
				zap.Int64("count", count),
				zap.Int("limit", max),
			)
			return http.ErrorLimited(c, []string{"Rate limit exceeded. Please try again later"})
		}

		return c.Next()
	}
}

// AuthRateLimiter creates a stricter rate limiter for auth endpoints
func AuthRateLimiter(redisClient *redis.Client) fiber.Handler {
	return NewRateLimiter(AuthRateLimitRequests, AuthRateLimitDuration, redisClient)
}

// APIRateLimiter creates a general rate limiter for API endpoints
func APIRateLimiter(redisClient *redis.Client) fiber.Handler {
	return NewRateLimiter(APIRateLimitRequests, APIRateLimitDuration, redisClient)
}

// resolveRateLimitKey builds the rate-limit key based on userID or IP
func resolveRateLimitKey(c fiber.Ctx) string {
	if userID := c.Locals("userID"); userID != nil {
		if v := fmt.Sprint(userID); v != "" {
			return v
		}
	}

	// Fallback to IP address
	return c.IP()
}
