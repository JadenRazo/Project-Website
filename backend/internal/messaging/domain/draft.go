package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrDraftNotFound is returned when a draft cannot be found
var ErrDraftNotFound = errors.New("draft not found")

// ErrInvalidDraft is returned when draft data is invalid
var ErrInvalidDraft = errors.New("invalid draft data")

// Draft represents a message draft
type Draft struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	UserID      uuid.UUID    `json:"user_id" db:"user_id"`
	ChannelID   uuid.UUID    `json:"channel_id" db:"channel_id"`
	ReplyToID   *uuid.UUID   `json:"reply_to_id,omitempty" db:"reply_to_id"`
	Content     string       `json:"content" db:"content"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	Attachments []Attachment `json:"attachments,omitempty" db:"-"`
}

// NewDraft creates a new message draft
func NewDraft(userID, channelID uuid.UUID, content string) *Draft {
	now := time.Now()
	return &Draft{
		ID:          uuid.New(),
		UserID:      userID,
		ChannelID:   channelID,
		Content:     content,
		CreatedAt:   now,
		UpdatedAt:   now,
		Attachments: []Attachment{},
	}
}

// NewReplyDraft creates a new reply draft
func NewReplyDraft(userID, channelID, replyToID uuid.UUID, content string) *Draft {
	draft := NewDraft(userID, channelID, content)
	draft.ReplyToID = &replyToID
	return draft
}

// Update updates the draft content
func (d *Draft) Update(content string) {
	d.Content = content
	d.UpdatedAt = time.Now()
}

// AddAttachment adds an attachment to a draft
func (d *Draft) AddAttachment(attachment *Attachment) {
	if d.Attachments == nil {
		d.Attachments = []Attachment{}
	}
	d.Attachments = append(d.Attachments, *attachment)
	d.UpdatedAt = time.Now()
}

// RemoveAttachment removes an attachment from a draft
func (d *Draft) RemoveAttachment(attachmentID uuid.UUID) {
	if d.Attachments == nil {
		return
	}

	updatedAttachments := make([]Attachment, 0, len(d.Attachments))
	for _, att := range d.Attachments {
		if att.ID != attachmentID {
			updatedAttachments = append(updatedAttachments, att)
		}
	}

	d.Attachments = updatedAttachments
	d.UpdatedAt = time.Now()
}

// ConvertToMessage converts a draft to a message
func (d *Draft) ConvertToMessage(messageType MessageType) (*Message, error) {
	if d.Content == "" && len(d.Attachments) == 0 {
		return nil, errors.New("draft has no content or attachments")
	}

	if messageType == "" {
		messageType = MessageTypeText
	}

	message := &Message{
		ID:          uuid.New(),
		ChannelID:   d.ChannelID,
		UserID:      d.UserID,
		Content:     d.Content,
		Type:        messageType,
		ReplyToID:   d.ReplyToID,
		CreatedAt:   time.Now(),
		IsDeleted:   false,
		Attachments: d.Attachments,
		Reactions:   []Reaction{},
	}

	return message, nil
}

// DraftRepository defines the interface for draft data access
type DraftRepository interface {
	// Create creates a new draft
	Create(ctx context.Context, draft *Draft) error

	// Get retrieves a draft by ID
	Get(ctx context.Context, id uuid.UUID) (*Draft, error)

	// Update updates a draft
	Update(ctx context.Context, draft *Draft) error

	// Delete deletes a draft
	Delete(ctx context.Context, id uuid.UUID) error

	// GetUserDrafts retrieves drafts for a user
	GetUserDrafts(ctx context.Context, userID uuid.UUID) ([]*Draft, error)

	// GetChannelDraft retrieves a draft for a user in a channel
	GetChannelDraft(ctx context.Context, userID, channelID uuid.UUID) (*Draft, error)

	// GetReplyDraft retrieves a reply draft for a user
	GetReplyDraft(ctx context.Context, userID, replyToID uuid.UUID) (*Draft, error)

	// AddAttachment adds an attachment to a draft
	AddAttachment(ctx context.Context, draftID uuid.UUID, attachment *Attachment) error

	// RemoveAttachment removes an attachment from a draft
	RemoveAttachment(ctx context.Context, draftID, attachmentID uuid.UUID) error
}

// DraftUseCase defines the business logic for draft operations
type DraftUseCase interface {
	// CreateDraft creates a new draft
	CreateDraft(ctx context.Context, userID, channelID uuid.UUID, content string) (*Draft, error)

	// CreateReplyDraft creates a new reply draft
	CreateReplyDraft(ctx context.Context, userID, channelID, replyToID uuid.UUID, content string) (*Draft, error)

	// GetDraft retrieves a draft by ID
	GetDraft(ctx context.Context, id string) (*Draft, error)

	// UpdateDraft updates a draft
	UpdateDraft(ctx context.Context, id string, content string) (*Draft, error)

	// DeleteDraft deletes a draft
	DeleteDraft(ctx context.Context, id string) error

	// GetUserDrafts retrieves drafts for a user
	GetUserDrafts(ctx context.Context, userID string) ([]*Draft, error)

	// GetChannelDraft retrieves a draft for a user in a channel
	GetChannelDraft(ctx context.Context, userID, channelID string) (*Draft, error)

	// GetReplyDraft retrieves a reply draft for a user
	GetReplyDraft(ctx context.Context, userID, replyToID string) (*Draft, error)

	// AddAttachment adds an attachment to a draft
	AddAttachment(ctx context.Context, draftID string, fileName string, fileSize int64, fileType, fileURL string) (*Draft, error)

	// RemoveAttachment removes an attachment from a draft
	RemoveAttachment(ctx context.Context, draftID, attachmentID string) (*Draft, error)

	// ConvertToMessage converts a draft to a message and sends it
	ConvertToMessage(ctx context.Context, draftID string) (*Message, error)
}
