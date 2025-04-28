package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrMessageNotFound is returned when a message cannot be found
var ErrMessageNotFound = errors.New("message not found")

// ErrInvalidMessage is returned when message data is invalid
var ErrInvalidMessage = errors.New("invalid message data")

// ErrMessageDeleted is returned when attempting to access a deleted message
var ErrMessageDeleted = errors.New("message has been deleted")

// NewMessage creates a new message
func NewMessage(channelID, userID uuid.UUID, content string, messageType MessageType) (*Message, error) {
	if channelID == uuid.Nil {
		return nil, errors.New("channel ID cannot be empty")
	}

	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}

	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	if messageType == "" {
		messageType = MessageTypeText
	}

	now := time.Now()

	return &Message{
		ID:          uuid.New(),
		ChannelID:   channelID,
		UserID:      userID,
		Content:     content,
		Type:        messageType,
		CreatedAt:   now,
		IsDeleted:   false,
		Attachments: []Attachment{},
		Reactions:   []Reaction{},
	}, nil
}

// NewSystemMessage creates a new system message
func NewSystemMessage(channelID uuid.UUID, content string) (*Message, error) {
	if channelID == uuid.Nil {
		return nil, errors.New("channel ID cannot be empty")
	}

	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	now := time.Now()

	// Use a zero UUID for system messages
	systemUserID := uuid.Nil

	return &Message{
		ID:          uuid.New(),
		ChannelID:   channelID,
		UserID:      systemUserID,
		Content:     content,
		Type:        MessageTypeSystem,
		CreatedAt:   now,
		IsDeleted:   false,
		Attachments: []Attachment{},
		Reactions:   []Reaction{},
	}, nil
}

// NewReplyMessage creates a new reply message
func NewReplyMessage(channelID, userID, replyToID uuid.UUID, content string) (*Message, error) {
	if channelID == uuid.Nil {
		return nil, errors.New("channel ID cannot be empty")
	}

	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}

	if replyToID == uuid.Nil {
		return nil, errors.New("reply to ID cannot be empty")
	}

	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	now := time.Now()

	return &Message{
		ID:          uuid.New(),
		ChannelID:   channelID,
		UserID:      userID,
		Content:     content,
		Type:        MessageTypeText,
		ReplyToID:   &replyToID,
		CreatedAt:   now,
		IsDeleted:   false,
		Attachments: []Attachment{},
		Reactions:   []Reaction{},
	}, nil
}

// NewAttachment creates a new file attachment
func NewAttachment(messageID, userID uuid.UUID, fileName string, fileSize int64, fileType, fileURL, thumbnailURL string) (*Attachment, error) {
	if messageID == uuid.Nil {
		return nil, errors.New("message ID cannot be empty")
	}

	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}

	if fileName == "" {
		return nil, errors.New("file name cannot be empty")
	}

	if fileURL == "" {
		return nil, errors.New("file URL cannot be empty")
	}

	return &Attachment{
		ID:           uuid.New(),
		MessageID:    messageID,
		UserID:       userID,
		FileName:     fileName,
		FileSize:     fileSize,
		FileType:     fileType,
		FileURL:      fileURL,
		ThumbnailURL: thumbnailURL,
		CreatedAt:    time.Now(),
	}, nil
}

// AddAttachment adds an attachment to a message
func (m *Message) AddAttachment(attachment *Attachment) {
	if m.Attachments == nil {
		m.Attachments = []Attachment{}
	}

	m.Attachments = append(m.Attachments, *attachment)
}

// AddReaction adds a reaction to a message
func (m *Message) AddReaction(reaction *Reaction) {
	if m.Reactions == nil {
		m.Reactions = []Reaction{}
	}

	// Check if the reaction already exists from this user
	for i, r := range m.Reactions {
		if r.Emoji == reaction.Emoji && r.UserID == reaction.UserID {
			// Update existing reaction
			m.Reactions[i] = *reaction
			return
		}
	}

	// Add new reaction
	m.Reactions = append(m.Reactions, *reaction)
}

// Edit updates the message content and marks it as edited
func (m *Message) Edit(content string) error {
	if content == "" {
		return errors.New("message content cannot be empty")
	}

	if m.IsDeleted {
		return ErrMessageDeleted
	}

	// Don't allow editing system messages
	if m.Type == MessageTypeSystem {
		return errors.New("cannot edit system messages")
	}

	m.Content = content
	now := time.Now()
	m.EditedAt = &now

	return nil
}

// SoftDelete marks a message as deleted without removing it from the database
func (m *Message) SoftDelete() {
	m.IsDeleted = true
	m.Content = "This message has been deleted"
	m.Attachments = []Attachment{}
	m.Reactions = []Reaction{}
}

// IsReply returns true if the message is a reply to another message
func (m *Message) IsReply() bool {
	return m.ReplyToID != nil
}

// IsThread returns true if the message is part of a thread
func (m *Message) IsThread() bool {
	return m.ThreadCount > 0
}

// HasAttachments returns true if the message has attachments
func (m *Message) HasAttachments() bool {
	return len(m.Attachments) > 0
}

// HasReactions returns true if the message has reactions
func (m *Message) HasReactions() bool {
	return len(m.Reactions) > 0
}
