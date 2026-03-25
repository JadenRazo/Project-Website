package core

import (
	"fmt"
	"sync"
	"time"
)

// ServiceStatus represents the current state of a service
type ServiceStatus struct {
	Running   bool
	Uptime    time.Duration
	Errors    []string
	LastCheck time.Time
}

// Service interface defines the methods that all services must implement
type Service interface {
	Start() error
	Stop() error
	Restart() error
	Status() ServiceStatus
	Name() string
	HealthCheck() error
}

// ServiceManager handles all services in the application
type ServiceManager struct {
	services            map[string]Service
	mu                  sync.RWMutex
	healthCheckInterval time.Duration
	stopHealth          chan struct{}
}

// NewServiceManager creates a new service manager instance
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services:            make(map[string]Service),
		healthCheckInterval: 30 * time.Second,
		stopHealth:          make(chan struct{}),
	}
}

// RegisterService adds a new service to the manager
func (sm *ServiceManager) RegisterService(service Service) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.services[service.Name()]; exists {
		return fmt.Errorf("service %s already registered", service.Name())
	}

	sm.services[service.Name()] = service
	return nil
}

// StartService starts a specific service
func (sm *ServiceManager) StartService(name string) error {
	sm.mu.RLock()
	service, exists := sm.services[name]
	sm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	return service.Start()
}

// StopService stops a specific service
func (sm *ServiceManager) StopService(name string) error {
	sm.mu.RLock()
	service, exists := sm.services[name]
	sm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	return service.Stop()
}

// RestartService restarts a specific service
func (sm *ServiceManager) RestartService(name string) error {
	sm.mu.RLock()
	service, exists := sm.services[name]
	sm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	return service.Restart()
}

// GetServiceStatus returns the status of a specific service
func (sm *ServiceManager) GetServiceStatus(name string) (ServiceStatus, error) {
	sm.mu.RLock()
	service, exists := sm.services[name]
	sm.mu.RUnlock()

	if !exists {
		return ServiceStatus{}, fmt.Errorf("service %s not found", name)
	}

	return service.Status(), nil
}

// GetAllServices returns a map of all registered services
func (sm *ServiceManager) GetAllServices() map[string]Service {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	services := make(map[string]Service, len(sm.services))
	for k, v := range sm.services {
		services[k] = v
	}
	return services
}

// StartHealthChecks begins monitoring all services
func (sm *ServiceManager) StartHealthChecks() {
	go func() {
		ticker := time.NewTicker(sm.healthCheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				sm.mu.RLock()
				for _, service := range sm.services {
					go func(s Service) {
						if err := s.HealthCheck(); err != nil {
							fmt.Printf("Health check failed for service %s: %v\n", s.Name(), err)
						}
					}(service)
				}
				sm.mu.RUnlock()
			case <-sm.stopHealth:
				return
			}
		}
	}()
}

// StopHealthChecks stops the health check goroutine
func (sm *ServiceManager) StopHealthChecks() {
	close(sm.stopHealth)
}

// StopAllServices gracefully stops all running services
func (sm *ServiceManager) StopAllServices() error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var errors []string
	for name, service := range sm.services {
		if err := service.Stop(); err != nil {
			errors = append(errors, fmt.Sprintf("failed to stop %s: %v", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors stopping services: %v", errors)
	}
	return nil
}
