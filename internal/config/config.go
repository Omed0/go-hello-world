package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Port            string
	DatabaseURL     string
	ServerTimeout   time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() *Config {
	return &Config{
		Port:            getEnvOrDefault("PORT", "8080"),
		DatabaseURL:     getEnvOrDefault("DB_URL", ""),
		ServerTimeout:   getEnvDurationOrDefault("SERVER_TIMEOUT", 30*time.Second),
		ReadTimeout:     getEnvDurationOrDefault("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    getEnvDurationOrDefault("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:     getEnvDurationOrDefault("IDLE_TIMEOUT", 60*time.Second),
		MaxOpenConns:    getEnvIntOrDefault("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvIntOrDefault("DB_MAX_IDLE_CONNS", 25),
		ConnMaxLifetime: getEnvDurationOrDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute),
	}
}

// Helper functions for environment variable parsing
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
