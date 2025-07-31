// Package config handles application configuration.
// It follows the 12-factor app principle of storing config in the environment.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Port string
	Host string

	// Database configuration
	DatabaseType string // "sqlite" or "postgres"
	DatabaseURL  string

	// Database connection pool settings
	DBMaxConns int
	DBMaxIdle  int

	// Logging
	LogLevel string

	// Environment
	Environment string // "development", "production", etc.
}

// Load reads configuration from environment variables.
// It provides sensible defaults for development.
func Load() (*Config, error) {
	cfg := &Config{
		// Server defaults
		Port: getEnv("PORT", "8080"),
		Host: getEnv("HOST", "0.0.0.0"),

		// Database defaults
		DatabaseType: getEnv("DATABASE_TYPE", "sqlite"),
		DatabaseURL:  getEnv("DATABASE_URL", "gitlab_readiness.db"),

		// Connection pool defaults
		DBMaxConns: getEnvAsInt("DB_MAX_CONNS", 25),
		DBMaxIdle:  getEnvAsInt("DB_MAX_IDLE", 5),

		// Logging defaults
		LogLevel: getEnv("LOG_LEVEL", "info"),

		// Environment
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate database type
	if c.DatabaseType != "sqlite" && c.DatabaseType != "postgres" {
		return fmt.Errorf("invalid DATABASE_TYPE: must be 'sqlite' or 'postgres'")
	}

	// Validate database URL
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	// Validate port
	port, err := strconv.Atoi(c.Port)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid PORT: must be between 1 and 65535")
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid LOG_LEVEL: must be debug, info, warn, or error")
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// Helper functions

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer with a fallback default
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}