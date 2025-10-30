package worker

import (
	"context"
	"sync"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/core"
	"github.com/JadenRazo/Project-Website/backend/internal/worker/tasks"
	"gorm.io/gorm"
)

type Service struct {
	*core.BaseService
	db             *gorm.DB
	scheduledTasks *tasks.ScheduledTasks
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.Mutex
}

func NewService(db *gorm.DB) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &Service{
		BaseService:    core.NewBaseService("worker"),
		db:             db,
		scheduledTasks: tasks.NewScheduledTasks(db),
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (s *Service) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.BaseService.Start(); err != nil {
		return err
	}

	go func() {
		if err := s.scheduledTasks.Start(s.ctx); err != nil && err != context.Canceled {
			s.AddError(err)
		}
	}()

	return nil
}

func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cancel()

	s.scheduledTasks.Stop()

	return s.BaseService.Stop()
}

func (s *Service) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}

	time.Sleep(500 * time.Millisecond)

	s.ctx, s.cancel = context.WithCancel(context.Background())

	return s.Start()
}

func (s *Service) HealthCheck() error {
	if err := s.BaseService.HealthCheck(); err != nil {
		return err
	}

	return nil
}
