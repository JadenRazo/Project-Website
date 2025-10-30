package core

import (
	"sync"
	"sync/atomic"
	"time"
)

// BaseService provides common functionality for all services
type BaseService struct {
	name         string
	running      bool
	startTime    time.Time
	errors       []string
	requestCount int64
	errorCount   int64
	mu           sync.RWMutex
}

// NewBaseService creates a new base service instance
func NewBaseService(name string) *BaseService {
	return &BaseService{
		name:   name,
		errors: make([]string, 0),
	}
}

// Name returns the service name
func (s *BaseService) Name() string {
	return s.name
}

// Start marks the service as running
func (s *BaseService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return nil
	}

	s.running = true
	s.startTime = time.Now()
	s.errors = make([]string, 0)
	return nil
}

// Stop marks the service as stopped
func (s *BaseService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.running = false
	return nil
}

// Restart restarts the service
func (s *BaseService) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}
	return s.Start()
}

// Status returns the current service status
func (s *BaseService) Status() ServiceStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var uptime time.Duration
	if s.running {
		uptime = time.Since(s.startTime)
	}

	return ServiceStatus{
		Running:   s.running,
		Uptime:    uptime,
		Errors:    s.errors,
		LastCheck: time.Now(),
	}
}

// HealthCheck performs a basic health check
func (s *BaseService) HealthCheck() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.running {
		return nil
	}

	return nil
}

// AddError adds an error to the service's error list
func (s *BaseService) AddError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err != nil {
		s.errors = append(s.errors, err.Error())
	}
}

// ClearErrors clears the service's error list
func (s *BaseService) ClearErrors() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.errors = make([]string, 0)
}

// IsRunning returns whether the service is currently running
func (s *BaseService) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.running
}

// GetStartTime returns when the service was started
func (s *BaseService) GetStartTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.startTime
}

// GetErrors returns the current list of errors
func (s *BaseService) GetErrors() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	errors := make([]string, len(s.errors))
	copy(errors, s.errors)
	return errors
}

// IncrementRequests increments the request counter
func (s *BaseService) IncrementRequests() {
	atomic.AddInt64(&s.requestCount, 1)
}

// IncrementErrors increments the error counter
func (s *BaseService) IncrementErrors() {
	atomic.AddInt64(&s.errorCount, 1)
}

// GetRequestCount returns the total number of requests processed
func (s *BaseService) GetRequestCount() int64 {
	return atomic.LoadInt64(&s.requestCount)
}

// GetErrorCount returns the total number of errors encountered
func (s *BaseService) GetErrorCount() int64 {
	return atomic.LoadInt64(&s.errorCount)
}
