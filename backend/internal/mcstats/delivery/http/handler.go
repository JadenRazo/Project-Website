package http

import (
	"encoding/json"
	"net/http"

	"github.com/JadenRazo/Project-Website/backend/internal/mcstats/domain"
)

// Handler handles HTTP requests for mcstats
type Handler struct {
	service domain.StatsService
}

// NewHandler creates a new handler
func NewHandler(service domain.StatsService) *Handler {
	return &Handler{service: service}
}

// GetServerStats handles GET /api/mc/stats
func (h *Handler) GetServerStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := h.service.GetServerStats(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch server stats",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	json.NewEncoder(w).Encode(stats)
}

// HealthCheck handles GET /api/mc/health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "mcstats",
	})
}
