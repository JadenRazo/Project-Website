package circuitbreaker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// State represents the circuit breaker state
type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

// String returns the string representation of the state
func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateHalfOpen:
		return "HALF_OPEN"
	case StateOpen:
		return "OPEN"
	default:
		return "UNKNOWN"
	}
}

// Config holds the circuit breaker configuration
type Config struct {
	// MaxFailures is the maximum number of failures before opening the circuit
	MaxFailures int
	// Timeout is how long to wait before transitioning from Open to Half-Open
	Timeout time.Duration
	// MaxRequests is the maximum number of requests allowed in Half-Open state
	MaxRequests int
	// SuccessThreshold is the number of consecutive successes needed to close the circuit
	SuccessThreshold int
}

// DefaultConfig returns a default circuit breaker configuration
func DefaultConfig() Config {
	return Config{
		MaxFailures:      5,
		Timeout:          30 * time.Second,
		MaxRequests:      3,
		SuccessThreshold: 2,
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config           Config
	mutex            sync.RWMutex
	state            State
	failures         int
	successes        int
	lastFailureTime  time.Time
	nextAttemptTime  time.Time
	halfOpenRequests int
}

// New creates a new circuit breaker with the given configuration
func New(config Config) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// NewDefault creates a new circuit breaker with default configuration
func NewDefault() *CircuitBreaker {
	return New(DefaultConfig())
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	// Check if circuit breaker allows the request
	if !cb.allowRequest() {
		return ErrCircuitBreakerOpen
	}

	// Execute the function
	err := fn()

	// Record the result
	cb.recordResult(err)

	return err
}

// allowRequest determines if a request should be allowed through
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		// Check if enough time has passed to transition to half-open
		if time.Now().After(cb.nextAttemptTime) {
			cb.state = StateHalfOpen
			cb.halfOpenRequests = 0
			cb.successes = 0
			return true
		}
		return false
	case StateHalfOpen:
		// Allow limited requests in half-open state
		if cb.halfOpenRequests < cb.config.MaxRequests {
			cb.halfOpenRequests++
			return true
		}
		return false
	default:
		return false
	}
}

// recordResult records the result of a request
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err != nil {
		cb.recordFailure()
	} else {
		cb.recordSuccess()
	}
}

// recordFailure records a failure
func (cb *CircuitBreaker) recordFailure() {
	cb.failures++
	cb.successes = 0
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.config.MaxFailures {
			cb.state = StateOpen
			cb.nextAttemptTime = time.Now().Add(cb.config.Timeout)
		}
	case StateHalfOpen:
		// Any failure in half-open state immediately opens the circuit
		cb.state = StateOpen
		cb.nextAttemptTime = time.Now().Add(cb.config.Timeout)
		cb.halfOpenRequests = 0
	}
}

// recordSuccess records a success
func (cb *CircuitBreaker) recordSuccess() {
	cb.successes++

	switch cb.state {
	case StateHalfOpen:
		if cb.successes >= cb.config.SuccessThreshold {
			cb.state = StateClosed
			cb.failures = 0
			cb.successes = 0
			cb.halfOpenRequests = 0
		}
	case StateClosed:
		// Reset failure count on success
		cb.failures = 0
	}
}

// GetState returns the current circuit breaker state
func (cb *CircuitBreaker) GetState() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// GetStats returns current statistics
func (cb *CircuitBreaker) GetStats() Stats {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	return Stats{
		State:            cb.state,
		Failures:         cb.failures,
		Successes:        cb.successes,
		LastFailureTime:  cb.lastFailureTime,
		NextAttemptTime:  cb.nextAttemptTime,
		HalfOpenRequests: cb.halfOpenRequests,
	}
}

// Stats holds circuit breaker statistics
type Stats struct {
	State            State
	Failures         int
	Successes        int
	LastFailureTime  time.Time
	NextAttemptTime  time.Time
	HalfOpenRequests int
}

// String returns a string representation of the stats
func (s Stats) String() string {
	return fmt.Sprintf("State: %s, Failures: %d, Successes: %d, HalfOpenRequests: %d",
		s.State, s.Failures, s.Successes, s.HalfOpenRequests)
}

// Reset resets the circuit breaker to its initial state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.state = StateClosed
	cb.failures = 0
	cb.successes = 0
	cb.halfOpenRequests = 0
	cb.lastFailureTime = time.Time{}
	cb.nextAttemptTime = time.Time{}
}

// ForceOpen forces the circuit breaker into the open state
func (cb *CircuitBreaker) ForceOpen() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.state = StateOpen
	cb.nextAttemptTime = time.Now().Add(cb.config.Timeout)
}

// ForceClose forces the circuit breaker into the closed state
func (cb *CircuitBreaker) ForceClose() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.state = StateClosed
	cb.failures = 0
	cb.successes = 0
	cb.halfOpenRequests = 0
}

// Predefined errors
var (
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
	ErrTooManyRequests    = errors.New("too many requests in half-open state")
)

// DatabaseCircuitBreaker is a specialized circuit breaker for database operations
type DatabaseCircuitBreaker struct {
	*CircuitBreaker
	name string
}

// NewDatabaseCircuitBreaker creates a new database-specific circuit breaker
func NewDatabaseCircuitBreaker(name string, config Config) *DatabaseCircuitBreaker {
	return &DatabaseCircuitBreaker{
		CircuitBreaker: New(config),
		name:           name,
	}
}

// ExecuteDB executes a database operation with circuit breaker protection
func (dcb *DatabaseCircuitBreaker) ExecuteDB(ctx context.Context, operation func() error) error {
	return dcb.Execute(ctx, func() error {
		// Add database-specific error handling
		err := operation()
		if err != nil {
			// Log database errors for monitoring
			fmt.Printf("Database circuit breaker '%s' recorded failure: %v\n", dcb.name, err)
		}
		return err
	})
}

// IsHealthy returns true if the circuit breaker is in a healthy state
func (dcb *DatabaseCircuitBreaker) IsHealthy() bool {
	return dcb.GetState() != StateOpen
}

// GetHealthStatus returns a health status string
func (dcb *DatabaseCircuitBreaker) GetHealthStatus() string {
	stats := dcb.GetStats()
	switch stats.State {
	case StateClosed:
		return "healthy"
	case StateHalfOpen:
		return "recovering"
	case StateOpen:
		return "unhealthy"
	default:
		return "unknown"
	}
}
