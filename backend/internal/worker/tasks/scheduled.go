package tasks

import (
	"context"

	"github.com/JadenRazo/Project-Website/backend/internal/codestats"
	"github.com/JadenRazo/Project-Website/backend/internal/codestats/projectpath/repository"
	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type ScheduledTasks struct {
	cron               *cron.Cron
	codeStatsService   *codestats.Service
	visitorMetricsTask *VisitorMetricsTask
}

func NewScheduledTasks(db *gorm.DB) *ScheduledTasks {
	projectPathRepo := repository.NewGormRepository(db)

	return &ScheduledTasks{
		cron:               cron.New(cron.WithSeconds()),
		codeStatsService:   codestats.NewService(db, projectPathRepo),
		visitorMetricsTask: NewVisitorMetricsTask(db),
	}
}

func (st *ScheduledTasks) Start(ctx context.Context) error {
	logger.Info("Starting scheduled tasks")

	_, err := st.cron.AddFunc("0 0 * * * *", func() {
		if err := st.visitorMetricsTask.AggregateHourlyMetrics(ctx); err != nil {
			logger.Error("Failed to aggregate hourly visitor metrics", "error", err)
		} else {
			logger.Info("Hourly visitor metrics aggregated successfully")
		}
	})
	if err != nil {
		logger.Error("Failed to schedule hourly metrics aggregation", "error", err)
	}

	_, err = st.cron.AddFunc("0 5 0 * * *", func() {
		if err := st.visitorMetricsTask.GenerateDailySummary(ctx); err != nil {
			logger.Error("Failed to generate daily visitor summary", "error", err)
		} else {
			logger.Info("Daily visitor summary generated successfully")
		}
	})
	if err != nil {
		logger.Error("Failed to schedule daily summary generation", "error", err)
	}

	_, err = st.cron.AddFunc("0 0 */6 * * *", func() {
		if err := st.visitorMetricsTask.CleanupExpiredSessions(ctx); err != nil {
			logger.Error("Failed to cleanup expired visitor sessions", "error", err)
		} else {
			logger.Info("Expired visitor sessions cleaned up successfully")
		}
	})
	if err != nil {
		logger.Error("Failed to schedule session cleanup", "error", err)
	}

	_, err = st.cron.AddFunc("0 0 2 * * *", func() {
		if err := st.visitorMetricsTask.UpdateLocationAggregates(ctx); err != nil {
			logger.Error("Failed to update visitor location aggregates", "error", err)
		} else {
			logger.Info("Visitor location aggregates updated successfully")
		}
	})
	if err != nil {
		logger.Error("Failed to schedule location aggregates update", "error", err)
	}

	_, err = st.cron.AddFunc("0 0 3 * * 0", func() {
		if err := st.visitorMetricsTask.CleanupOldMetrics(ctx, 90); err != nil {
			logger.Error("Failed to cleanup old visitor metrics", "error", err)
		} else {
			logger.Info("Old visitor metrics cleaned up successfully")
		}
	})
	if err != nil {
		logger.Error("Failed to schedule old metrics cleanup", "error", err)
	}

	logger.Info("Visitor analytics scheduled tasks registered",
		"hourly_aggregation", "0 0 * * * *",
		"daily_summary", "0 5 0 * * *",
		"session_cleanup", "0 0 */6 * * *",
		"location_aggregates", "0 0 2 * * *",
		"metrics_cleanup", "0 0 3 * * 0",
	)

	st.cron.Start()

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
