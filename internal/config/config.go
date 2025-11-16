package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                string
	DatabaseURL         string
	DBMaxOpenConns      int
	DBMaxIdleConns      int
	DBConnMaxLifetime   time.Duration
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
	RateLimitMax        int
	RateLimitExpiration time.Duration
	CORSOrigins         string
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "3000"),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		DBMaxOpenConns:      getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:      getEnvInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime:   getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		ReadTimeout:         getEnvDuration("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:        getEnvDuration("WRITE_DURATION", 10*time.Second),
		IdleTimeout:         getEnvDuration("IDLE_TIMEOUT", 120*time.Second),
		RateLimitMax:        getEnvInt("RATE_LIMIT_MAX", 100),
		RateLimitExpiration: getEnvDuration("RATE_LIMIT_EXPIRATION", 1*time.Minute),
		CORSOrigins:         getEnv("CORS_ORIGINS", "*"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
