package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrThreadNotFound is returned when a thread cannot be found
var ErrThreadNotFound = errors.New("thread not found")

// ErrInvalidThread is returned when thread data is invalid
var ErrInvalidThread = errors.New("invalid thread data")

// ThreadParticipant represents a participant in a message thread
type ThreadParticipant struct {
	ThreadID    uuid.UUID `json:"thread_id" db:"thread_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	LastReadID  uuid.UUID `json:"last_read_id,omitempty" db:"last_read_id"`
	LastReadAt  time.Time `json:"last_read_at" db:"last_read_at"`
	JoinedAt    time.Time `json:"joined_at" db:"joined_at"`
	IsFollowing bool      `json:"is_following" db:"is_following"`
}

// ThreadSummary represents a summary of a thread
type ThreadSummary struct {
	ParentMessage      *Message    `json:"parent_message"`
	ReplyCount         int         `json:"reply_count"`
	ParticipantCount   int         `json:"participant_count"`
	LatestReplyAt      time.Time   `json:"latest_reply_at"`
	LatestReplyID      uuid.UUID   `json:"latest_reply_id"`
	UnreadCount        int         `json:"unread_count"`
	HasUnread          bool        `json:"has_unread"`
	RecentParticipants []uuid.UUID `json:"recent_participants"`
	IsFollowing        bool        `json:"is_following"`
}

// NewThreadParticipant creates a new thread participant
func NewThreadParticipant(threadID, userID uuid.UUID) *ThreadParticipant {
	now := time.Now()
	return &ThreadParticipant{
		ThreadID:    threadID,
		UserID:      userID,
		LastReadAt:  now,
		JoinedAt:    now,
		IsFollowing: true,
	}
}

// MarkAsRead marks the thread as read up to a specific message
func (p *ThreadParticipant) MarkAsRead(messageID uuid.UUID) {
	p.LastReadID = messageID
	p.LastReadAt = time.Now()
}

// SetFollowing sets whether the participant is following the thread
func (p *ThreadParticipant) SetFollowing(following bool) {
	p.IsFollowing = following
}

// ThreadRepository defines the interface for thread data access
type ThreadRepository interface {
	// GetThreadMessages retrieves messages in a thread
	GetThreadMessages(ctx context.Context, parentID uuid.UUID, pagination Pagination) ([]*Message, int, error)

	// GetThreadSummary retrieves a summary of a thread
	GetThreadSummary(ctx context.Context, parentID uuid.UUID, userID uuid.UUID) (*ThreadSummary, error)

	// GetThreadParticipant retrieves a thread participant
	GetThreadParticipant(ctx context.Context, threadID, userID uuid.UUID) (*ThreadParticipant, error)

	// AddThreadParticipant adds a participant to a thread
	AddThreadParticipant(ctx context.Context, participant *ThreadParticipant) error

	// UpdateThreadParticipant updates a thread participant
	UpdateThreadParticipant(ctx context.Context, participant *ThreadParticipant) error

	// MarkThreadAsRead marks a thread as read for a user
	MarkThreadAsRead(ctx context.Context, threadID, userID, messageID uuid.UUID) error

	// SetThreadFollowing sets whether a user is following a thread
	SetThreadFollowing(ctx context.Context, threadID, userID uuid.UUID, following bool) error

	// GetUserThreads retrieves threads that a user is participating in
	GetUserThreads(ctx context.Context, userID uuid.UUID, pagination Pagination) ([]*ThreadSummary, int, error)
}

// ThreadUseCase defines the business logic for thread operations
type ThreadUseCase interface {
	// GetThreadMessages retrieves messages in a thread
	GetThreadMessages(ctx context.Context, parentID string, page, pageSize int) ([]*Message, int, error)

	// GetThreadSummary retrieves a summary of a thread
	GetThreadSummary(ctx context.Context, parentID string, userID string) (*ThreadSummary, error)

	// ReplyToThread creates a reply in a thread
	ReplyToThread(ctx context.Context, parentID string, userID string, content string) (*Message, error)

	// MarkThreadAsRead marks a thread as read for a user
	MarkThreadAsRead(ctx context.Context, threadID string, userID string) error

	// SetThreadFollowing sets whether a user is following a thread
	SetThreadFollowing(ctx context.Context, threadID string, userID string, following bool) error

	// GetUserThreads retrieves threads that a user is participating in
	GetUserThreads(ctx context.Context, userID string, page, pageSize int) ([]*ThreadSummary, int, error)
}

// IsThreadReply checks if a message is a thread reply
func IsThreadReply(message *Message) bool {
	return message != nil && message.ReplyToID != nil
}

// GetThreadID returns the thread ID (parent message ID) for a message
func GetThreadID(message *Message) uuid.UUID {
	if message == nil || message.ReplyToID == nil {
		return message.ID // This is a parent message, so its own ID is the thread ID
	}
	return *message.ReplyToID // This is a reply, so the parent message ID is the thread ID
}

// OrganizeThreadMessages groups messages into threads
func OrganizeThreadMessages(messages []*Message) map[uuid.UUID][]*Message {
	threads := make(map[uuid.UUID][]*Message)

	for _, msg := range messages {
		threadID := GetThreadID(msg)

		if _, exists := threads[threadID]; !exists {
			threads[threadID] = make([]*Message, 0)
		}

		threads[threadID] = append(threads[threadID], msg)
	}

	return threads
}
