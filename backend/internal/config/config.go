package config

import (
	"os"
	"strconv"
	"time"
)

// Config はアプリケーション全体の設定を保持します
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Log      LogConfig
}

// ServerConfig はサーバー関連の設定です
type ServerConfig struct {
	Env          string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig はデータベース関連の設定です
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig はRedis関連の設定です
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig はJWT認証関連の設定です
type JWTConfig struct {
	Secret                   string
	AccessTokenExpiration    time.Duration
	RefreshTokenExpiration   time.Duration
	RefreshTokenRotation     bool
	RefreshTokenReuseWindow  time.Duration
	AccessTokenCookieDomain  string
	RefreshTokenCookieDomain string
}

// LogConfig はログ関連の設定です
type LogConfig struct {
	Level      string
	Format     string
	OutputPath string
}

// Load は環境変数から設定を読み込みます
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Env:          getEnv("ENV", "development"),
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			Name:            getEnv("DB_NAME", "effisio_dev"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:                   getEnv("JWT_SECRET", "your-secret-key-change-this"),
			AccessTokenExpiration:    getDurationEnv("JWT_ACCESS_TOKEN_EXPIRATION", 15*time.Minute),
			RefreshTokenExpiration:   getDurationEnv("JWT_REFRESH_TOKEN_EXPIRATION", 7*24*time.Hour),
			RefreshTokenRotation:     getBoolEnv("JWT_REFRESH_TOKEN_ROTATION", true),
			RefreshTokenReuseWindow:  getDurationEnv("JWT_REFRESH_TOKEN_REUSE_WINDOW", 10*time.Second),
			AccessTokenCookieDomain:  getEnv("JWT_ACCESS_TOKEN_COOKIE_DOMAIN", ""),
			RefreshTokenCookieDomain: getEnv("JWT_REFRESH_TOKEN_COOKIE_DOMAIN", ""),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			OutputPath: getEnv("LOG_OUTPUT_PATH", "stdout"),
		},
	}
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返します
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv は環境変数を整数として取得します
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getBoolEnv は環境変数を真偽値として取得します
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getDurationEnv は環境変数を時間として取得します
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		// 秒数として解釈を試みる
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}
