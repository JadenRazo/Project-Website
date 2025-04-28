package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrInvalidReadReceipt is returned when a read receipt is invalid
var ErrInvalidReadReceipt = errors.New("invalid read receipt")

// ErrReadReceiptNotFound is returned when a read receipt is not found
var ErrReadReceiptNotFound = errors.New("read receipt not found")

// NewReadReceipt creates a new read receipt
func NewReadReceipt(channelID, userID, messageID uuid.UUID) (*ReadReceipt, error) {
	if channelID == uuid.Nil {
		return nil, errors.New("channel ID cannot be empty")
	}

	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}

	if messageID == uuid.Nil {
		return nil, errors.New("message ID cannot be empty")
	}

	return &ReadReceipt{
		ID:        uuid.New(),
		ChannelID: channelID,
		UserID:    userID,
		MessageID: messageID,
		ReadAt:    time.Now(),
	}, nil
}

// NewChannelReadStatus creates a new channel read status
func NewChannelReadStatus(channelID, userID, lastReadID uuid.UUID, unreadCount, mentionCount int) *ChannelReadStatus {
	return &ChannelReadStatus{
		ChannelID:    channelID,
		UserID:       userID,
		LastReadID:   lastReadID,
		LastReadAt:   time.Now(),
		UnreadCount:  unreadCount,
		MentionCount: mentionCount,
	}
}

// HasUnread checks if the channel has unread messages
func (c *ChannelReadStatus) HasUnread() bool {
	return c.UnreadCount > 0
}

// HasMentions checks if the channel has unread mentions
func (c *ChannelReadStatus) HasMentions() bool {
	return c.MentionCount > 0
}

// UpdateReadStatus updates the channel read status with a new message
func (c *ChannelReadStatus) UpdateReadStatus(messageID uuid.UUID) {
	c.LastReadID = messageID
	c.LastReadAt = time.Now()
	c.UnreadCount = 0
	c.MentionCount = 0
}

// IncrementUnread increments the unread count
func (c *ChannelReadStatus) IncrementUnread() {
	c.UnreadCount++
}

// IncrementMention increments the mention count
func (c *ChannelReadStatus) IncrementMention() {
	c.MentionCount++
}

// CalculateReadStatus calculates the read status for a list of messages
func CalculateReadStatus(messages []Message, receipts []ReadReceipt, userID uuid.UUID) ChannelReadStatus {
	if len(messages) == 0 {
		return ChannelReadStatus{
			UserID:       userID,
			UnreadCount:  0,
			MentionCount: 0,
		}
	}

	// Find the latest read receipt for the user
	var latestReceipt *ReadReceipt
	for i := range receipts {
		if receipts[i].UserID == userID {
			if latestReceipt == nil || receipts[i].ReadAt.After(latestReceipt.ReadAt) {
				latestReceipt = &receipts[i]
			}
		}
	}

	// If no read receipt is found, all messages are unread
	if latestReceipt == nil {
		unreadCount := len(messages)
		mentionCount := countMentions(messages, userID)

		return ChannelReadStatus{
			ChannelID:    messages[0].ChannelID,
			UserID:       userID,
			UnreadCount:  unreadCount,
			MentionCount: mentionCount,
		}
	}

	// Count unread messages and mentions
	unreadCount := 0
	mentionCount := 0
	lastReadID := latestReceipt.MessageID
	lastReadAt := latestReceipt.ReadAt

	for _, msg := range messages {
		if msg.CreatedAt.After(lastReadAt) {
			unreadCount++

			// Check for mentions (this is simplified, actual implementation would need to parse the message content)
			if containsMention(msg, userID) {
				mentionCount++
			}
		}
	}

	return ChannelReadStatus{
		ChannelID:    messages[0].ChannelID,
		UserID:       userID,
		LastReadID:   lastReadID,
		LastReadAt:   lastReadAt,
		UnreadCount:  unreadCount,
		MentionCount: mentionCount,
	}
}

// countMentions counts the number of mentions for a user in a list of messages
func countMentions(messages []Message, userID uuid.UUID) int {
	count := 0

	for _, msg := range messages {
		if containsMention(msg, userID) {
			count++
		}
	}

	return count
}

// containsMention checks if a message mentions a specific user
// This is a simplified implementation, in a real app you would parse the message content
func containsMention(message Message, userID uuid.UUID) bool {
	// In a real implementation, you would parse the message content for @mentions
	// or use a precomputed mentions table
	return false
}

// GetChannelReadCount gets the number of users who have read a message
func GetChannelReadCount(receipts []ReadReceipt, messageID uuid.UUID) int {
	count := 0

	for _, receipt := range receipts {
		if receipt.MessageID == messageID {
			count++
		}
	}

	return count
}

// GetReadUserIDs gets the IDs of users who have read a message
func GetReadUserIDs(receipts []ReadReceipt, messageID uuid.UUID) []uuid.UUID {
	userIDs := make([]uuid.UUID, 0)

	for _, receipt := range receipts {
		if receipt.MessageID == messageID {
			userIDs = append(userIDs, receipt.UserID)
		}
	}

	return userIDs
}
