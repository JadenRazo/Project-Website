package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/mcstats/domain"
)

// StatsRepository implements domain.StatsRepository
type StatsRepository struct {
	db *sql.DB
}

// NewStatsRepository creates a new stats repository
func NewStatsRepository(mariadb *MariaDB) *StatsRepository {
	return &StatsRepository{db: mariadb.DB()}
}

// GetServerStats retrieves all server statistics from MariaDB
func (r *StatsRepository) GetServerStats(ctx context.Context) (*domain.ServerStats, error) {
	stats := &domain.ServerStats{
		Timestamp: time.Now().UnixMilli(),
	}

	// Total unique players
	if err := r.queryLong(ctx, &stats.TotalPlayers,
		"SELECT COUNT(DISTINCT uuid) FROM weenie_playtime"); err != nil {
		// Table might not exist yet, default to 0
		stats.TotalPlayers = 0
	}

	// Total playtime in hours
	var totalSeconds sql.NullInt64
	if err := r.queryNullableLong(ctx, &totalSeconds,
		"SELECT SUM(total_seconds) FROM weenie_playtime"); err == nil && totalSeconds.Valid {
		stats.TotalPlaytime = totalSeconds.Int64 / 3600
	}

	// Blocks placed
	var blocksPlaced sql.NullInt64
	if err := r.queryNullableLong(ctx, &blocksPlaced,
		"SELECT SUM(xp_amount) FROM weenie_claim_xp_history WHERE xp_source = 'BLOCKS_PLACED'"); err == nil && blocksPlaced.Valid {
		stats.BlocksPlaced = blocksPlaced.Int64
	}

	// Mobs killed (bounty claims as proxy)
	var mobsKilled sql.NullInt64
	if err := r.queryNullableLong(ctx, &mobsKilled,
		"SELECT COUNT(*) FROM weenie_bounty_claims"); err == nil && mobsKilled.Valid {
		stats.MobsKilled = mobsKilled.Int64
	}

	// Unique players today
	if err := r.queryLong(ctx, &stats.UniquePlayersToday,
		"SELECT COUNT(DISTINCT uuid) FROM weenie_playtime WHERE last_join >= DATE_SUB(NOW(), INTERVAL 24 HOUR)"); err != nil {
		stats.UniquePlayersToday = 0
	}

	// Peak online today (estimated as 1/4 of unique players today, minimum 1)
	stats.PeakOnlineToday = stats.UniquePlayersToday / 4
	if stats.PeakOnlineToday < 1 && stats.UniquePlayersToday > 0 {
		stats.PeakOnlineToday = 1
	}

	// Total claims
	if err := r.queryLong(ctx, &stats.TotalClaims,
		"SELECT COUNT(*) FROM weenie_claims"); err != nil {
		stats.TotalClaims = 0
	}

	// Total chunks claimed
	if err := r.queryLong(ctx, &stats.TotalChunksClaimed,
		"SELECT COUNT(*) FROM weenie_chunks"); err != nil {
		stats.TotalChunksClaimed = 0
	}

	return stats, nil
}

func (r *StatsRepository) queryLong(ctx context.Context, dest *int64, query string) error {
	return r.db.QueryRowContext(ctx, query).Scan(dest)
}

func (r *StatsRepository) queryNullableLong(ctx context.Context, dest *sql.NullInt64, query string) error {
	return r.db.QueryRowContext(ctx, query).Scan(dest)
}

// Close is a no-op as connection is managed by MariaDB wrapper
func (r *StatsRepository) Close() error {
	return nil
}

// Ping checks database connectivity
func (r *StatsRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
