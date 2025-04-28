package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrInvalidReaction is returned when a reaction is invalid
var ErrInvalidReaction = errors.New("invalid reaction")

// ErrReactionNotFound is returned when a reaction is not found
var ErrReactionNotFound = errors.New("reaction not found")

// NewReaction creates a new reaction
func NewReaction(messageID, userID uuid.UUID, emoji string) (*Reaction, error) {
	if messageID == uuid.Nil {
		return nil, errors.New("message ID cannot be empty")
	}

	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}

	if emoji == "" {
		return nil, errors.New("emoji cannot be empty")
	}

	return &Reaction{
		ID:        uuid.New(),
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
		CreatedAt: time.Now(),
	}, nil
}

// ValidateEmoji checks if the emoji is valid
// This implementation is simple, but in production you might want to use a proper emoji validation library
func ValidateEmoji(emoji string) bool {
	if emoji == "" {
		return false
	}

	// Simplified validation: just check for reasonable length
	// A full validation would use unicode categories
	return len(emoji) >= 1 && len(emoji) <= 8
}

// CountReactions groups reactions by emoji and counts them
func CountReactions(reactions []Reaction, currentUserID uuid.UUID) []ReactionCount {
	counts := make(map[string]*ReactionCount)

	for _, reaction := range reactions {
		if count, exists := counts[reaction.Emoji]; exists {
			count.Count++
			if reaction.UserID == currentUserID {
				count.HasReacted = true
			}
		} else {
			hasReacted := reaction.UserID == currentUserID
			counts[reaction.Emoji] = &ReactionCount{
				Emoji:      reaction.Emoji,
				Count:      1,
				Users:      []string{},
				HasReacted: hasReacted,
			}
		}
	}

	result := make([]ReactionCount, 0, len(counts))
	for _, count := range counts {
		result = append(result, *count)
	}

	return result
}

// RemoveReaction removes a reaction from a list of reactions
func RemoveReaction(reactions []Reaction, userID uuid.UUID, emoji string) []Reaction {
	result := make([]Reaction, 0, len(reactions))

	for _, reaction := range reactions {
		if reaction.UserID != userID || reaction.Emoji != emoji {
			result = append(result, reaction)
		}
	}

	return result
}

// HasUserReacted checks if a user has reacted with a specific emoji
func HasUserReacted(reactions []Reaction, userID uuid.UUID, emoji string) bool {
	for _, reaction := range reactions {
		if reaction.UserID == userID && reaction.Emoji == emoji {
			return true
		}
	}

	return false
}

// GetUserReactions gets all reactions from a specific user
func GetUserReactions(reactions []Reaction, userID uuid.UUID) []Reaction {
	result := make([]Reaction, 0)

	for _, reaction := range reactions {
		if reaction.UserID == userID {
			result = append(result, reaction)
		}
	}

	return result
}

// GetEmojiReactions gets all reactions with a specific emoji
func GetEmojiReactions(reactions []Reaction, emoji string) []Reaction {
	result := make([]Reaction, 0)

	for _, reaction := range reactions {
		if reaction.Emoji == emoji {
			result = append(result, reaction)
		}
	}

	return result
}
