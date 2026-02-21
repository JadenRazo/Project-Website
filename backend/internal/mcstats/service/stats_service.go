package service

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/common/cache"
	"github.com/JadenRazo/Project-Website/backend/internal/mcstats/domain"
)

const (
	statsCacheKey = "mcstats:server_stats"
	statsCacheTTL = 60 * time.Second
)

// StatsService implements domain.StatsService with caching
type StatsService struct {
	repo  domain.StatsRepository
	cache cache.Cache
}

// NewStatsService creates a new stats service
func NewStatsService(repo domain.StatsRepository, c cache.Cache) *StatsService {
	return &StatsService{
		repo:  repo,
		cache: c,
	}
}

// GetServerStats retrieves stats with caching
func (s *StatsService) GetServerStats(ctx context.Context) (*domain.ServerStats, error) {
	// Try cache first
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, statsCacheKey); err == nil && cached != nil {
			// Handle map[string]interface{} from JSON unmarshaling
			if statsMap, ok := cached.(map[string]interface{}); ok {
				return mapToServerStats(statsMap), nil
			}
		}
	}

	// Fetch from database
	stats, err := s.repo.GetServerStats(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if s.cache != nil {
		_ = s.cache.Set(ctx, statsCacheKey, stats, statsCacheTTL)
	}

	return stats, nil
}

func mapToServerStats(m map[string]interface{}) *domain.ServerStats {
	getInt64 := func(key string) int64 {
		if v, ok := m[key]; ok {
			switch val := v.(type) {
			case float64:
				return int64(val)
			case int64:
				return val
			case int:
				return int64(val)
			}
		}
		return 0
	}

	return &domain.ServerStats{
		TotalPlayers:       getInt64("totalPlayers"),
		TotalPlaytime:      getInt64("totalPlaytime"),
		BlocksPlaced:       getInt64("blocksPlaced"),
		MobsKilled:         getInt64("mobsKilled"),
		UniquePlayersToday: getInt64("uniquePlayersToday"),
		PeakOnlineToday:    getInt64("peakOnlineToday"),
		TotalClaims:        getInt64("totalClaims"),
		TotalChunksClaimed: getInt64("totalChunksClaimed"),
		Timestamp:          getInt64("timestamp"),
	}
}
