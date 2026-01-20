# Redis Integration with Google Wire

This document explains how Redis is integrated with Google Wire for dependency injection, allowing all services to access Redis easily.

## Architecture

Redis is now managed by Google Wire and available in the `Application` struct:

```go
type Application struct {
    Config         *config.Config
    DB             *gorm.DB
    Redis          *redis.Client        // Redis client injected by Wire
    AuthHandler    *http.AuthHandler
    AuthMiddleware *middleware.AuthMiddleware
}
```

## How It Works

1. **Configuration** - Redis config is loaded from environment variables in `internal/infrastructure/config/config.go`
2. **Initialization** - Wire creates the Redis client in `InitializeApplication()` via `redis.NewRedisClient()`
3. **Dependency Injection** - All services that need Redis receive it through Wire

## Using Redis in Your Services

### Example 1: Accessing Redis from Routes

```go
// Routes automatically receive Redis from the Application
routes.SetupRoutes(fiberApp, app.AuthHandler, app.AuthMiddleware, app.Redis)

// Inside routes/routes.go
func SetupRoutes(
    app *fiber.App,
    authHandler *http.AuthHandler,
    authMiddleware *middleware.AuthMiddleware,
    redisClient *redis.Client,  // Redis injected here
) {
    // Use Redis for rate limiting
    auth := app.Group("/auth", middleware.AuthRateLimiter(redisClient))
}
```

### Example 2: Adding Redis to a New Service

1. **Update wire.go** to add your service as a provider:

```go
var ProviderSet = wire.NewSet(
    // ... existing providers ...
    myservice.NewMyService,  // Add your service
)
```

2. **Create your service with Redis dependency**:

```go
package myservice

import (
    "github.com/redis/go-redis/v9"
)

type MyService struct {
    redis *redis.Client
}

func NewMyService(redis *redis.Client) *MyService {
    return &MyService{
        redis: redis,
    }
}

func (s *MyService) SomeMethod(ctx context.Context) error {
    // Use Redis
    err := s.redis.Set(ctx, "key", "value", time.Hour).Err()
    return err
}
```

3. **Regenerate Wire dependencies**:

```bash
make wire
```

### Example 3: Using TokenRedisRepository

```go
// TokenRedisRepository is available for token management
tokenRepo := repositories.NewTokenRedisRepository(app.Redis)

// Invalidate a token on logout
tokenRepo.InvalidateToken(ctx, accessToken, tokenExpiry)

// Check if token is valid
isValid, err := tokenRepo.IsTokenValid(ctx, token)
```

## Available Redis Methods

The Redis client provides standard operations:

```go
// Set a key with expiration
app.Redis.Set(ctx, "key", "value", 1*time.Hour).Err()

// Get a value
val, err := app.Redis.Get(ctx, "key").Result()

// Increment counter (for rate limiting)
count, err := app.Redis.Incr(ctx, "counter").Result()

// Expire a key
app.Redis.Expire(ctx, "key", 1*time.Hour).Err()

// Delete a key
app.Redis.Del(ctx, "key").Err()

// Scan keys
iter := app.Redis.Scan(ctx, 0, "pattern:*", 0).Iterator()
```

## Wire Dependency Flow

```
config.LoadConfig()
    ↓
config.RedisConfig
    ↓
redis.NewRedisClient()  ← Handles connection & initialization
    ↓
Application.Redis       ← Available in Application struct
    ↓
All Services            ← Can receive Redis through dependency injection
```

## Environment Configuration

```dotenv
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=              # Leave empty for no password
REDIS_DB=0                   # Database number
```

## Error Handling

If Redis is unavailable:
- Rate limiter will allow requests to proceed (graceful degradation)
- Token operations will return errors which should be handled by services

## Cleanup

Redis connection is automatically closed when the application stops:

```go
defer redis.CloseRedis(app.Redis)
```

## Testing

To test with Redis locally:

```bash
# Start Redis container
docker run -d -p 6379:6379 redis:latest

# Or use Homebrew
brew services start redis

# Run app
make dev
```

## Benefits

✅ **Centralized Management** - Single point of Redis configuration
✅ **Dependency Injection** - Clean architecture with Wire
✅ **Type Safety** - Redis client is properly typed
✅ **Easy Testing** - Can mock Redis in tests
✅ **Scalability** - Works across multiple instances with same Redis server
✅ **Connection Pooling** - Automatically managed by go-redis
