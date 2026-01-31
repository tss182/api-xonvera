package redis

import (
	"context"
	"fmt"
	"time"

	"app/xonvera-core/internal/infrastructure/config"
	"app/xonvera-core/internal/infrastructure/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// NewRedisClient creates a new Redis client with connection verification
func NewRedisClient(cfg *config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis",
			zap.String("host", cfg.Host),
			zap.String("port", cfg.Port),
			zap.Error(err),
		)
	}

	logger.Info("Connected to Redis",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.Int("db", cfg.DB),
	)

	return client
}

// CloseRedis closes the Redis connection
func CloseRedis(client *redis.Client) error {
	if client == nil {
		return nil
	}
	return client.Close()
}
