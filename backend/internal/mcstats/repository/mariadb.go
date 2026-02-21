package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	_ "github.com/go-sql-driver/mysql"
)

// MariaDB wraps a MariaDB connection
type MariaDB struct {
	db *sql.DB
}

// NewMariaDB creates a new MariaDB connection
func NewMariaDB(cfg *config.MariaDBConfig) (*MariaDB, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("MariaDB is disabled in configuration")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&timeout=10s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open MariaDB connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeMinutes) * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping MariaDB: %w", err)
	}

	return &MariaDB{db: db}, nil
}

// DB returns the underlying sql.DB
func (m *MariaDB) DB() *sql.DB {
	return m.db
}

// Close closes the database connection
func (m *MariaDB) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

// Ping checks database connectivity
func (m *MariaDB) Ping(ctx context.Context) error {
	return m.db.PingContext(ctx)
}
