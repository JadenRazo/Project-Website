package logger

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/natefinch/lumberjack"
    "github.com/sirupsen/logrus"
    "github.com/JadenRazo/Project-Website/backend/internal/app/config"
)

var log *logrus.Logger

// InitLogger initializes the logger with the given configuration
func InitLogger(cfg *config.LoggingConfig) error {
    log = logrus.New()
    
    // Set log level
    level, err := logrus.ParseLevel(cfg.Level)
    if err != nil {
        return fmt.Errorf("invalid log level: %v", err)
    }
    log.SetLevel(level)

    // Configure formatter
    if cfg.Format == "json" {
        log.SetFormatter(&logrus.JSONFormatter{
            TimestampFormat: time.RFC3339,
        })
    } else {
        log.SetFormatter(&logrus.TextFormatter{
            TimestampFormat: time.RFC3339,
            FullTimestamp: true,
        })
    }

    // Configure output
    if cfg.Output == "file" {
        // Ensure directory exists
        dir := filepath.Dir(cfg.FilePath)
        if err := os.MkdirAll(dir, 0755); err != nil {
            return fmt.Errorf("failed to create log directory: %v", err)
        }

        // Configure log rotation
        log.SetOutput(&lumberjack.Logger{
            Filename:   cfg.FilePath,
            MaxSize:    cfg.MaxSize,    // MB
            MaxBackups: cfg.MaxBackups,
            MaxAge:     cfg.MaxAge,     // days
            Compress:   cfg.Compress,
        })
    }

    return nil
}

// WithFields creates a new entry with fields
func WithFields(fields map[string]interface{}) *logrus.Entry {
    return log.WithFields(fields)
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

// ApplyFilters applies the configured log filters
func ApplyFilters(entry *logrus.Entry, filters []config.LogFilter) bool {
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
