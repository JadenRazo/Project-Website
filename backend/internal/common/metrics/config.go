package metrics

// Config contains configuration for metrics collection
type Config struct {
	// Enabled indicates if metrics are enabled
	Enabled bool `json:"enabled" mapstructure:"enabled"`

	// Grafana configuration
	Grafana GrafanaConfig `json:"grafana" mapstructure:"grafana"`
}

// GrafanaConfig contains configuration for Grafana metrics
type GrafanaConfig struct {
	// Enabled determines if Grafana metrics are enabled
	Enabled bool `json:"enabled" mapstructure:"enabled"`

	// Namespace is the metrics namespace
	Namespace string `json:"namespace" mapstructure:"namespace"`

	// Subsystem is the metrics subsystem
	Subsystem string `json:"subsystem" mapstructure:"subsystem"`

	// Address is the Prometheus scrape endpoint
	Address string `json:"address" mapstructure:"address"`

	// Endpoint is the HTTP endpoint to expose metrics on
	Endpoint string `json:"endpoint" mapstructure:"endpoint"`

	// MetricsPrefix is the prefix added to all metrics
	MetricsPrefix string `json:"metrics_prefix" mapstructure:"metrics_prefix"`

	// DefaultLabels are labels added to all metrics
	DefaultLabels map[string]string `json:"default_labels" mapstructure:"default_labels"`
}

// NewDefaultConfig returns a new default configuration
func NewDefaultConfig() *Config {
	return &Config{
		Enabled: true,
		Grafana: GrafanaConfig{
			Enabled:       false,
			Namespace:     "api",
			Subsystem:     "http",
			Address:       ":9090",
			Endpoint:      "/metrics",
			MetricsPrefix: "app_",
			DefaultLabels: map[string]string{
				"service": "backend",
			},
		},
	}
}
