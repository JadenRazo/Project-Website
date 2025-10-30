package metrics

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/common/response"
)

// SystemMetricsResponse represents the response for system metrics
type SystemMetricsResponse struct {
	Period            string          `json:"period"`
	HasSufficientData bool            `json:"has_sufficient_data"`
	Metrics           []LatencyMetric `json:"metrics"`
	Stats             LatencyStats    `json:"stats"`
	LastUpdated       time.Time       `json:"last_updated"`
	Message           string          `json:"message,omitempty"`
}

// SystemMetricsHandler handles requests for system metrics
type SystemMetricsHandler struct {
	manager *Manager
}

// NewSystemMetricsHandler creates a new system metrics handler
func NewSystemMetricsHandler(manager *Manager) *SystemMetricsHandler {
	return &SystemMetricsHandler{
		manager: manager,
	}
}

// HandleMetrics handles GET /api/metrics/{period}
func (h *SystemMetricsHandler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, "Method not allowed")
		return
	}

	// Extract period from URL path
	period := extractPeriodFromPath(r.URL.Path)
	if period == "" {
		period = string(TimePeriodDay) // Default to day
	}

	timePeriod := TimePeriod(period)

	// Validate period
	if !isValidPeriod(timePeriod) {
		response.BadRequest(w, "Invalid time period. Use 'day', 'week', or 'month'")
		return
	}

	metricsResponse := h.buildMetricsResponse(timePeriod)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	if err := json.NewEncoder(w).Encode(metricsResponse); err != nil {
		response.InternalError(w, "Failed to encode response")
		return
	}
}

// HandleMetricsDay handles GET /api/metrics/day
func (h *SystemMetricsHandler) HandleMetricsDay(w http.ResponseWriter, r *http.Request) {
	response := h.buildMetricsResponse(TimePeriodDay)
	h.writeJSONResponse(w, response)
}

// HandleMetricsWeek handles GET /api/metrics/week
func (h *SystemMetricsHandler) HandleMetricsWeek(w http.ResponseWriter, r *http.Request) {
	response := h.buildMetricsResponse(TimePeriodWeek)
	h.writeJSONResponse(w, response)
}

// HandleMetricsMonth handles GET /api/metrics/month
func (h *SystemMetricsHandler) HandleMetricsMonth(w http.ResponseWriter, r *http.Request) {
	response := h.buildMetricsResponse(TimePeriodMonth)
	h.writeJSONResponse(w, response)
}

// HandleLatestMetrics handles GET /api/metrics/latest
func (h *SystemMetricsHandler) HandleLatestMetrics(w http.ResponseWriter, r *http.Request) {
	latest := h.manager.GetLatestLatency()

	response := map[string]interface{}{
		"latest_metric": latest,
		"timestamp":     time.Now(),
	}

	if latest == nil {
		response["message"] = "No latency data available yet"
	}

	h.writeJSONResponse(w, response)
}

// buildMetricsResponse builds the metrics response for a given period
func (h *SystemMetricsHandler) buildMetricsResponse(period TimePeriod) SystemMetricsResponse {
	hasSufficientData := h.manager.HasSufficientLatencyData(period)

	response := SystemMetricsResponse{
		Period:            string(period),
		HasSufficientData: hasSufficientData,
		LastUpdated:       time.Now(),
	}

	if hasSufficientData {
		response.Metrics = h.manager.GetLatencyMetrics(period)
		response.Stats = h.manager.GetLatencyStats(period)
	} else {
		response.Metrics = []LatencyMetric{}
		response.Stats = LatencyStats{Period: string(period)}
		response.Message = h.getInsufficientDataMessage(period)
	}

	return response
}

// getInsufficientDataMessage returns an appropriate message for insufficient data
func (h *SystemMetricsHandler) getInsufficientDataMessage(period TimePeriod) string {
	switch period {
	case TimePeriodDay:
		return "We're still collecting data for the 24-hour view. Please check back in a few hours as we gather sufficient metrics to display meaningful latency trends."
	case TimePeriodWeek:
		return "Weekly analytics are not yet available. We need at least 24 hours of uptime data before we can provide comprehensive weekly latency insights."
	case TimePeriodMonth:
		return "Monthly analytics are not yet available. We need at least one week of uptime data before we can provide comprehensive monthly latency insights and trends."
	default:
		return "Insufficient data available for the requested time period."
	}
}

// writeJSONResponse writes a JSON response
func (h *SystemMetricsHandler) writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		response.InternalError(w, "Failed to encode response")
		return
	}
}

// RegisterRoutes registers the metrics routes
func (h *SystemMetricsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/metrics/day", h.HandleMetricsDay)
	mux.HandleFunc("/api/metrics/week", h.HandleMetricsWeek)
	mux.HandleFunc("/api/metrics/month", h.HandleMetricsMonth)
	mux.HandleFunc("/api/metrics/latest", h.HandleLatestMetrics)
	mux.HandleFunc("/api/metrics/", h.HandleMetrics)
}

// extractPeriodFromPath extracts the time period from URL path
func extractPeriodFromPath(path string) string {
	// Extract period from paths like "/api/metrics/day", "/api/metrics/week", etc.
	if len(path) > 13 { // "/api/metrics/" is 13 characters
		return path[13:]
	}
	return ""
}

// isValidPeriod checks if the period is valid
func isValidPeriod(period TimePeriod) bool {
	switch period {
	case TimePeriodDay, TimePeriodWeek, TimePeriodMonth:
		return true
	default:
		return false
	}
}
