package websocket

import (
	"sync"
	"time"
)

// PresenceManager manages user presence status
type PresenceManager struct {
	mu       sync.RWMutex
	presence map[string]*UserPresence
}

// UserPresence represents a user's presence information
type UserPresence struct {
	UserID       string    `json:"user_id"`
	Status       string    `json:"status"` // online, away, offline
	LastActivity time.Time `json:"last_activity"`
}

// NewPresenceManager creates a new presence manager
func NewPresenceManager() *PresenceManager {
	return &PresenceManager{
		presence: make(map[string]*UserPresence),
	}
}

// UpdatePresence updates a user's presence status
func (pm *PresenceManager) UpdatePresence(userID string, status string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.presence[userID] = &UserPresence{
		UserID:       userID,
		Status:       status,
		LastActivity: time.Now(),
	}
}

// GetPresence returns a user's presence status
func (pm *PresenceManager) GetPresence(userID string) (*UserPresence, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	presence, exists := pm.presence[userID]
	return presence, exists
}

// RemovePresence removes a user from presence tracking
func (pm *PresenceManager) RemovePresence(userID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	delete(pm.presence, userID)
}

// GetAllPresence returns all user presence information
func (pm *PresenceManager) GetAllPresence() map[string]*UserPresence {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make(map[string]*UserPresence)
	for k, v := range pm.presence {
		result[k] = v
	}
	return result
}