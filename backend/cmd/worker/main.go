package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	appconfig "github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/JadenRazo/Project-Website/backend/internal/common/database"
	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"gorm.io/gorm"
)

const (
	appName    = "project-website-worker"
	appVersion = "1.0.0"
)

// Job represents a background task to be processed
type Job struct {
	ID          string
	Type        string
	Payload     []byte
	CreatedAt   time.Time
	Status      string
	RetryCount  int
	NextRetryAt *time.Time
	Error       string
}

// Worker handles processing of background jobs
type Worker struct {
	dbConn        *gorm.DB
	quit          chan struct{}
	processingWg  sync.WaitGroup
	pollInterval  time.Duration
	maxRetries    int
	processorsMap map[string]JobProcessor
	batchSize     int
	logger        *log.Logger
}

// JobProcessor defines the interface for job type-specific processors
type JobProcessor interface {
	Process(ctx context.Context, job Job) error
}

// EmailProcessor processes email jobs
type EmailProcessor struct{}

func (p *EmailProcessor) Process(ctx context.Context, job Job) error {
	// Implement email sending logic
	log.Printf("Processing email job: %s", job.ID)
	return nil
}

// NotificationProcessor processes notification jobs
type NotificationProcessor struct{}

func (p *NotificationProcessor) Process(ctx context.Context, job Job) error {
	// Implement notification sending logic
	log.Printf("Processing notification job: %s", job.ID)
	return nil
}

// NewWorker creates a new worker instance
func NewWorker(database *gorm.DB) *Worker {
	return &Worker{
		dbConn:       database,
		quit:         make(chan struct{}),
		pollInterval: 5 * time.Second,
		maxRetries:   3,
		batchSize:    10,
		logger:       log.New(os.Stdout, "[WORKER] ", log.LstdFlags),
		processorsMap: map[string]JobProcessor{
			"email":        &EmailProcessor{},
			"notification": &NotificationProcessor{},
		},
	}
}

// Start begins the worker's job processing loop
func (w *Worker) Start(ctx context.Context) error {
	w.logger.Println("Starting worker service...")

	// Create a ticker for regular job polling
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	// Start processing loop
	for {
		select {
		case <-ticker.C:
			w.processBatch(ctx)
		case <-w.quit:
			w.logger.Println("Worker shutting down...")
			return nil
		case <-ctx.Done():
			w.logger.Println("Context cancelled, worker shutting down...")
			return ctx.Err()
		}
	}
}

// processBatch processes a batch of jobs
func (w *Worker) processBatch(ctx context.Context) {
	jobs, err := w.fetchJobs(ctx, w.batchSize)
	if err != nil {
		w.logger.Printf("Error fetching jobs: %v", err)
		return
	}

	if len(jobs) == 0 {
		return // No jobs to process
	}

	w.logger.Printf("Processing batch of %d jobs", len(jobs))

	for _, job := range jobs {
		// Process each job in its own goroutine with tracking
		w.processingWg.Add(1)
		go func(j Job) {
			defer w.processingWg.Done()
			defer func() {
				if r := recover(); r != nil {
					w.logger.Printf("Recovered from panic in job %s: %v", j.ID, r)
					w.markJobFailed(ctx, j.ID, fmt.Sprintf("panic: %v", r))
				}
			}()

			if err := w.processJob(ctx, j); err != nil {
				w.logger.Printf("Failed to process job %s: %v", j.ID, err)
				w.markJobFailed(ctx, j.ID, err.Error())
			} else {
				w.markJobComplete(ctx, j.ID)
			}
		}(job)
	}
}

// fetchJobs retrieves pending jobs from the database
func (w *Worker) fetchJobs(ctx context.Context, limit int) ([]Job, error) {
	var jobs []Job

	// In a production environment, we would use a transaction and update job status
	// to prevent multiple workers from processing the same jobs
	tx := w.dbConn.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// This is where you would implement your database query
	// Example query (commented out for now):
	/*
		err := tx.Raw(`
			UPDATE jobs
			SET status = 'processing', updated_at = NOW()
			WHERE id IN (
				SELECT id FROM jobs
				WHERE status = 'pending'
				AND next_retry_at <= NOW()
				ORDER BY priority DESC, created_at ASC
				LIMIT ?
				FOR UPDATE SKIP LOCKED
			)
			RETURNING *
		`, limit).Scan(&jobs).Error

		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to fetch jobs: %w", err)
		}
	*/

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return empty slice for now
	return jobs, nil
}

// processJob handles processing a single job
func (w *Worker) processJob(ctx context.Context, job Job) error {
	processor, exists := w.processorsMap[job.Type]
	if !exists {
		return fmt.Errorf("no processor found for job type: %s", job.Type)
	}

	w.logger.Printf("Processing job %s of type %s", job.ID, job.Type)

	// Create a timeout context for individual job processing
	jobCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	return processor.Process(jobCtx, job)
}

// markJobComplete marks a job as successfully completed
func (w *Worker) markJobComplete(ctx context.Context, jobID string) {
	// Update job status in database
	tx := w.dbConn.WithContext(ctx).Begin()
	if tx.Error != nil {
		w.logger.Printf("Failed to begin transaction for completing job %s: %v", jobID, tx.Error)
		return
	}

	// Example update query (commented out for now)
	/*
		err := tx.Exec(`
			UPDATE jobs
			SET status = 'completed', completed_at = NOW(), updated_at = NOW()
			WHERE id = ?
		`, jobID).Error

		if err != nil {
			tx.Rollback()
			w.logger.Printf("Failed to mark job %s as complete: %v", jobID, err)
			return
		}
	*/

	if err := tx.Commit().Error; err != nil {
		w.logger.Printf("Failed to commit transaction for completing job %s: %v", jobID, err)
		return
	}

	w.logger.Printf("Marked job %s as complete", jobID)
}

// markJobFailed marks a job as failed and updates retry information
func (w *Worker) markJobFailed(ctx context.Context, jobID string, errMsg string) {
	tx := w.dbConn.WithContext(ctx).Begin()
	if tx.Error != nil {
		w.logger.Printf("Failed to begin transaction for failing job %s: %v", jobID, tx.Error)
		return
	}

	// Example update query with retry logic (commented out for now)
	/*
		var job Job
		err := tx.First(&job, "id = ?", jobID).Error
		if err != nil {
			tx.Rollback()
			w.logger.Printf("Failed to find job %s: %v", jobID, err)
			return
		}

		job.RetryCount++
		job.Error = errMsg

		if job.RetryCount >= w.maxRetries {
			job.Status = "failed"
		} else {
			job.Status = "pending"
			// Calculate exponential backoff for next retry
			backoff := time.Duration(1<<uint(job.RetryCount-1)) * time.Second
			nextRetry := time.Now().Add(backoff)
			job.NextRetryAt = &nextRetry
		}

		err = tx.Save(&job).Error
		if err != nil {
			tx.Rollback()
			w.logger.Printf("Failed to update job %s after failure: %v", jobID, err)
			return
		}
	*/

	if err := tx.Commit().Error; err != nil {
		w.logger.Printf("Failed to commit transaction for failing job %s: %v", jobID, err)
		return
	}

	w.logger.Printf("Marked job %s as failed: %s", jobID, errMsg)
}

// GracefulShutdown waits for all in-progress jobs to complete
func (w *Worker) GracefulShutdown(timeout time.Duration) {
	w.logger.Printf("Graceful shutdown initiated with %s timeout", timeout)

	// Signal the main loop to stop polling for new jobs
	close(w.quit)

	// Wait for in-progress jobs with timeout
	done := make(chan struct{})
	go func() {
		w.processingWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		w.logger.Println("All jobs completed successfully")
	case <-time.After(timeout):
		w.logger.Println("Shutdown timeout reached, some jobs may still be running")
	}
}

func main() {
	// Capture and handle panics
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("PANIC: %v", r)
		}
	}()

	// Setup context for the entire application
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load environment configuration
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Initialize the logging system
	logConfig := &appconfig.LoggingConfig{
		Level:      "info",
		Format:     "json",
		Output:     "file",
		TimeFormat: time.RFC3339,
		Filename:   "logs/worker.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}

	err := logger.InitLogger(logConfig, appName, appVersion)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Shutdown()

	// Log startup
	logger.Info("Worker service starting up", "name", appName, "version", appVersion, "environment", env)

	var dbConn *gorm.DB

	// Try to use the new database package first
	dbConfig := &appconfig.DatabaseConfig{
		Driver:                 "postgres",
		DSN:                    os.Getenv("DATABASE_URL"),
		MaxIdleConns:           5,
		MaxOpenConns:           20,
		ConnMaxLifetimeMinutes: 60,
		LogLevel:               "warn",
		SlowThresholdMs:        200,
	}

	// Try to connect using the database package
	logger.Info("Connecting to database...")

	// First attempt with the common database package
	dbInstance, err := database.NewDatabase(dbConfig)
	if err == nil {
		dbConn = dbInstance.GetDB()
		logger.Info("Connected to database using database package")
	} else {
		// Fallback to core/db package
		logger.Warn("Failed to connect with database package, trying fallback", "error", err.Error())

		dbConn, err = db.GetDB()
		if err != nil {
			logger.Fatal("Failed to connect to database", "error", err)
		}
		logger.Info("Connected to database using fallback method")
	}

	// Initialize worker with database connection
	worker := NewWorker(dbConn)

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start worker in background
	go func() {
		if err := worker.Start(ctx); err != nil && err != context.Canceled {
			logger.Fatal("Worker failed", "error", err)
		}
	}()

	// Wait for termination signal
	sig := <-quit
	logger.Info("Received shutdown signal", "signal", sig.String())

	// Cancel context to signal worker to stop
	cancel()

	// Start graceful shutdown with timeout
	worker.GracefulShutdown(30 * time.Second)

	logger.Info("Worker service shutdown complete")
}
