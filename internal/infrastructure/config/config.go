package config

import (
	"fmt"
	"time"

	"app/xonvera-core/internal/infrastructure/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type (
	Config struct {
		App        AppConfig        `mapstructure:",squash"`
		Database   DatabaseConfig   `mapstructure:",squash"`
		Token      TokenConfig      `mapstructure:",squash"`
		Redis      RedisConfig      `mapstructure:",squash"`
		Pagination PaginationConfig `mapstructure:",squash"`
	}

	AppConfig struct {
		Name           string `mapstructure:"APP_NAME"`
		Version        string `mapstructure:"APP_VERSION"`
		Port           string `mapstructure:"APP_PORT"`
		Env            string `mapstructure:"APP_ENV"`
		AllowedOrigins string `mapstructure:"APP_ALLOWED_ORIGINS"`
		RequestTimeout time.Duration
		CORSOrigins    string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	}

	DatabaseConfig struct {
		Host            string `mapstructure:"DB_HOST"`
		Port            string `mapstructure:"DB_PORT"`
		User            string `mapstructure:"DB_USER"`
		Password        string `mapstructure:"DB_PASSWORD"`
		Name            string `mapstructure:"DB_NAME"`
		SSLMode         string `mapstructure:"DB_SSLMODE"`
		MaxIdleConn     int    `mapstructure:"DB_OPT_MAX_IDLE_CONN"`
		MaxOpenConn     int    `mapstructure:"DB_OPT_MAX_OPEN_CONN"`
		ConnMaxIdleTime time.Duration
		ConnMaxLifeTime time.Duration
	}

	TokenConfig struct {
		SecretKey      string `mapstructure:"TOKEN_SECRET_KEY"`
		Expired        time.Duration
		RefreshExpired time.Duration
	}

	RedisConfig struct {
		Host     string `mapstructure:"REDIS_HOST"`
		Port     string `mapstructure:"REDIS_PORT"`
		Password string `mapstructure:"REDIS_PASSWORD"`
		DB       int    `mapstructure:"REDIS_DB"`
	}

	PaginationConfig struct {
		DefaultLimit int `mapstructure:"PAGINATION_DEFAULT_LIMIT"`
		MaxLimit     int `mapstructure:"PAGINATION_MAX_LIMIT"`
	}
)

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Warn("Error reading config file", zap.Error(err))
	}

	// Set defaults
	setDefaults()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Fatal("Failed to unmarshal config", zap.Error(err))
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Fatal("Invalid configuration", zap.Error(err))
	}

	// Parse duration configurations
	parseDurationConfigs(&cfg)

	return &cfg
}

// parseDurationConfigs parses all duration-based environment variables
func parseDurationConfigs(cfg *Config) {
	durationConfigs := []struct {
		envKey    string
		target    *time.Duration
		fieldName string
	}{
		{
			envKey:    "APP_REQUEST_TIMEOUT",
			target:    &cfg.App.RequestTimeout,
			fieldName: "APP_REQUEST_TIMEOUT",
		},
		{
			envKey:    "DB_OPT_CONN_MAX_IDLE_TIME",
			target:    &cfg.Database.ConnMaxIdleTime,
			fieldName: "DB_OPT_CONN_MAX_IDLE_TIME",
		},
		{
			envKey:    "DB_OPT_CONN_MAX_LIFE_TIME",
			target:    &cfg.Database.ConnMaxLifeTime,
			fieldName: "DB_OPT_CONN_MAX_LIFE_TIME",
		},
		{
			envKey:    "TOKEN_EXPIRE",
			target:    &cfg.Token.Expired,
			fieldName: "TOKEN_EXPIRE",
		},
		{
			envKey:    "TOKEN_REFRESH_EXPIRE",
			target:    &cfg.Token.RefreshExpired,
			fieldName: "TOKEN_REFRESH_EXPIRE",
		},
	}

	for _, dc := range durationConfigs {
		if timeoutStr := viper.GetString(dc.envKey); timeoutStr != "" {
			duration, err := time.ParseDuration(timeoutStr)
			if err != nil {
				logger.Fatal("Failed to parse duration config",
					zap.String("config", dc.fieldName),
					zap.String("value", timeoutStr),
					zap.Error(err),
				)
			}
			*dc.target = duration
		}
	}
}

// Validate checks configuration for production readiness
func (c *Config) Validate() error {
	if c.App.Env == "production" {
		// Validate secret key
		if c.Token.SecretKey == "" || len(c.Token.SecretKey) < 32 {
			return fmt.Errorf("TOKEN_SECRET_KEY must be at least 32 characters in production, got %d", len(c.Token.SecretKey))
		}
		// Don't allow default values in production
		if c.Token.SecretKey == "your-super-secret-key-min-32-chars!!" {
			return fmt.Errorf("TOKEN_SECRET_KEY must be changed from default value in production")
		}
		// Validate database SSL
		if c.Database.SSLMode == "disable" {
			return fmt.Errorf("DB_SSLMODE must not be 'disable' in production")
		}
	}
	return nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("APP_NAME", "xonvera")
	viper.SetDefault("APP_PORT", "3000")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_REQUEST_TIMEOUT", "30s")
	viper.SetDefault("APP_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080")

	// Database defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "xonvera_db")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_OPT_MAX_IDLE_CONN", 10)
	viper.SetDefault("DB_OPT_MAX_OPEN_CONN", 100)
	viper.SetDefault("DB_OPT_CONN_MAX_IDLE_TIME", "10m")
	viper.SetDefault("DB_OPT_CONN_MAX_LIFE_TIME", "1h")

	// Token defaults
	viper.SetDefault("TOKEN_SECRET_KEY", "your-super-secret-key-min-32-chars!!")
	viper.SetDefault("TOKEN_EXPIRE", "24h")
	viper.SetDefault("TOKEN_REFRESH_EXPIRE", "168h")

	// Redis defaults
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0) // 7 days
}
