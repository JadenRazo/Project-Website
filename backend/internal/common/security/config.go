package security

import (
	"time"
)

type SecurityConfig struct {
	DevPanel    DevPanelSecurity `yaml:"devpanel" json:"devpanel"`
	RateLimit   RateLimitConfig  `yaml:"rate_limit" json:"rate_limit"`
	CORS        CORSConfig       `yaml:"cors" json:"cors"`
	AuditLog    AuditLogConfig   `yaml:"audit_log" json:"audit_log"`
	Environment string           `yaml:"environment" json:"environment"`
}

type DevPanelSecurity struct {
	IPWhitelist    IPWhitelistConfig `yaml:"ip_whitelist" json:"ip_whitelist"`
	BindAddress    string            `yaml:"bind_address" json:"bind_address"`
	LocalhostOnly  bool              `yaml:"localhost_only" json:"localhost_only"`
	RequireTLS     bool              `yaml:"require_tls" json:"require_tls"`
	SessionTimeout time.Duration     `yaml:"session_timeout" json:"session_timeout"`
	MaxSessions    int               `yaml:"max_sessions" json:"max_sessions"`
}

type RateLimitConfig struct {
	Enabled         bool          `yaml:"enabled" json:"enabled"`
	RequestsPerMin  int           `yaml:"requests_per_minute" json:"requests_per_minute"`
	BurstSize       int           `yaml:"burst_size" json:"burst_size"`
	CleanupInterval time.Duration `yaml:"cleanup_interval" json:"cleanup_interval"`
	AuthEndpoints   struct {
		RequestsPerMin int `yaml:"requests_per_minute" json:"requests_per_minute"`
		BurstSize      int `yaml:"burst_size" json:"burst_size"`
	} `yaml:"auth_endpoints" json:"auth_endpoints"`
}

type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins" json:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods" json:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers" json:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials" json:"allow_credentials"`
	MaxAge           int      `yaml:"max_age" json:"max_age"`
}

type AuditLogConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	LogFile    string `yaml:"log_file" json:"log_file"`
	MaxSize    int    `yaml:"max_size_mb" json:"max_size_mb"`
	MaxBackups int    `yaml:"max_backups" json:"max_backups"`
	MaxAge     int    `yaml:"max_age_days" json:"max_age_days"`
	Compress   bool   `yaml:"compress" json:"compress"`
}

func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		Environment: "development",
		DevPanel: DevPanelSecurity{
			IPWhitelist:    DefaultLocalhostConfig(),
			BindAddress:    "127.0.0.1",
			LocalhostOnly:  true,
			RequireTLS:     false,
			SessionTimeout: 24 * time.Hour,
			MaxSessions:    5,
		},
		RateLimit: RateLimitConfig{
			Enabled:         true,
			RequestsPerMin:  100,
			BurstSize:       10,
			CleanupInterval: 5 * time.Minute,
			AuthEndpoints: struct {
				RequestsPerMin int `yaml:"requests_per_minute" json:"requests_per_minute"`
				BurstSize      int `yaml:"burst_size" json:"burst_size"`
			}{
				RequestsPerMin: 10,
				BurstSize:      3,
			},
		},
		CORS: CORSConfig{
			AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
			AllowCredentials: true,
			MaxAge:           3600,
		},
		AuditLog: AuditLogConfig{
			Enabled:    true,
			LogFile:    "/var/log/devpanel-audit.log",
			MaxSize:    100,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		},
	}
}

func ProductionSecurityConfig() SecurityConfig {
	config := DefaultSecurityConfig()
	config.Environment = "production"
	config.DevPanel.RequireTLS = true
	config.DevPanel.SessionTimeout = 8 * time.Hour
	config.DevPanel.MaxSessions = 2
	config.RateLimit.RequestsPerMin = 50
	config.RateLimit.AuthEndpoints.RequestsPerMin = 5
	config.RateLimit.AuthEndpoints.BurstSize = 2
	return config
}

func (sc *SecurityConfig) IsProduction() bool {
	return sc.Environment == "production"
}

func (sc *SecurityConfig) IsDevelopment() bool {
	return sc.Environment == "development"
}

func (sc *SecurityConfig) GetBindAddress(port int) string {
	if sc.DevPanel.LocalhostOnly {
		return "127.0.0.1"
	}
	return sc.DevPanel.BindAddress
}
