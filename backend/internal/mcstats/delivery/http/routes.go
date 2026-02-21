package http

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers mcstats routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/mc", func(r chi.Router) {
		r.Get("/stats", h.GetServerStats)
		r.Get("/health", h.HealthCheck)
	})
}
