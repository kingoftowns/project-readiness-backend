// Package main is the entry point for the GitLab Readiness API.
// It demonstrates clean architecture and Go best practices for educational purposes.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/user/go-backend/internal/config"
	"github.com/user/go-backend/internal/database"
	"github.com/user/go-backend/internal/handlers"
	"github.com/user/go-backend/internal/repository"
	"github.com/user/go-backend/internal/router"
)

func main() {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist, we'll use environment variables
		if !os.IsNotExist(err) {
			slog.Warn("failed to load .env file", "error", err)
		}
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Set up structured logging
	logger := setupLogger(cfg.LogLevel)
	logger.Info("starting gitlab readiness api",
		"environment", cfg.Environment,
		"port", cfg.Port,
		"database_type", cfg.DatabaseType,
	)

	// Connect to database
	dbConfig := database.Config{
		Type:     cfg.DatabaseType,
		URL:      cfg.DatabaseURL,
		MaxConns: cfg.DBMaxConns,
		MaxIdle:  cfg.DBMaxIdle,
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	logger.Info("running database migrations")
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		// Auto-detect migrations path relative to project root
		if _, err := os.Stat("migrations"); err == nil {
			migrationsPath = "migrations"
		} else {
			migrationsPath = "../../migrations"
		}
	}
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		logger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Initialize repository
	projectRepo := repository.NewProjectRepository(db)

	// Initialize handlers
	projectHandler := handlers.NewProjectHandler(projectRepo, logger)

	// Set up routes
	handler := router.New(projectHandler, logger)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("server starting", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}

// setupLogger configures structured logging with slog
func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
		// Add source file information in debug mode
		AddSource: level == "debug",
	}

	// Use JSON format in production, text format in development
	var handler slog.Handler
	if os.Getenv("ENVIRONMENT") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}