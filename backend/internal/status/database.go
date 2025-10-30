package status

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/JadenRazo/Project-Website/backend/internal/common/circuitbreaker"
	"github.com/JadenRazo/Project-Website/backend/internal/common/database"
	"gorm.io/gorm"
)

// DatabaseStatus represents the current status of the database connection
type DatabaseStatus int

const (
	DatabaseStatusOnline DatabaseStatus = iota
	DatabaseStatusDegraded
	DatabaseStatusOffline
	DatabaseStatusRecovering
)

// StatusDatabase wraps the database connection with health monitoring and retry logic
type StatusDatabase struct {
	db             database.Database
	originalDB     *gorm.DB
	status         DatabaseStatus
	lastError      error
	lastPing       time.Time
	mu             sync.RWMutex
	config         *config.DatabaseConfig
	retryAttempts  int
	maxRetries     int
	backoffBase    time.Duration
	circuitBreaker *circuitbreaker.DatabaseCircuitBreaker
}

// NewStatusDatabase creates a new database wrapper with health monitoring
func NewStatusDatabase(cfg *config.DatabaseConfig, originalDB *gorm.DB) *StatusDatabase {
	// Configure circuit breaker for database operations
	cbConfig := circuitbreaker.Config{
		MaxFailures:      3,                // Open circuit after 3 failures
		Timeout:          30 * time.Second, // Wait 30 seconds before trying again
		MaxRequests:      2,                // Allow 2 requests in half-open state
		SuccessThreshold: 2,                // Need 2 successes to close circuit
	}

	sdb := &StatusDatabase{
		originalDB:     originalDB,
		status:         DatabaseStatusOnline,
		config:         cfg,
		maxRetries:     5,
		backoffBase:    time.Second,
		circuitBreaker: circuitbreaker.NewDatabaseCircuitBreaker("status-db", cbConfig),
	}

	// Try to create a proper database interface
	if cfg != nil {
		if db, err := database.NewDatabase(cfg); err == nil {
			sdb.db = db
		}
	}

	// Start health monitoring
	go sdb.startHealthMonitoring()

	return sdb
}

// GetDB returns the GORM database instance
func (sdb *StatusDatabase) GetDB() *gorm.DB {
	return sdb.originalDB
}

// GetStatus returns the current database status
func (sdb *StatusDatabase) GetStatus() DatabaseStatus {
	sdb.mu.RLock()
	defer sdb.mu.RUnlock()
	return sdb.status
}

// GetLastError returns the last database error
func (sdb *StatusDatabase) GetLastError() error {
	sdb.mu.RLock()
	defer sdb.mu.RUnlock()
	return sdb.lastError
}

// IsHealthy returns true if the database is online
func (sdb *StatusDatabase) IsHealthy() bool {
	return sdb.GetStatus() == DatabaseStatusOnline
}

// Ping performs a database health check with retry logic and circuit breaker protection
func (sdb *StatusDatabase) Ping(ctx context.Context) error {
	// Use circuit breaker to protect ping operations
	err := sdb.circuitBreaker.ExecuteDB(ctx, func() error {
		return sdb.performPing(ctx)
	})

	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	if err == nil {
		// Success - reset retry attempts and update status
		sdb.retryAttempts = 0
		sdb.status = DatabaseStatusOnline
		sdb.lastError = nil
		sdb.lastPing = time.Now()
		return nil
	}

	// Handle failure
	sdb.lastError = err
	sdb.retryAttempts++

	// Determine new status based on circuit breaker state and retry attempts
	if !sdb.circuitBreaker.IsHealthy() || sdb.retryAttempts >= sdb.maxRetries {
		sdb.status = DatabaseStatusOffline
	} else {
		sdb.status = DatabaseStatusDegraded
	}

	return err
}

// performPing executes the actual database ping
func (sdb *StatusDatabase) performPing(ctx context.Context) error {
	// Try the proper database interface first
	if sdb.db != nil {
		return sdb.db.Ping(ctx)
	}

	// Fallback to direct GORM ping
	if sdb.originalDB != nil {
		sqlDB, err := sdb.originalDB.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql.DB: %w", err)
		}
		return sqlDB.PingContext(ctx)
	}

	return fmt.Errorf("no database connection available")
}

// ExecuteWithRetry executes a database operation with retry logic and circuit breaker protection
func (sdb *StatusDatabase) ExecuteWithRetry(ctx context.Context, operation func(*gorm.DB) error) error {
	// Use circuit breaker to protect database operations
	return sdb.circuitBreaker.ExecuteDB(ctx, func() error {
		// Check if database is available
		if !sdb.IsHealthy() {
			// Try to recover connection first
			if err := sdb.attemptRecovery(ctx); err != nil {
				return fmt.Errorf("database unavailable: %w", err)
			}
		}

		// Execute the operation
		err := operation(sdb.originalDB)
		if err != nil {
			// If operation fails, mark as degraded
			sdb.mu.Lock()
			sdb.status = DatabaseStatusDegraded
			sdb.lastError = err
			sdb.mu.Unlock()
		}

		return err
	})
}

// attemptRecovery tries to recover the database connection
func (sdb *StatusDatabase) attemptRecovery(ctx context.Context) error {
	sdb.mu.Lock()
	sdb.status = DatabaseStatusRecovering
	sdb.mu.Unlock()

	// Calculate backoff time
	backoffTime := sdb.backoffBase * time.Duration(sdb.retryAttempts*sdb.retryAttempts)
	if backoffTime > 30*time.Second {
		backoffTime = 30 * time.Second
	}

	// Wait before attempting recovery
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(backoffTime):
	}

	// Try to ping the database
	return sdb.Ping(ctx)
}

// startHealthMonitoring starts background health monitoring
func (sdb *StatusDatabase) startHealthMonitoring() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// Perform health check
		sdb.Ping(ctx)

		// If database is offline, try recovery
		if sdb.GetStatus() == DatabaseStatusOffline {
			sdb.attemptRecovery(ctx)
		}

		cancel()
	}
}

// GetStatusString returns a human-readable status string
func (sdb *StatusDatabase) GetStatusString() string {
	cbStatus := sdb.circuitBreaker.GetHealthStatus()
	dbStatus := sdb.GetStatus()

	switch dbStatus {
	case DatabaseStatusOnline:
		if cbStatus == "healthy" {
			return "online"
		}
		return "online (circuit breaker: " + cbStatus + ")"
	case DatabaseStatusDegraded:
		return "degraded (circuit breaker: " + cbStatus + ")"
	case DatabaseStatusOffline:
		return "offline (circuit breaker: " + cbStatus + ")"
	case DatabaseStatusRecovering:
		return "recovering (circuit breaker: " + cbStatus + ")"
	default:
		return "unknown"
	}
}

// GetCircuitBreakerStats returns circuit breaker statistics
func (sdb *StatusDatabase) GetCircuitBreakerStats() circuitbreaker.Stats {
	return sdb.circuitBreaker.GetStats()
}

// IsCircuitBreakerHealthy returns true if the circuit breaker is in a healthy state
func (sdb *StatusDatabase) IsCircuitBreakerHealthy() bool {
	return sdb.circuitBreaker.IsHealthy()
}
