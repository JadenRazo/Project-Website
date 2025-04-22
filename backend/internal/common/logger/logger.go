package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	globalconfig "github.com/JadenRazo/Project-Website/backend/config"
	appconfig "github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/google/uuid"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// Key constants for context values
const (
	RequestIDKey  = "request_id"
	UserIDKey     = "user_id"
	CorrelationID = "correlation_id"
	SessionID     = "session_id"
	ServiceKey    = "service"
	VersionKey    = "version"
)

var log *logrus.Logger
var serviceName string
var serviceVersion string

// Logger defines the logging interface
type Logger interface {
	WithContext(ctx context.Context) *Entry
	WithField(key string, value interface{}) *Entry
	WithFields(fields map[string]interface{}) *Entry
	WithError(err error) *Entry
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// Entry is a wrapper around logrus.Entry with added functionality
type Entry struct {
	*logrus.Entry
}

// InitLogger initializes the logger with the given configuration
func InitLogger(cfg *appconfig.LoggingConfig, appName, version string) error {
	log = logrus.New()
	serviceName = appName
	serviceVersion = version

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %v", err)
	}
	log.SetLevel(level)

	// Configure formatter
	if cfg.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: cfg.TimeFormat,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := filepath.Base(f.File)
				return "", fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: cfg.TimeFormat,
			FullTimestamp:   true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := filepath.Base(f.File)
				return "", fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
	}

	// Enable caller info
	log.SetReportCaller(true)

	// Configure output
	if cfg.Output == "file" {
		// Ensure directory exists
		dir := filepath.Dir(cfg.Filename)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}

		// Configure log rotation
		log.SetOutput(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize, // MB
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge, // days
			Compress:   cfg.Compress,
		})
	} else if cfg.Output == "stdout" {
		log.SetOutput(os.Stdout)
	} else if cfg.Output == "stderr" {
		log.SetOutput(os.Stderr)
	} else {
		// Default to stdout
		log.SetOutput(os.Stdout)
	}

	Info("Logger initialized successfully", "level", level, "output", cfg.Output)
	return nil
}

// WithContext extracts context values and creates an entry with those fields
func WithContext(ctx context.Context) *Entry {
	entry := log.WithFields(logrus.Fields{
		ServiceKey:  serviceName,
		VersionKey:  serviceVersion,
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
	})

	// Extract common context values if they exist
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		entry = entry.WithField(RequestIDKey, reqID)
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		entry = entry.WithField(UserIDKey, userID)
	}

	if corrID, ok := ctx.Value(CorrelationID).(string); ok {
		entry = entry.WithField(CorrelationID, corrID)
	} else {
		// Generate correlation ID if not present
		entry = entry.WithField(CorrelationID, uuid.New().String())
	}

	if sessID, ok := ctx.Value(SessionID).(string); ok {
		entry = entry.WithField(SessionID, sessID)
	}

	return &Entry{entry}
}

// WithField adds a field to the log entry
func WithField(key string, value interface{}) *Entry {
	return &Entry{log.WithField(key, value)}
}

// WithFields adds multiple fields to the log entry
func WithFields(fields map[string]interface{}) *Entry {
	baseFields := logrus.Fields{
		ServiceKey:  serviceName,
		VersionKey:  serviceVersion,
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
	}

	// Merge with provided fields
	for k, v := range fields {
		baseFields[k] = v
	}

	return &Entry{log.WithFields(baseFields)}
}

// WithError adds an error to the log entry
func WithError(err error) *Entry {
	return &Entry{log.WithError(err)}
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	log.Error(args...)
}

// Fatal logs a fatal message and exits
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Debug logs a debug message with the given entry
func (e *Entry) Debug(args ...interface{}) {
	e.Entry.Debug(args...)
}

// Info logs an info message with the given entry
func (e *Entry) Info(args ...interface{}) {
	e.Entry.Info(args...)
}

// Warn logs a warning message with the given entry
func (e *Entry) Warn(args ...interface{}) {
	e.Entry.Warn(args...)
}

// Error logs an error message with the given entry
func (e *Entry) Error(args ...interface{}) {
	e.Entry.Error(args...)
}

// Fatal logs a fatal message and exits with the given entry
func (e *Entry) Fatal(args ...interface{}) {
	e.Entry.Fatal(args...)
}

// ApplyFilters applies the configured log filters
func ApplyFilters(entry *logrus.Entry, filters []globalconfig.LogFilter) bool {
	for _, filter := range filters {
		value, exists := entry.Data[filter.Field]
		if !exists {
			return false
		}

		strValue := fmt.Sprintf("%v", value)
		switch filter.Operator {
		case "eq":
			if strValue != filter.Value {
				return false
			}
		case "ne":
			if strValue == filter.Value {
				return false
			}
		case "contains":
			if !strings.Contains(strValue, filter.Value) {
				return false
			}
		case "gt":
			if strValue <= filter.Value {
				return false
			}
		case "lt":
			if strValue >= filter.Value {
				return false
			}
		}
	}
	return true
}

// Shutdown performs any cleanup needed for the logger
func Shutdown() error {
	Info("Logger shutting down")
	return nil
}
