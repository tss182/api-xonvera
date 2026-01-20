package redis

import (
	"context"
	"fmt"

	"app/xonvera-core/internal/infrastructure/config"
	"app/xonvera-core/internal/infrastructure/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*1000000000)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	logger.Info("Connected to Redis",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
	)

	return client
}

// CloseRedis closes the Redis connection
func CloseRedis(client *redis.Client) error {
	return client.Close()
}
