package middleware

import (
	"app/xonvera-core/internal/adapters/handler/http"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// NewRateLimiter creates a rate limiting middleware with Redis storage
func NewRateLimiter(max int, duration time.Duration, redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := fmt.Sprintf("ratelimit:%s", c.IP())
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// Increment counter
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			// If Redis fails, allow the request to proceed
			return c.Next()
		}

		// Set expiration on first request
		if count == 1 {
			redisClient.Expire(ctx, key, duration)
		}

		// Check if limit exceeded
		if count > int64(max) {
			return http.ErrorLimited(c, []string{"Rate limit exceeded. Please try again later"})
		}

		return c.Next()
	}
}

// AuthRateLimiter creates a stricter rate limiter for auth endpoints
func AuthRateLimiter(redisClient *redis.Client) fiber.Handler {
	return NewRateLimiter(5000, 15*time.Minute, redisClient) // 10 requests per 15 minutes
}

// APIRateLimiter creates a general rate limiter for API endpoints
func APIRateLimiter(redisClient *redis.Client) fiber.Handler {
	return NewRateLimiter(100, 1*time.Minute, redisClient) // 100 requests per minute
}
