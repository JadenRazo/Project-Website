package domain

import "context"

// ServerStats represents aggregated Minecraft server statistics
type ServerStats struct {
	TotalPlayers       int64 `json:"totalPlayers"`
	TotalPlaytime      int64 `json:"totalPlaytime"`       // hours
	BlocksPlaced       int64 `json:"blocksPlaced"`
	MobsKilled         int64 `json:"mobsKilled"`
	UniquePlayersToday int64 `json:"uniquePlayersToday"`
	PeakOnlineToday    int64 `json:"peakOnlineToday"`
	TotalClaims        int64 `json:"totalClaims"`
	TotalChunksClaimed int64 `json:"totalChunksClaimed"`
	Timestamp          int64 `json:"timestamp"`
}

// StatsRepository defines the interface for stats data access
type StatsRepository interface {
	GetServerStats(ctx context.Context) (*ServerStats, error)
	Close() error
	Ping(ctx context.Context) error
}

// StatsService defines the interface for stats business logic
type StatsService interface {
	GetServerStats(ctx context.Context) (*ServerStats, error)
}
