package performance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/JadenRazo/Project-Website/backend/internal/messaging"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/entity"
	urlshortener "github.com/JadenRazo/Project-Website/backend/internal/urlshortener"
	urlEntity "github.com/JadenRazo/Project-Website/backend/internal/urlshortener/entity"
)

// BenchmarkResult stores performance test results
type BenchmarkResult struct {
	TestName        string        `json:"test_name"`
	Duration        time.Duration `json:"duration"`
	RequestsPerSec  float64       `json:"requests_per_sec"`
	MemoryUsage     MemoryStats   `json:"memory_usage"`
	DatabaseQueries int           `json:"database_queries"`
	Errors          int           `json:"errors"`
	P95Latency      time.Duration `json:"p95_latency"`
	P99Latency      time.Duration `json:"p99_latency"`
}

type MemoryStats struct {
	AllocMB      float64 `json:"alloc_mb"`
	TotalAllocMB float64 `json:"total_alloc_mb"`
	SysMB        float64 `json:"sys_mb"`
	NumGC        uint32  `json:"num_gc"`
}

// Helper function to capture memory stats
func getMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryStats{
		AllocMB:      float64(m.Alloc) / 1024 / 1024,
		TotalAllocMB: float64(m.TotalAlloc) / 1024 / 1024,
		SysMB:        float64(m.Sys) / 1024 / 1024,
		NumGC:        m.NumGC,
	}
}

// Setup test database
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Reduce log noise during benchmarks
	})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Auto migrate
	db.AutoMigrate(
		&entity.Channel{},
		&entity.ChannelMember{},
		&entity.Message{},
		&entity.MessageReaction{},
		&urlEntity.ShortenedURL{},
		&urlEntity.URLClick{},
	)

	return db
}

// Setup test router with services
func setupTestRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB()
	router := gin.New()

	// Add middleware to track database queries
	queryCount := 0
	db.Callback().Query().Before("gorm:query").Register("count_queries", func(db *gorm.DB) {
		queryCount++
	})

	// Setup services
	messagingService := messaging.NewService(db, messaging.Config{
		WebSocketPort:    8081,
		MaxMessageSize:   1024,
		MaxAttachments:   10,
		AllowedFileTypes: []string{"image/jpeg", "image/png"},
	})

	urlService := urlshortener.NewService(db, urlshortener.Config{
		BaseURL:      "http://localhost:8080",
		MaxURLLength: 2048,
		MinURLLength: 5,
	})

	// Register routes
	api := router.Group("/api/v1")
	messagingService.RegisterRoutes(api)
	urlService.RegisterRoutes(api)

	return router, db
}

// Benchmark messaging service endpoints
func BenchmarkMessagingService(b *testing.B) {
	router, db := setupTestRouter()

	// Create test user and channel
	userID := uuid.New()
	channelID := uuid.New()

	channel := &entity.Channel{
		ID:          channelID,
		Name:        "test-channel",
		Description: "Test channel for benchmarking",
		Type:        "public",
		CreatedBy:   userID,
	}
	db.Create(channel)

	member := &entity.ChannelMember{
		ChannelID: channelID,
		UserID:    userID,
		Role:      "owner",
	}
	db.Create(member)

	b.Run("SendMessage", func(b *testing.B) {
		memBefore := getMemoryStats()
		latencies := make([]time.Duration, b.N)

		b.ResetTimer()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			start := time.Now()

			reqBody := map[string]interface{}{
				"content":      fmt.Sprintf("Test message %d", i),
				"message_type": "text",
			}
			jsonBody, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/messaging/channels/%s/messages", channelID), bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")

			// Mock authentication middleware
			req = req.WithContext(req.Context())

			w := httptest.NewRecorder()

			// Add user_id to context manually for testing
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("user_id", userID.String())

			router.ServeHTTP(w, req)

			latencies[i] = time.Since(start)

			if w.Code != http.StatusCreated {
				b.Errorf("Expected status 201, got %d", w.Code)
			}
		}

		b.StopTimer()

		memAfter := getMemoryStats()

		// Calculate latency percentiles
		p95, p99 := calculatePercentiles(latencies)

		result := BenchmarkResult{
			TestName:       "MessagingService.SendMessage",
			Duration:       time.Duration(b.N) * time.Nanosecond,
			RequestsPerSec: float64(b.N) / b.Elapsed().Seconds(),
			MemoryUsage:    memAfter,
			P95Latency:     p95,
			P99Latency:     p99,
		}

		logBenchmarkResult(b, result)

		// Performance assertions
		if result.RequestsPerSec < 100 {
			b.Errorf("Low throughput: %.2f req/sec (expected > 100)", result.RequestsPerSec)
		}

		if memAfter.AllocMB-memBefore.AllocMB > 50 {
			b.Errorf("High memory usage: %.2f MB increase", memAfter.AllocMB-memBefore.AllocMB)
		}
	})

	b.Run("GetMessages", func(b *testing.B) {
		// Pre-populate with test messages
		for i := 0; i < 100; i++ {
			msg := &entity.Message{
				ChannelID:   channelID,
				UserID:      userID,
				Content:     fmt.Sprintf("Benchmark message %d", i),
				MessageType: "text",
			}
			db.Create(msg)
		}

		memBefore := getMemoryStats()
		latencies := make([]time.Duration, b.N)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			start := time.Now()

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/messaging/channels/%s/messages?page=1&page_size=20", channelID), nil)
			req.Header.Set("Authorization", "Bearer test-token")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("user_id", userID.String())

			router.ServeHTTP(w, req)

			latencies[i] = time.Since(start)

			if w.Code != http.StatusOK {
				b.Errorf("Expected status 200, got %d", w.Code)
			}
		}

		memAfter := getMemoryStats()
		p95, p99 := calculatePercentiles(latencies)

		result := BenchmarkResult{
			TestName:       "MessagingService.GetMessages",
			Duration:       b.Elapsed(),
			RequestsPerSec: float64(b.N) / b.Elapsed().Seconds(),
			MemoryUsage:    memAfter,
			P95Latency:     p95,
			P99Latency:     p99,
		}

		logBenchmarkResult(b, result)
	})
}

// Benchmark URL shortener service
func BenchmarkURLShortenerService(b *testing.B) {
	router, _ := setupTestRouter()

	b.Run("ShortenURL", func(b *testing.B) {
		memBefore := getMemoryStats()
		latencies := make([]time.Duration, b.N)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			start := time.Now()

			reqBody := map[string]interface{}{
				"url":         fmt.Sprintf("https://example.com/test-page-%d", i),
				"title":       fmt.Sprintf("Test Page %d", i),
				"description": "Test page for benchmarking",
			}
			jsonBody, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/api/v1/urls/shorten", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			latencies[i] = time.Since(start)

			if w.Code != http.StatusCreated {
				b.Errorf("Expected status 201, got %d", w.Code)
			}
		}

		memAfter := getMemoryStats()
		p95, p99 := calculatePercentiles(latencies)

		result := BenchmarkResult{
			TestName:       "URLShortenerService.ShortenURL",
			Duration:       b.Elapsed(),
			RequestsPerSec: float64(b.N) / b.Elapsed().Seconds(),
			MemoryUsage:    memAfter,
			P95Latency:     p95,
			P99Latency:     p99,
		}

		logBenchmarkResult(b, result)

		// Performance assertions
		if result.RequestsPerSec < 200 {
			b.Errorf("Low throughput: %.2f req/sec (expected > 200)", result.RequestsPerSec)
		}
	})
}

// Benchmark concurrent operations
func BenchmarkConcurrentOperations(b *testing.B) {
	router, db := setupTestRouter()

	// Setup test data
	userID := uuid.New()
	channelID := uuid.New()

	channel := &entity.Channel{
		ID:          channelID,
		Name:        "concurrent-test",
		Description: "Test channel for concurrent benchmarking",
		Type:        "public",
		CreatedBy:   userID,
	}
	db.Create(channel)

	member := &entity.ChannelMember{
		ChannelID: channelID,
		UserID:    userID,
		Role:      "owner",
	}
	db.Create(member)

	b.Run("ConcurrentMessageSending", func(b *testing.B) {
		memBefore := getMemoryStats()

		concurrency := 10
		messagesPerGoroutine := b.N / concurrency

		var wg sync.WaitGroup
		latencies := make([]time.Duration, b.N)
		latencyMutex := sync.Mutex{}
		latencyIndex := 0

		b.ResetTimer()

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				for j := 0; j < messagesPerGoroutine; j++ {
					start := time.Now()

					reqBody := map[string]interface{}{
						"content":      fmt.Sprintf("Concurrent message %d-%d", goroutineID, j),
						"message_type": "text",
					}
					jsonBody, _ := json.Marshal(reqBody)

					req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/messaging/channels/%s/messages", channelID), bytes.NewBuffer(jsonBody))
					req.Header.Set("Content-Type", "application/json")

					w := httptest.NewRecorder()
					c, _ := gin.CreateTestContext(w)
					c.Request = req
					c.Set("user_id", userID.String())

					router.ServeHTTP(w, req)

					duration := time.Since(start)

					latencyMutex.Lock()
					if latencyIndex < len(latencies) {
						latencies[latencyIndex] = duration
						latencyIndex++
					}
					latencyMutex.Unlock()

					if w.Code != http.StatusCreated {
						b.Errorf("Expected status 201, got %d", w.Code)
					}
				}
			}(i)
		}

		wg.Wait()

		memAfter := getMemoryStats()
		p95, p99 := calculatePercentiles(latencies[:latencyIndex])

		result := BenchmarkResult{
			TestName:       "ConcurrentMessageSending",
			Duration:       b.Elapsed(),
			RequestsPerSec: float64(b.N) / b.Elapsed().Seconds(),
			MemoryUsage:    memAfter,
			P95Latency:     p95,
			P99Latency:     p99,
		}

		logBenchmarkResult(b, result)

		// Performance assertions for concurrent operations
		if result.RequestsPerSec < 50 {
			b.Errorf("Low concurrent throughput: %.2f req/sec (expected > 50)", result.RequestsPerSec)
		}

		if memAfter.AllocMB-memBefore.AllocMB > 100 {
			b.Errorf("High memory usage in concurrent test: %.2f MB increase", memAfter.AllocMB-memBefore.AllocMB)
		}
	})
}

// Memory leak detection test
func BenchmarkMemoryLeakDetection(b *testing.B) {
	router, _ := setupTestRouter()

	b.Run("MemoryLeakTest", func(b *testing.B) {
		memBefore := getMemoryStats()

		// Perform many operations
		for i := 0; i < b.N; i++ {
			reqBody := map[string]interface{}{
				"url":   fmt.Sprintf("https://example.com/page-%d", i),
				"title": fmt.Sprintf("Page %d", i),
			}
			jsonBody, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/api/v1/urls/shorten", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if i%100 == 0 {
				runtime.GC() // Force garbage collection periodically
			}
		}

		// Force final GC
		runtime.GC()
		runtime.GC() // Run twice to ensure cleanup

		memAfter := getMemoryStats()
		memoryGrowth := memAfter.AllocMB - memBefore.AllocMB

		b.Logf("Memory growth: %.2f MB (before: %.2f MB, after: %.2f MB)",
			memoryGrowth, memBefore.AllocMB, memAfter.AllocMB)

		// Memory growth should be reasonable
		if memoryGrowth > 50 {
			b.Errorf("Potential memory leak detected: %.2f MB growth", memoryGrowth)
		}

		// GC should have run at least once
		if memAfter.NumGC <= memBefore.NumGC {
			b.Logf("Warning: GC did not run during test")
		}
	})
}

// Database query performance test
func BenchmarkDatabaseQueries(b *testing.B) {
	db := setupTestDB()

	// Pre-populate with test data
	userID := uuid.New()
	channelID := uuid.New()

	// Create channel
	channel := &entity.Channel{
		ID:          channelID,
		Name:        "query-test",
		Description: "Channel for query benchmarking",
		Type:        "public",
		CreatedBy:   userID,
	}
	db.Create(channel)

	// Create many messages
	for i := 0; i < 1000; i++ {
		msg := &entity.Message{
			ChannelID:   channelID,
			UserID:      userID,
			Content:     fmt.Sprintf("Query test message %d", i),
			MessageType: "text",
		}
		db.Create(msg)
	}

	b.Run("PaginatedMessageQuery", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var messages []entity.Message
			err := db.Where("channel_id = ?", channelID).
				Order("created_at DESC").
				Limit(20).
				Offset(i % 10 * 20). // Simulate different pages
				Find(&messages).Error

			if err != nil {
				b.Errorf("Database query failed: %v", err)
			}

			if len(messages) == 0 {
				b.Errorf("No messages returned")
			}
		}
	})

	b.Run("MessageCountQuery", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var count int64
			err := db.Model(&entity.Message{}).
				Where("channel_id = ?", channelID).
				Count(&count).Error

			if err != nil {
				b.Errorf("Count query failed: %v", err)
			}

			if count == 0 {
				b.Errorf("Expected non-zero count")
			}
		}
	})
}

// Helper functions

func calculatePercentiles(latencies []time.Duration) (p95, p99 time.Duration) {
	if len(latencies) == 0 {
		return 0, 0
	}

	// Sort latencies
	for i := 0; i < len(latencies)-1; i++ {
		for j := i + 1; j < len(latencies); j++ {
			if latencies[i] > latencies[j] {
				latencies[i], latencies[j] = latencies[j], latencies[i]
			}
		}
	}

	p95Index := int(0.95 * float64(len(latencies)))
	p99Index := int(0.99 * float64(len(latencies)))

	if p95Index >= len(latencies) {
		p95Index = len(latencies) - 1
	}
	if p99Index >= len(latencies) {
		p99Index = len(latencies) - 1
	}

	return latencies[p95Index], latencies[p99Index]
}

func logBenchmarkResult(b *testing.B, result BenchmarkResult) {
	b.Logf("Benchmark Results for %s:", result.TestName)
	b.Logf("  Duration: %v", result.Duration)
	b.Logf("  Requests/sec: %.2f", result.RequestsPerSec)
	b.Logf("  Memory (Alloc): %.2f MB", result.MemoryUsage.AllocMB)
	b.Logf("  Memory (Sys): %.2f MB", result.MemoryUsage.SysMB)
	b.Logf("  GC runs: %d", result.MemoryUsage.NumGC)
	b.Logf("  P95 Latency: %v", result.P95Latency)
	b.Logf("  P99 Latency: %v", result.P99Latency)

	// Write results to file for analysis
	if testing.Verbose() {
		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		log.Printf("Benchmark Result JSON:\n%s", string(resultJSON))
	}
}
