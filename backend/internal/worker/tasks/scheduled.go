package tasks

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/codestats"
	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type ScheduledTasks struct {
	cron            *cron.Cron
	codeStatsService *codestats.Service
}

func NewScheduledTasks(db *gorm.DB) *ScheduledTasks {
	return &ScheduledTasks{
		cron: cron.New(cron.WithSeconds()),
		codeStatsService: codestats.NewService(db, codestats.Config{
			ProjectDir:     "/main/Project-Website",
			UpdateInterval: 1 * time.Hour,
		}),
	}
}

func (st *ScheduledTasks) Start(ctx context.Context) error {
	logger.Info("Starting scheduled tasks")

	// Update code stats every hour
	_, err := st.cron.AddFunc("0 0 * * * *", func() {
		logger.Info("Running scheduled code stats update")
		if err := st.codeStatsService.UpdateStats(); err != nil {
			logger.Error("Failed to update code stats", "error", err)
		}
	})
	if err != nil {
		return err
	}

	// Run immediately on startup
	go func() {
		logger.Info("Running initial code stats update")
		if err := st.codeStatsService.UpdateStats(); err != nil {
			logger.Error("Failed to update code stats on startup", "error", err)
		}
	}()

	st.cron.Start()

	// Wait for context cancellation
	<-ctx.Done()
	
	logger.Info("Stopping scheduled tasks")
	ctx = st.cron.Stop()
	
	return nil
}

func (st *ScheduledTasks) Stop() {
	if st.cron != nil {
		st.cron.Stop()
	}
}