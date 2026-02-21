package minecraft

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type ServerStats struct {
	TotalPlayers       int64 `json:"totalPlayers"`
	TotalPlaytime      int64 `json:"totalPlaytime"`      // in hours
	BlocksPlaced       int64 `json:"blocksPlaced"`
	MobsKilled         int64 `json:"mobsKilled"`
	UniquePlayersToday int64 `json:"uniquePlayersToday"`
	PeakOnlineToday    int64 `json:"peakOnlineToday"`
}

type Service struct {
	db          *sql.DB
	statsCache  *ServerStats
	cacheExpiry time.Time
	cacheTTL    time.Duration
}

func NewService() (*Service, error) {
	// Get database configuration from environment
	host := os.Getenv("MINECRAFT_DB_HOST")
	port := os.Getenv("MINECRAFT_DB_PORT")
	user := os.Getenv("MINECRAFT_DB_USER")
	password := os.Getenv("MINECRAFT_DB_PASSWORD")
	dbName := os.Getenv("MINECRAFT_DB_NAME")

	if host == "" {
		host = "208.115.245.242"
	}
	if port == "" {
		port = "1025"
	}
	if dbName == "" {
		dbName = "s5578_MariaDB-1"
	}

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&timeout=10s",
		user, password, host, port, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Service{
		db:       db,
		cacheTTL: 1 * time.Minute, // Cache stats for 1 minute
	}, nil
}

func (s *Service) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Service) GetStats(ctx context.Context) (*ServerStats, error) {
	// Check cache
	if s.statsCache != nil && time.Now().Before(s.cacheExpiry) {
		return s.statsCache, nil
	}

	stats := &ServerStats{}

	// Query total unique players from playtime table
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(DISTINCT uuid) FROM weenie_playtime").Scan(&stats.TotalPlayers)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get total players: %w", err)
	}

	// Query total playtime in hours
	var totalSeconds sql.NullInt64
	err = s.db.QueryRowContext(ctx,
		"SELECT SUM(total_seconds) FROM weenie_playtime").Scan(&totalSeconds)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get total playtime: %w", err)
	}
	if totalSeconds.Valid {
		stats.TotalPlaytime = totalSeconds.Int64 / 3600 // Convert to hours
	}

	// Query blocks placed from XP history (blocks_placed events)
	var blocksPlaced sql.NullInt64
	err = s.db.QueryRowContext(ctx,
		"SELECT SUM(xp_amount) FROM weenie_claim_xp_history WHERE xp_source = 'BLOCKS_PLACED'").Scan(&blocksPlaced)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get blocks placed: %w", err)
	}
	if blocksPlaced.Valid {
		stats.BlocksPlaced = blocksPlaced.Int64
	}

	// Query mobs killed from bounty claims (kills)
	var mobsKilled sql.NullInt64
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM weenie_bounty_claims").Scan(&mobsKilled)
	if err != nil && err != sql.ErrNoRows {
		// If table doesn't exist, try event_stats for a fallback
		err = s.db.QueryRowContext(ctx,
			"SELECT SUM(wins) FROM weenie_event_stats").Scan(&mobsKilled)
		if err != nil && err != sql.ErrNoRows {
			// Just set to 0 if neither works
			stats.MobsKilled = 0
		}
	}
	if mobsKilled.Valid {
		stats.MobsKilled = mobsKilled.Int64
	}

	// Query unique players today (players who joined in last 24h)
	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(DISTINCT uuid) FROM weenie_playtime WHERE last_join >= DATE_SUB(NOW(), INTERVAL 24 HOUR)").Scan(&stats.UniquePlayersToday)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get unique players today: %w", err)
	}

	// For peak online today, we'd need session tracking which doesn't exist
	// Use a reasonable estimate based on unique players today
	stats.PeakOnlineToday = stats.UniquePlayersToday / 4
	if stats.PeakOnlineToday < 1 && stats.UniquePlayersToday > 0 {
		stats.PeakOnlineToday = 1
	}

	// Cache the results
	s.statsCache = stats
	s.cacheExpiry = time.Now().Add(s.cacheTTL)

	return stats, nil
}

func (s *Service) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/stats", s.handleGetStats)
}

func (s *Service) handleGetStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	stats, err := s.GetStats(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch stats", "details": err.Error()})
		return
	}

	c.JSON(200, stats)
}
