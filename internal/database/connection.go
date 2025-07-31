// Package database handles database connections and migrations.
// It supports both SQLite (for development/testing) and PostgreSQL (for production).
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	// Database drivers
	_ "github.com/lib/pq"          // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Config holds database configuration
type Config struct {
	Type     string // "sqlite" or "postgres"
	URL      string // Connection string
	MaxConns int    // Maximum number of open connections
	MaxIdle  int    // Maximum number of idle connections
}

// DB wraps the standard sql.DB with additional functionality
type DB struct {
	*sql.DB
	dbType string
}

// NewConnection creates a new database connection based on the configuration.
// It supports both SQLite and PostgreSQL databases.
func NewConnection(cfg Config) (*DB, error) {
	var driverName string
	
	switch cfg.Type {
	case "sqlite":
		driverName = "sqlite3"
		// Add SQLite specific pragmas for better performance
		if cfg.URL != ":memory:" {
			cfg.URL += "?_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=ON"
		}
	case "postgres":
		driverName = "postgres"
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	// Open database connection
	db, err := sql.Open(driverName, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	if cfg.MaxConns > 0 {
		db.SetMaxOpenConns(cfg.MaxConns)
	} else {
		db.SetMaxOpenConns(25) // Default
	}
	
	if cfg.MaxIdle > 0 {
		db.SetMaxIdleConns(cfg.MaxIdle)
	} else {
		db.SetMaxIdleConns(5) // Default
	}
	
	// Set connection lifetime
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("database connection established", 
		"type", cfg.Type,
		"max_conns", db.Stats().MaxOpenConnections,
	)

	return &DB{
		DB:     db,
		dbType: cfg.Type,
	}, nil
}

// Type returns the database type (sqlite or postgres)
func (db *DB) Type() string {
	return db.dbType
}

// Close closes the database connection
func (db *DB) Close() error {
	slog.Info("closing database connection")
	return db.DB.Close()
}

// HealthCheck performs a simple health check on the database
func (db *DB) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Run a simple query to ensure the database is responsive
	var result int
	query := "SELECT 1"
	if err := db.QueryRowContext(ctx, query).Scan(&result); err != nil {
		return fmt.Errorf("database query check failed: %w", err)
	}

	return nil
}

