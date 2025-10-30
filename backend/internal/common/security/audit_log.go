package security

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)

type AuditLogger struct {
	logger *log.Logger
	config AuditLogConfig
}

type AuditEntry struct {
	Timestamp  time.Time              `json:"timestamp"`
	ClientIP   string                 `json:"client_ip"`
	UserAgent  string                 `json:"user_agent"`
	Method     string                 `json:"method"`
	Path       string                 `json:"path"`
	StatusCode int                    `json:"status_code,omitempty"`
	UserID     string                 `json:"user_id,omitempty"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Duration   time.Duration          `json:"duration,omitempty"`
	Error      string                 `json:"error,omitempty"`
	SessionID  string                 `json:"session_id,omitempty"`
}

func NewAuditLogger(config AuditLogConfig) (*AuditLogger, error) {
	if !config.Enabled {
		return &AuditLogger{config: config}, nil
	}

	// Ensure log directory exists
	logDir := filepath.Dir(config.LogFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logger := log.New(&lumberjack.Logger{
		Filename:   config.LogFile,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}, "", 0)

	return &AuditLogger{
		logger: logger,
		config: config,
	}, nil
}

func (al *AuditLogger) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !al.config.Enabled {
			c.Next()
			return
		}

		start := time.Now()

		c.Next()

		al.LogRequest(c, start)
	}
}

func (al *AuditLogger) LogRequest(c *gin.Context, start time.Time) {
	if !al.config.Enabled || al.logger == nil {
		return
	}

	entry := AuditEntry{
		Timestamp:  time.Now(),
		ClientIP:   getClientIP(c),
		UserAgent:  c.GetHeader("User-Agent"),
		Method:     c.Request.Method,
		Path:       c.Request.URL.Path,
		StatusCode: c.Writer.Status(),
		Action:     "http_request",
		Duration:   time.Since(start),
	}

	// Add user info if available
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			entry.UserID = uid
		}
	}

	// Add session ID if available
	if sessionID, exists := c.Get("session_id"); exists {
		if sid, ok := sessionID.(string); ok {
			entry.SessionID = sid
		}
	}

	// Add error if request failed
	if c.Writer.Status() >= 400 {
		if err, exists := c.Get("error"); exists {
			if errStr, ok := err.(string); ok {
				entry.Error = errStr
			}
		}
	}

	al.logEntry(entry)
}

func (al *AuditLogger) LogAdminAction(c *gin.Context, action, resource string, details map[string]interface{}) {
	if !al.config.Enabled || al.logger == nil {
		return
	}

	entry := AuditEntry{
		Timestamp: time.Now(),
		ClientIP:  getClientIP(c),
		UserAgent: c.GetHeader("User-Agent"),
		Action:    action,
		Resource:  resource,
		Details:   details,
	}

	// Add user info if available
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			entry.UserID = uid
		}
	}

	// Add session ID if available
	if sessionID, exists := c.Get("session_id"); exists {
		if sid, ok := sessionID.(string); ok {
			entry.SessionID = sid
		}
	}

	al.logEntry(entry)
}

func (al *AuditLogger) LogAuthEvent(clientIP, action, userID string, success bool, details map[string]interface{}) {
	if !al.config.Enabled || al.logger == nil {
		return
	}

	entry := AuditEntry{
		Timestamp: time.Now(),
		ClientIP:  clientIP,
		Action:    action,
		UserID:    userID,
		Details:   details,
	}

	if !success {
		entry.Error = "Authentication failed"
	}

	al.logEntry(entry)
}

func (al *AuditLogger) LogSecurityEvent(clientIP, action string, details map[string]interface{}) {
	if !al.config.Enabled || al.logger == nil {
		return
	}

	entry := AuditEntry{
		Timestamp: time.Now(),
		ClientIP:  clientIP,
		Action:    action,
		Details:   details,
	}

	al.logEntry(entry)
}

func (al *AuditLogger) logEntry(entry AuditEntry) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal audit entry: %v", err)
		return
	}

	al.logger.Println(string(jsonData))
}

func (al *AuditLogger) IsEnabled() bool {
	return al.config.Enabled
}
