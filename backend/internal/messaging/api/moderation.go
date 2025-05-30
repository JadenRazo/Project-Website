package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/messaging/domain"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// ModerationHandler handles moderation-related HTTP requests
type ModerationHandler struct {
	moderationUC domain.ModerationUseCase
	wordFilterUC domain.WordFilterUseCase
}

// NewModerationHandler creates a new ModerationHandler
func NewModerationHandler(moderationUC domain.ModerationUseCase, wordFilterUC domain.WordFilterUseCase) *ModerationHandler {
	return &ModerationHandler{
		moderationUC: moderationUC,
		wordFilterUC: wordFilterUC,
	}
}

// CreateRuleRequest represents a request to create a moderation rule
type CreateRuleRequest struct {
	Pattern     string                      `json:"pattern" validate:"required"`
	Description string                      `json:"description" validate:"required"`
	Action      domain.ModerationActionType `json:"action" validate:"required,oneof=warn delete mute ban"`
	Level       domain.ModerationLevel      `json:"level" validate:"required,oneof=low medium high"`
	Duration    *time.Duration              `json:"duration,omitempty"`
}

// CreateRule creates a new moderation rule
func (h *ModerationHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
	var req CreateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rule := &domain.ModerationRule{
		ID:          uuid.New(),
		Pattern:     req.Pattern,
		Description: req.Description,
		Action:      req.Action,
		Level:       req.Level,
		Duration:    req.Duration,
		CreatedAt:   time.Now(),
		IsActive:    true,
	}

	if err := h.moderationUC.CreateRule(r.Context(), rule); err != nil {
		http.Error(w, "Failed to create rule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rule)
}

// ListRules lists all moderation rules
func (h *ModerationHandler) ListRules(w http.ResponseWriter, r *http.Request) {
	serverID, err := uuid.Parse(mux.Vars(r)["server_id"])
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}

	rules, err := h.moderationUC.ListRules(r.Context(), serverID)
	if err != nil {
		http.Error(w, "Failed to list rules", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

// GetRule retrieves a specific moderation rule
func (h *ModerationHandler) GetRule(w http.ResponseWriter, r *http.Request) {
	ruleID, err := uuid.Parse(mux.Vars(r)["rule_id"])
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	rule, err := h.moderationUC.GetRule(r.Context(), ruleID)
	if err != nil {
		http.Error(w, "Failed to get rule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rule)
}

// UpdateRuleRequest represents a request to update a moderation rule
type UpdateRuleRequest struct {
	Pattern     string                      `json:"pattern,omitempty"`
	Description string                      `json:"description,omitempty"`
	Action      domain.ModerationActionType `json:"action,omitempty"`
	Level       domain.ModerationLevel      `json:"level,omitempty"`
	Duration    *time.Duration              `json:"duration,omitempty"`
	IsActive    *bool                       `json:"is_active,omitempty"`
}

// UpdateRule updates an existing moderation rule
func (h *ModerationHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	ruleID, err := uuid.Parse(mux.Vars(r)["rule_id"])
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	var req UpdateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rule, err := h.moderationUC.GetRule(r.Context(), ruleID)
	if err != nil {
		http.Error(w, "Rule not found", http.StatusNotFound)
		return
	}

	if req.Pattern != "" {
		rule.Pattern = req.Pattern
	}
	if req.Description != "" {
		rule.Description = req.Description
	}
	if req.Action != "" {
		rule.Action = req.Action
	}
	if req.Level != "" {
		rule.Level = req.Level
	}
	if req.Duration != nil {
		rule.Duration = req.Duration
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}

	rule.UpdatedAt = time.Now()

	if err := h.moderationUC.UpdateRule(r.Context(), rule); err != nil {
		http.Error(w, "Failed to update rule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rule)
}

// DeleteRule deletes a moderation rule
func (h *ModerationHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	ruleID, err := uuid.Parse(mux.Vars(r)["rule_id"])
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	if err := h.moderationUC.DeleteRule(r.Context(), ruleID); err != nil {
		http.Error(w, "Failed to delete rule", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetViolations retrieves moderation violations
func (h *ModerationHandler) GetViolations(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	channelID := r.URL.Query().Get("channel_id")

	var violations []*domain.ModerationViolation
	var err error

	if userID != "" {
		uid, err := uuid.Parse(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		violations, err = h.moderationUC.GetUserViolations(r.Context(), uid)
	} else if channelID != "" {
		cid, err := uuid.Parse(channelID)
		if err != nil {
			http.Error(w, "Invalid channel ID", http.StatusBadRequest)
			return
		}
		violations, err = h.moderationUC.GetChannelViolations(r.Context(), cid)
	} else {
		http.Error(w, "Must specify user_id or channel_id", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to get violations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(violations)
}

// GetModerationLogs retrieves moderation logs
func (h *ModerationHandler) GetModerationLogs(w http.ResponseWriter, r *http.Request) {
	filter := &domain.ModerationLogFilter{
		Limit:  100,
		Offset: 0,
	}

	if userID := r.URL.Query().Get("user_id"); userID != "" {
		uid, err := uuid.Parse(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		filter.UserID = &uid
	}

	if moderatorID := r.URL.Query().Get("moderator_id"); moderatorID != "" {
		mid, err := uuid.Parse(moderatorID)
		if err != nil {
			http.Error(w, "Invalid moderator ID", http.StatusBadRequest)
			return
		}
		filter.ModeratorID = &mid
	}

	if action := r.URL.Query().Get("action"); action != "" {
		actionType := domain.ModerationActionType(action)
		filter.Action = &actionType
	}

	logs, err := h.moderationUC.GetModerationLogs(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to get moderation logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
