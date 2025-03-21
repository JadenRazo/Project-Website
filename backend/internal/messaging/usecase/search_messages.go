package usecase

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// SearchMessagesUseCase handles message searching with advanced filters
type SearchMessagesUseCase struct {
	messageRepo repository.MessagingMessageRepository
}

// NewSearchMessagesUseCase creates a new search messages usecase
func NewSearchMessagesUseCase(messageRepo repository.MessagingMessageRepository) *SearchMessagesUseCase {
	return &SearchMessagesUseCase{
		messageRepo: messageRepo,
	}
}

// SearchMessagesInput represents the input for searching messages
type SearchMessagesInput struct {
	ChannelID      *uint
	UserID         *uint
	Query          string
	StartTime      *time.Time
	EndTime        *time.Time
	HasAttachments bool
	HasMentions    bool
	IsPinned       bool
	IsNSFW         bool
	IsSpoiler      bool
	Limit          int
	Offset         int
}

// SearchResult represents a search result with metadata
type SearchResult struct {
	Messages []domain.MessagingMessage
	Total    int
	HasMore  bool
}

// Execute searches messages with advanced filtering
func (uc *SearchMessagesUseCase) Execute(ctx context.Context, input SearchMessagesInput) (*SearchResult, error) {
	// Build search filters
	filters := repository.MessageSearchFilters{
		ChannelID:      input.ChannelID,
		UserID:         input.UserID,
		Query:          input.Query,
		StartTime:      input.StartTime,
		EndTime:        input.EndTime,
		HasAttachments: input.HasAttachments,
		HasMentions:    input.HasMentions,
		IsPinned:       input.IsPinned,
		IsNSFW:         input.IsNSFW,
		IsSpoiler:      input.IsSpoiler,
		Limit:          input.Limit,
		Offset:         input.Offset,
	}

	// Search messages
	messages, total, err := uc.messageRepo.SearchMessages(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Check if there are more results
	hasMore := total > (input.Offset + input.Limit)

	return &SearchResult{
		Messages: messages,
		Total:    total,
		HasMore:  hasMore,
	}, nil
}

// SearchInThreadInput represents the input for searching messages in a thread
type SearchInThreadInput struct {
	ThreadID       uint
	Query          string
	StartTime      *time.Time
	EndTime        *time.Time
	HasAttachments bool
	HasMentions    bool
	Limit          int
	Offset         int
}

// SearchInThread searches messages within a specific thread
func (uc *SearchMessagesUseCase) SearchInThread(ctx context.Context, input SearchInThreadInput) (*SearchResult, error) {
	// Build thread search filters
	filters := repository.MessageSearchFilters{
		ThreadID:       &input.ThreadID,
		Query:          input.Query,
		StartTime:      input.StartTime,
		EndTime:        input.EndTime,
		HasAttachments: input.HasAttachments,
		HasMentions:    input.HasMentions,
		Limit:          input.Limit,
		Offset:         input.Offset,
	}

	// Search messages in thread
	messages, total, err := uc.messageRepo.SearchMessages(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Check if there are more results
	hasMore := total > (input.Offset + input.Limit)

	return &SearchResult{
		Messages: messages,
		Total:    total,
		HasMore:  hasMore,
	}, nil
}
