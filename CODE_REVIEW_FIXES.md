# Code Review - Issues & Recommended Fixes

## 🔴 CRITICAL - Must Fix Before Production

### 1. Weak Default Secret Key
**File:** `internal/infrastructure/config/config.go:137`
**Issue:** Default secret key is weak and committed to git
```go
// ❌ Current
viper.SetDefault("TOKEN_SECRET_KEY", "your-super-secret-key-min-32-chars!!")

// ✅ Fix: Require secret key in production
if cfg.App.Env == "production" && cfg.Token.SecretKey == "your-super-secret-key-min-32-chars!!" {
    logger.Fatal("TOKEN_SECRET_KEY must be set in production environment")
}
```

### 2. CORS Wildcard Configuration
**File:** `internal/infrastructure/server/fiber.go:30`
**Issue:** Allows any origin to access API
```go
// ❌ Current
AllowOrigins: "*"

// ✅ Fix: Make configurable
type AppConfig struct {
    // ... existing fields
    AllowedOrigins string `mapstructure:"APP_ALLOWED_ORIGINS"`
}
// In fiber.go:
AllowOrigins: cfg.App.AllowedOrigins // e.g., "https://yourdomain.com,https://app.yourdomain.com"
```

### 3. Password Logging Risk
**File:** `internal/adapters/middleware/body_logger.go`
**Issue:** Request body logging might expose passwords
```go
// ✅ Already partially implemented with maskSensitiveFields
// Ensure these fields are always masked:
sensitiveFields := []string{"password", "token", "secret", "authorization", "refresh_token", "access_token"}
```

---

## 🟠 HIGH PRIORITY - Bugs & Stability

### 4. Race Condition in Rate Limiter
**File:** `internal/adapters/middleware/rate_limiter_redis.go:18-27`
**Issue:** INCR and EXPIRE not atomic
```go
// ❌ Current
count, err := redisClient.Incr(ctx, key).Result()
if count == 1 {
    redisClient.Expire(ctx, key, duration)
}

// ✅ Fix: Use Lua script for atomic operation
const luaScript = `
local current = redis.call('INCR', KEYS[1])
if current == 1 then
    redis.call('EXPIRE', KEYS[1], ARGV[1])
end
return current
`
count, err := redisClient.Eval(ctx, luaScript, []string{key}, int(duration.Seconds())).Int64()
```

### 5. Wrong Timeout Constant
**File:** `internal/infrastructure/redis/redis.go:20`
```go
// ❌ Current
ctx, cancel := context.WithTimeout(context.Background(), 5*1000000000)

// ✅ Fix
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
```

### 6. Graceful Shutdown Not Implemented
**File:** `cmd/main.go`
```go
// ✅ Add at end of main():
// Start server in goroutine
go func() {
    addr := ":" + app.Config.App.Port
    logger.Info("Starting server", zap.String("addr", addr))
    if err := app.FiberApp.Listen(addr); err != nil {
        logger.Fatal("Failed to start server", zap.Error(err))
    }
}()

// Handle graceful shutdown
graceful.Shutdown(app.FiberApp, app.DB)
```

### 7. Unsafe Type Assertion
**File:** `internal/adapters/routes/routes.go:36`
```go
// ❌ Current
userID := c.Locals("userID").(uint)

// ✅ Fix
userID, ok := c.Locals("userID").(uint)
if !ok {
    return http.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user context")
}
```

### 8. Token Update Not Atomic
**File:** `internal/adapters/repositories/redis/token.go`
```go
// ✅ Use Redis transaction (MULTI/EXEC)
func (r *TokenRepository) Update(ctx context.Context, token *domain.Token) error {
    pipe := r.client.TxPipeline()
    
    // Delete old keys
    oldAccessKey := fmt.Sprintf("token:access:%s", token.AccessToken)
    oldRefreshKey := fmt.Sprintf("token:refresh:%s", token.RefreshToken)
    pipe.Del(ctx, oldAccessKey, oldRefreshKey)
    
    // Set new keys
    data, _ := json.Marshal(token)
    accessKey := fmt.Sprintf("token:access:%s", token.AccessToken)
    refreshKey := fmt.Sprintf("token:refresh:%s", token.RefreshToken)
    pipe.Set(ctx, accessKey, data, time.Until(token.ExpiresAt))
    pipe.Set(ctx, refreshKey, data, time.Until(token.RefreshExpiresAt))
    
    _, err := pipe.Exec(ctx)
    return err
}
```

---

## 🟡 MEDIUM PRIORITY - Improvements

### 9. Add Request Size Limits
**File:** `internal/infrastructure/server/fiber.go`
```go
app := fiber.New(fiber.Config{
    AppName:       cfg.App.Name,
    BodyLimit:     4 * 1024 * 1024,  // 4MB max
    ReadTimeout:   time.Second * 10,
    WriteTimeout:  time.Second * 10,
    StrictRouting: true,
    CaseSensitive: true,
})
```

### 10. Enhanced Health Check
**File:** `internal/infrastructure/server/fiber.go`
```go
app.Get("/health", func(c *fiber.Ctx) error {
    uptime := time.Since(startTime)
    
    // Check DB
    sqlDB, _ := db.DB()
    dbHealth := "healthy"
    if err := sqlDB.Ping(); err != nil {
        dbHealth = "unhealthy"
    }
    
    // Check Redis
    redisHealth := "healthy"
    if err := redisClient.Ping(c.Context()).Err(); err != nil {
        redisHealth = "unhealthy"
    }
    
    return c.JSON(fiber.Map{
        "status":  "healthy",
        "app":     cfg.App.Name,
        "version": cfg.App.Version,
        "uptime":  uptime.String(),
        "checks": fiber.Map{
            "database": dbHealth,
            "redis":    redisHealth,
        },
    })
})
```

### 11. Better Error Messages (Security)
**File:** `internal/core/services/auth.go`
```go
// ✅ Current approach is mostly good - keep generic messages
// But ensure consistency:
if exists {
    return nil, errors.New("username already taken") // Generic, don't specify email/phone
}
```

### 12. Phone Validation
**File:** `internal/core/domain/user.go`
```go
// Add custom validator
type RegisterRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email,max=255"`
    Phone    string `json:"phone" validate:"required,e164"` // E.164 format
    Password string `json:"password" validate:"required,min=8,max=100,password"` // Stronger min
}
```

### 13. Logger Early Initialization
**File:** `cmd/main.go`
```go
func main() {
    // Initialize logger FIRST with default settings
    logger.Init("development") // Default
    defer logger.Sync()
    
    // Load config
    app, err := dependencies.InitializeApplication()
    // ...
    
    // Re-initialize logger with correct env
    logger.Init(app.Config.App.Env)
}
```

### 14. Add Database Indexes
**Migration:** Create new migration
```sql
-- Add composite index for common queries
CREATE INDEX idx_users_email_deleted ON auth.users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_phone_deleted ON auth.users(phone) WHERE deleted_at IS NULL;
```

### 15. Environment Variable Validation
**File:** `internal/infrastructure/config/config.go`
```go
func (c *Config) Validate() error {
    if c.App.Env == "production" {
        if c.Token.SecretKey == "" || len(c.Token.SecretKey) < 32 {
            return errors.New("TOKEN_SECRET_KEY must be at least 32 characters in production")
        }
        if c.Database.SSLMode == "disable" {
            return errors.New("DB_SSLMODE=disable not allowed in production")
        }
    }
    return nil
}
```

---

## 📊 Performance Optimizations

### 16. Connection Pool Monitoring
```go
// Add periodic logging in main.go
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    for range ticker.C {
        sqlDB, _ := app.DB.DB()
        stats := sqlDB.Stats()
        logger.Info("DB Pool Stats",
            zap.Int("open_connections", stats.OpenConnections),
            zap.Int("in_use", stats.InUse),
            zap.Int("idle", stats.Idle),
        )
    }
}()
```

### 17. Redis Token Storage Optimization
Consider using Redis Hash instead of 3 separate keys:
```go
// Store as hash for one user
HSET token:user:123 access "token..." refresh "token..." expires_at "timestamp"
```

---

## 📝 Documentation Needed

1. **README.md** - Setup instructions, environment variables, running migrations
2. **API Documentation** - Swagger/OpenAPI spec
3. **Architecture Diagram** - Hexagonal architecture explanation
4. **SECURITY.md** - Security best practices, reporting vulnerabilities
5. **.env.example** - Template with safe defaults

---

## ✅ What You Did Well

1. ✅ **Hexagonal Architecture** - Clean separation of concerns
2. ✅ **Dependency Injection** - Using Wire properly
3. ✅ **Middleware Organization** - Well-structured
4. ✅ **Password Hashing** - Using bcrypt correctly
5. ✅ **Context Timeouts** - Implemented throughout
6. ✅ **PASETO Tokens** - Better than JWT for your use case
7. ✅ **Structured Logging** - Using Zap properly
8. ✅ **Database Migrations** - Proper migration system
9. ✅ **Redis for Sessions** - Good choice for token storage
10. ✅ **Rate Limiting** - Implemented (though needs atomic fix)

---

## Priority Order

1. Fix secret key validation (#1)
2. Fix rate limiter race condition (#4)
3. Fix Redis timeout constant (#5)
4. Implement graceful shutdown (#6)
5. Fix type assertion panic (#7)
6. Configure CORS properly (#2)
7. Add request size limits (#9)
8. Enhance health checks (#10)
9. Add database indexes (#14)
10. Write documentation
