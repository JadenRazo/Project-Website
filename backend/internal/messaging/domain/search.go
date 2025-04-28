package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SearchOperator represents a search operator
type SearchOperator string

const (
	// SearchOperatorAnd represents the AND operator
	SearchOperatorAnd SearchOperator = "AND"

	// SearchOperatorOr represents the OR operator
	SearchOperatorOr SearchOperator = "OR"
)

// SearchFilter defines the criteria for searching messages
type SearchFilter struct {
	UserID        *uuid.UUID    `json:"user_id,omitempty"`
	ChannelID     *uuid.UUID    `json:"channel_id,omitempty"`
	Types         []MessageType `json:"types,omitempty"`
	HasAttachment *bool         `json:"has_attachment,omitempty"`
	FromDate      *time.Time    `json:"from_date,omitempty"`
	ToDate        *time.Time    `json:"to_date,omitempty"`
	Query         string        `json:"query,omitempty"`
}

// SearchResult represents a search result
type SearchResult struct {
	Message    *Message `json:"message"`
	Channel    *Channel `json:"channel"`
	User       *User    `json:"user"`
	Score      float64  `json:"score"`
	Highlights []string `json:"highlights"`
}

// SearchOptions represents options to customize search results
type SearchOptions struct {
	Pagination      Pagination `json:"pagination"`
	IncludeDeleted  bool       `json:"include_deleted,omitempty"`
	HighlightFields []string   `json:"highlight_fields,omitempty"`
	SortField       string     `json:"sort_field,omitempty"`
	SortOrder       string     `json:"sort_order,omitempty"`
}

// SearchRepository defines the interface for search functionality
type SearchRepository interface {
	// SearchMessages searches for messages based on criteria
	SearchMessages(ctx context.Context, filter SearchFilter, pagination Pagination) ([]*SearchResult, int, error)

	// SearchUsers searches for users
	SearchUsers(ctx context.Context, query string, pagination Pagination) ([]*User, int, error)

	// SearchChannels searches for channels
	SearchChannels(ctx context.Context, query string, pagination Pagination) ([]*Channel, int, error)

	// GetSearchSuggestions gets search suggestions
	GetSearchSuggestions(ctx context.Context, query string) ([]string, error)
}

// MessageWithHighlights represents a message with highlighted text snippets
type MessageWithHighlights struct {
	Message       *Message          `json:"message"`
	Highlights    map[string]string `json:"highlights,omitempty"`
	ContextBefore []*Message        `json:"context_before,omitempty"`
	ContextAfter  []*Message        `json:"context_after,omitempty"`
}

// SearchUseCase defines the interface for search business logic
type SearchUseCase interface {
	// SearchMessages searches for messages based on criteria
	SearchMessages(ctx context.Context, filter SearchFilter, pagination Pagination) ([]*SearchResult, int, error)

	// SearchUsers searches for users
	SearchUsers(ctx context.Context, query string, pagination Pagination) ([]*User, int, error)

	// SearchChannels searches for channels
	SearchChannels(ctx context.Context, query string, pagination Pagination) ([]*Channel, int, error)

	// GetSearchSuggestions gets search suggestions
	GetSearchSuggestions(ctx context.Context, query string) ([]string, error)
}

// NewSearchFilter creates a new search filter with default values
func NewSearchFilter(query string) SearchFilter {
	return SearchFilter{
		Query: query,
	}
}

// WithChannels adds channel filtering to a search filter
func (f SearchFilter) WithChannels(channelIDs ...uuid.UUID) SearchFilter {
	if len(channelIDs) > 0 {
		f.ChannelID = &channelIDs[0] // Using first channel ID for now
	}
	return f
}

// WithUsers adds user filtering to a search filter
func (f SearchFilter) WithUsers(userIDs ...uuid.UUID) SearchFilter {
	if len(userIDs) > 0 {
		f.UserID = &userIDs[0] // Using first user ID for now
	}
	return f
}

// WithDateRange adds date range filtering to a search filter
func (f SearchFilter) WithDateRange(from, to time.Time) SearchFilter {
	f.FromDate = &from
	f.ToDate = &to
	return f
}

// WithFiles filters messages that have file attachments
func (f SearchFilter) WithFiles(hasFiles bool) SearchFilter {
	f.HasAttachment = &hasFiles
	return f
}

// WithMessageType filters messages by type
func (f SearchFilter) WithMessageType(messageType MessageType) SearchFilter {
	f.Types = []MessageType{messageType}
	return f
}

// NewSearchOptions creates a new search options with default values
func NewSearchOptions() SearchOptions {
	return SearchOptions{
		Pagination:      NewDefaultPagination(),
		IncludeDeleted:  false,
		HighlightFields: []string{"content"},
		SortField:       "created_at",
		SortOrder:       "desc",
	}
}

// WithPagination sets pagination options for search
func (o SearchOptions) WithPagination(page, pageSize int) SearchOptions {
	o.Pagination = Pagination{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    o.SortField,
		SortOrder: o.SortOrder,
	}
	return o
}

// WithHighlightFields sets fields to highlight in search results
func (o SearchOptions) WithHighlightFields(fields ...string) SearchOptions {
	o.HighlightFields = fields
	return o
}

// WithSort sets sorting options for search results
func (o SearchOptions) WithSort(field, order string) SearchOptions {
	o.SortField = field
	o.SortOrder = order
	o.Pagination.SortBy = field
	o.Pagination.SortOrder = order
	return o
}

// WithIncludeDeleted includes deleted messages in search results
func (o SearchOptions) WithIncludeDeleted(include bool) SearchOptions {
	o.IncludeDeleted = include
	return o
}
