//go:build wireinject
// +build wireinject

package dependencies

import (
	"time"

	"app/xonvera-core/internal/adapters/handler/http"
	"app/xonvera-core/internal/adapters/middleware"
	repositoriesRedis "app/xonvera-core/internal/adapters/repositories/redis"
	repositoriesSql "app/xonvera-core/internal/adapters/repositories/sql"
	"app/xonvera-core/internal/core/services"
	"app/xonvera-core/internal/infrastructure/config"
	"app/xonvera-core/internal/infrastructure/database"
	"app/xonvera-core/internal/infrastructure/redis"
	"app/xonvera-core/internal/infrastructure/server"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ProviderSet is the set of all providers
var ProviderSet = wire.NewSet(
	// Config
	config.LoadConfig,
	ProvideDBConfig,
	ProvideTokenConfig,
	ProvideRedisConfig,
	ProvideRequestTimeout,

	// Database
	database.NewConnection,

	// Redis
	redis.NewRedisClient,

	// Server
	server.NewFiberApp,

	// Repositories
	repositoriesSql.NewUserRepository,
	repositoriesSql.NewPackageRepository,
	repositoriesSql.NewInvoiceRepository,
	repositoriesSql.NewTxRepository,
	repositoriesRedis.NewTokenRepository,

	// Services
	services.NewTokenService,
	services.NewAuthService,
	services.NewPackageService,
	services.NewInvoiceService,

	// Handlers
	http.NewAuthHandler,
	http.NewPackageHandler,
	http.NewInvoiceHandler,

	// Middleware
	middleware.NewAuthMiddleware,
)

// ProvideDBConfig extracts DatabaseConfig from Config
func ProvideDBConfig(cfg *config.Config) *config.DatabaseConfig {
	return &cfg.Database
}

// ProvideTokenConfig extracts TokenConfig from Config
func ProvideTokenConfig(cfg *config.Config) *config.TokenConfig {
	return &cfg.Token
}

// ProvideRedisConfig extracts RedisConfig from Config
func ProvideRedisConfig(cfg *config.Config) *config.RedisConfig {
	return &cfg.Redis
}

// ProvideRequestTimeout extracts request timeout from Config
func ProvideRequestTimeout(cfg *config.Config) time.Duration {
	return cfg.App.RequestTimeout
}

// Application holds all the dependencies
type Application struct {
	Config         *config.Config
	DB             *gorm.DB
	Redis          *goredis.Client
	FiberApp       *fiber.App
	AuthHandler    *http.AuthHandler
	PackageHandler *http.PackageHandler
	InvoiceHandler *http.InvoiceHandler
	AuthMiddleware *middleware.AuthMiddleware
}

// InitializeApplication creates a new Application with all dependencies wired
func InitializeApplication() (*Application, error) {
	wire.Build(
		ProviderSet,
		wire.Struct(new(Application), "*"),
	)
	return nil, nil
}
