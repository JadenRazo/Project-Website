package entity

import (
	"time"

	"github.com/google/uuid"
)

// AuditAction defines the type of action performed
type AuditAction string

const (
	// AuditActionCreate represents a creation action
	AuditActionCreate AuditAction = "create"

	// AuditActionUpdate represents an update action
	AuditActionUpdate AuditAction = "update"

	// AuditActionDelete represents a deletion action
	AuditActionDelete AuditAction = "delete"

	// AuditActionLogin represents a login action
	AuditActionLogin AuditAction = "login"

	// AuditActionLogout represents a logout action
	AuditActionLogout AuditAction = "logout"

	// AuditActionFailedLogin represents a failed login attempt
	AuditActionFailedLogin AuditAction = "failed_login"
)

// AuditEntityType defines the type of entity being audited
type AuditEntityType string

const (
	// AuditEntityUser represents user-related actions
	AuditEntityUser AuditEntityType = "user"

	// AuditEntityURL represents URL-related actions
	AuditEntityURL AuditEntityType = "url"

	// AuditEntityMessage represents message-related actions
	AuditEntityMessage AuditEntityType = "message"

	// AuditEntityChannel represents channel-related actions
	AuditEntityChannel AuditEntityType = "channel"

	// AuditEntitySystem represents system-level actions
	AuditEntitySystem AuditEntityType = "system"
)

// AuditLog represents a record of system activity
type AuditLog struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	UserID       *uuid.UUID      `json:"user_id,omitempty" db:"user_id"`
	Action       AuditAction     `json:"action" db:"action"`
	EntityType   AuditEntityType `json:"entity_type" db:"entity_type"`
	EntityID     string          `json:"entity_id" db:"entity_id"`
	IPAddress    string          `json:"ip_address" db:"ip_address"`
	UserAgent    string          `json:"user_agent" db:"user_agent"`
	OldValue     interface{}     `json:"old_value,omitempty" db:"old_value"`
	NewValue     interface{}     `json:"new_value,omitempty" db:"new_value"`
	ErrorMessage string          `json:"error_message,omitempty" db:"error_message"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(
	userID *uuid.UUID,
	action AuditAction,
	entityType AuditEntityType,
	entityID string,
	ipAddress string,
	userAgent string,
) *AuditLog {
	return &AuditLog{
		ID:         uuid.New(),
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		CreatedAt:  time.Now(),
	}
}

// WithValues adds old and new values to the audit log
func (a *AuditLog) WithValues(oldValue, newValue interface{}) *AuditLog {
	a.OldValue = oldValue
	a.NewValue = newValue
	return a
}

// WithError adds an error message to the audit log
func (a *AuditLog) WithError(errorMessage string) *AuditLog {
	a.ErrorMessage = errorMessage
	return a
}

// IsSuccess checks if the audit log represents a successful operation
func (a *AuditLog) IsSuccess() bool {
	return a.ErrorMessage == ""
}

// IsUserAction checks if the audit log is related to a user action
func (a *AuditLog) IsUserAction() bool {
	return a.UserID != nil
}

// IsSystemAction checks if the audit log is related to a system action
func (a *AuditLog) IsSystemAction() bool {
	return a.EntityType == AuditEntitySystem
}

// FormatTimestamp returns the created timestamp in a specific format
func (a *AuditLog) FormatTimestamp(layout string) string {
	return a.CreatedAt.Format(layout)
}
