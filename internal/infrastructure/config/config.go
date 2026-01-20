package config

import (
	"time"

	"app/xonvera-core/internal/infrastructure/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	App      AppConfig      `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
	Token    TokenConfig    `mapstructure:",squash"`
	Redis    RedisConfig    `mapstructure:",squash"`
}

type AppConfig struct {
	Name           string `mapstructure:"APP_NAME"`
	Version        string `mapstructure:"APP_VERSION"`
	Port           string `mapstructure:"APP_PORT"`
	Env            string `mapstructure:"APP_ENV"`
	RequestTimeout time.Duration
}

type DatabaseConfig struct {
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

type TokenConfig struct {
	SecretKey      string `mapstructure:"TOKEN_SECRET_KEY"`
	Expired        time.Duration
	RefreshExpired time.Duration
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	// viper.SetConfigType("env")
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

	// Parse duration strings
	if timeoutStr := viper.GetString("APP_REQUEST_TIMEOUT"); timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			logger.Fatal("Failed to parse APP_REQUEST_TIMEOUT", zap.Error(err))
		}
		cfg.App.RequestTimeout = timeout
	}

	if idleTimeStr := viper.GetString("DB_OPT_CONN_MAX_IDLE_TIME"); idleTimeStr != "" {
		idleTime, err := time.ParseDuration(idleTimeStr)
		if err != nil {
			logger.Fatal("Failed to parse DB_OPT_CONN_MAX_IDLE_TIME", zap.Error(err))
		}
		cfg.Database.ConnMaxIdleTime = idleTime
	}

	if lifeTimeStr := viper.GetString("DB_OPT_CONN_MAX_LIFE_TIME"); lifeTimeStr != "" {
		lifeTime, err := time.ParseDuration(lifeTimeStr)
		if err != nil {
			logger.Fatal("Failed to parse DB_OPT_CONN_MAX_LIFE_TIME", zap.Error(err))
		}
		cfg.Database.ConnMaxLifeTime = lifeTime
	}

	if expire := viper.GetString("TOKEN_EXPIRE"); expire != "" {
		duration, err := time.ParseDuration(expire)
		if err != nil {
			logger.Fatal("Failed to parse TOKEN_EXPIRE", zap.Error(err))
		}
		cfg.Token.Expired = duration
	}

	if refreshExpire := viper.GetString("TOKEN_REFRESH_EXPIRE"); refreshExpire != "" {
		duration, err := time.ParseDuration(refreshExpire)
		if err != nil {
			logger.Fatal("Failed to parse TOKEN_REFRESH_EXPIRE", zap.Error(err))
		}
		cfg.Token.RefreshExpired = duration
	}

	return &cfg
}

func setDefaults() {
	// App defaults
	viper.SetDefault("APP_NAME", "xonvera")
	viper.SetDefault("APP_PORT", "3000")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_REQUEST_TIMEOUT", "30s")

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
