package status

import "fmt"

// ErrorCode represents a status service error code
type ErrorCode string

const (
	// Database errors
	ErrCodeDatabaseConnection  ErrorCode = "DB_CONNECTION_FAILED"
	ErrCodeDatabaseTimeout     ErrorCode = "DB_TIMEOUT"
	ErrCodeDatabaseUnavailable ErrorCode = "DB_UNAVAILABLE"

	// Service errors
	ErrCodeServiceTimeout     ErrorCode = "SERVICE_TIMEOUT"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeServiceDegraded    ErrorCode = "SERVICE_DEGRADED"

	// General errors
	ErrCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrCodeConfigurationError ErrorCode = "CONFIGURATION_ERROR"
)

// ErrorInfo contains sanitized error information for public consumption
type ErrorInfo struct {
	Code        ErrorCode `json:"code"`
	Message     string    `json:"message"`
	Recoverable bool      `json:"recoverable"`
}

// errorCodeMap maps error codes to user-friendly messages
var errorCodeMap = map[ErrorCode]ErrorInfo{
	ErrCodeDatabaseConnection: {
		Code:        ErrCodeDatabaseConnection,
		Message:     "Database connection temporarily unavailable",
		Recoverable: true,
	},
	ErrCodeDatabaseTimeout: {
		Code:        ErrCodeDatabaseTimeout,
		Message:     "Database response timeout",
		Recoverable: true,
	},
	ErrCodeDatabaseUnavailable: {
		Code:        ErrCodeDatabaseUnavailable,
		Message:     "Database service unavailable",
		Recoverable: true,
	},
	ErrCodeServiceTimeout: {
		Code:        ErrCodeServiceTimeout,
		Message:     "Service response timeout",
		Recoverable: true,
	},
	ErrCodeServiceUnavailable: {
		Code:        ErrCodeServiceUnavailable,
		Message:     "Service temporarily unavailable",
		Recoverable: true,
	},
	ErrCodeServiceDegraded: {
		Code:        ErrCodeServiceDegraded,
		Message:     "Service performance degraded",
		Recoverable: true,
	},
	ErrCodeInternalError: {
		Code:        ErrCodeInternalError,
		Message:     "Internal service error",
		Recoverable: false,
	},
	ErrCodeConfigurationError: {
		Code:        ErrCodeConfigurationError,
		Message:     "Service configuration error",
		Recoverable: false,
	},
}

// SanitizeError converts raw errors into sanitized error information
func SanitizeError(err error) ErrorInfo {
	if err == nil {
		return ErrorInfo{}
	}

	errStr := err.Error()

	// Check for database-related errors
	if containsAny(errStr, []string{"connection refused", "connect:", "dial tcp", "no such host"}) {
		return errorCodeMap[ErrCodeDatabaseConnection]
	}

	if containsAny(errStr, []string{"timeout", "deadline exceeded", "context deadline exceeded"}) {
		return errorCodeMap[ErrCodeDatabaseTimeout]
	}

	if containsAny(errStr, []string{"database", "sql", "postgres", "sqlite", "gorm"}) {
		return errorCodeMap[ErrCodeDatabaseUnavailable]
	}

	// Default to internal error for unknown errors
	return errorCodeMap[ErrCodeInternalError]
}

// SanitizeErrorMessage returns a safe error message for public consumption
func SanitizeErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	sanitized := SanitizeError(err)
	return sanitized.Message
}

// containsAny checks if a string contains any of the specified substrings
func containsAny(s string, substrings []string) bool {
	for _, substring := range substrings {
		if fmt.Sprintf("%s", s) != s {
			continue
		}
		if len(s) >= len(substring) {
			for i := 0; i <= len(s)-len(substring); i++ {
				if s[i:i+len(substring)] == substring {
					return true
				}
			}
		}
	}
	return false
}
