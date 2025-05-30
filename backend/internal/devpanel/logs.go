package devpanel

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"
)

// LogManager handles service log collection and retrieval
type LogManager struct {
    logDir     string
    maxLines   int
    retention  time.Duration
    mu         sync.RWMutex
    logStreams map[string]*LogStream
}

// LogStream represents a streaming log connection
type LogStream struct {
    Service    string
    Lines      chan string
    Done       chan struct{}
    LastOffset int64
}

// NewLogManager creates a new log manager
func NewLogManager(logDir string, maxLines int, retention time.Duration) *LogManager {
    return &LogManager{
        logDir:     logDir,
        maxLines:   maxLines,
        retention:  retention,
        logStreams: make(map[string]*LogStream),
    }
}

// WriteLog writes a log entry for a service
func (lm *LogManager) WriteLog(service string, level string, message string) error {
    lm.mu.Lock()
    defer lm.mu.Unlock()

    logFile := filepath.Join(lm.logDir, service+".log")
    timestamp := time.Now().Format(time.RFC3339)
    logEntry := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, message)

    // Append to log file
    f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer f.Close()

    if _, err := f.WriteString(logEntry); err != nil {
        return err
    }

    // Broadcast to active streams
    if stream, exists := lm.logStreams[service]; exists {
        select {
        case stream.Lines <- logEntry:
        default:
            // Skip if channel is full
        }
    }

    return nil
}

// GetLogs retrieves logs for a service
func (lm *LogManager) GetLogs(service string, lines int) ([]string, error) {
    lm.mu.RLock()
    defer lm.mu.RUnlock()

    logFile := filepath.Join(lm.logDir, service+".log")
    if _, err := os.Stat(logFile); os.IsNotExist(err) {
        return nil, nil
    }

    file, err := os.Open(logFile)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var logLines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        logLines = append(logLines, scanner.Text())
    }

    // Trim to requested number of lines
    if len(logLines) > lines {
        logLines = logLines[len(logLines)-lines:]
    }

    return logLines, nil
}

// StreamLogs creates a new log stream for a service
func (lm *LogManager) StreamLogs(service string) *LogStream {
    lm.mu.Lock()
    defer lm.mu.Unlock()

    stream := &LogStream{
        Service:    service,
        Lines:      make(chan string, 1000),
        Done:       make(chan struct{}),
        LastOffset: 0,
    }

    lm.logStreams[service] = stream
    return stream
}

// StopStreaming stops a log stream
func (lm *LogManager) StopStreaming(service string) {
    lm.mu.Lock()
    defer lm.mu.Unlock()

    if stream, exists := lm.logStreams[service]; exists {
        close(stream.Done)
        delete(lm.logStreams, service)
    }
}

// CleanupOldLogs removes log files older than the retention period
func (lm *LogManager) CleanupOldLogs() error {
    lm.mu.Lock()
    defer lm.mu.Unlock()

    cutoff := time.Now().Add(-lm.retention)
    files, err := os.ReadDir(lm.logDir)
    if err != nil {
        return err
    }

    for _, file := range files {
        if !strings.HasSuffix(file.Name(), ".log") {
            continue
        }

        filePath := filepath.Join(lm.logDir, file.Name())
        info, err := file.Info()
        if err != nil {
            continue
        }

        if info.ModTime().Before(cutoff) {
            os.Remove(filePath)
        }
    }

    return nil
}

// RotateLogs rotates log files if they exceed the maximum size
func (lm *LogManager) RotateLogs() error {
    lm.mu.Lock()
    defer lm.mu.Unlock()

    files, err := os.ReadDir(lm.logDir)
    if err != nil {
        return err
    }

    for _, file := range files {
        if !strings.HasSuffix(file.Name(), ".log") {
            continue
        }

        filePath := filepath.Join(lm.logDir, file.Name())
        info, err := file.Info()
        if err != nil {
            continue
        }

        // Rotate if file is too large (e.g., > 100MB)
        if info.Size() > 100*1024*1024 {
            timestamp := time.Now().Format("20060102-150405")
            newPath := filePath + "." + timestamp
            os.Rename(filePath, newPath)
        }
    }

    return nil
}

// StartCleanup begins periodic log cleanup
func (lm *LogManager) StartCleanup() {
    go func() {
        ticker := time.NewTicker(24 * time.Hour)
        defer ticker.Stop()

        for range ticker.C {
            lm.CleanupOldLogs()
            lm.RotateLogs()
        }
    }()
}

// GetLogStats returns statistics about service logs
func (lm *LogManager) GetLogStats(service string) (map[string]interface{}, error) {
    lm.mu.RLock()
    defer lm.mu.RUnlock()

    logFile := filepath.Join(lm.logDir, service+".log")
    if _, err := os.Stat(logFile); os.IsNotExist(err) {
        return nil, nil
    }

    file, err := os.Open(logFile)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    stats := make(map[string]interface{})
    stats["file_size"] = 0
    stats["line_count"] = 0
    stats["error_count"] = 0
    stats["warning_count"] = 0

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        stats["line_count"] = stats["line_count"].(int) + 1
        if strings.Contains(line, "ERROR") {
            stats["error_count"] = stats["error_count"].(int) + 1
        }
        if strings.Contains(line, "WARN") {
            stats["warning_count"] = stats["warning_count"].(int) + 1
        }
    }

    if info, err := file.Stat(); err == nil {
        stats["file_size"] = info.Size()
    }

    return stats, nil
} 