package repository

import (
	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"gorm.io/gorm"
)

// InitRepositories initializes all repositories and their tables
func InitRepositories(db *gorm.DB) error {
	// Auto-migrate all domain models
	err := db.AutoMigrate(
		&domain.MessagingMessage{},
		&domain.MessagingChannel{},
		&domain.MessagingChannelMember{},
		&domain.MessagingPinnedMessage{},
		&domain.MessagingAttachment{},
		&domain.MessagingEmbed{},
		&domain.MessagingReaction{},
		&domain.MessagingReadReceipt{},
		&domain.User{},
	)
	if err != nil {
		return err
	}

	// Create indexes
	err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_messages_channel_id ON messaging_messages(channel_id);
		CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messaging_messages(sender_id);
		CREATE INDEX IF NOT EXISTS idx_messages_reply_to_id ON messaging_messages(reply_to_id);
		CREATE INDEX IF NOT EXISTS idx_messages_thread_id ON messaging_messages(thread_id);
		CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messaging_messages(created_at);
		CREATE INDEX IF NOT EXISTS idx_messages_is_pinned ON messaging_messages(is_pinned);
		CREATE INDEX IF NOT EXISTS idx_messages_is_nsfw ON messaging_messages(is_nsfw);
		CREATE INDEX IF NOT EXISTS idx_messages_is_spoiler ON messaging_messages(is_spoiler);

		CREATE INDEX IF NOT EXISTS idx_channels_type ON messaging_channels(type);
		CREATE INDEX IF NOT EXISTS idx_channels_is_archived ON messaging_channels(is_archived);
		CREATE INDEX IF NOT EXISTS idx_channels_is_private ON messaging_channels(is_private);
		CREATE INDEX IF NOT EXISTS idx_channels_is_read_only ON messaging_channels(is_read_only);

		CREATE INDEX IF NOT EXISTS idx_channel_members_channel_id ON messaging_channel_members(channel_id);
		CREATE INDEX IF NOT EXISTS idx_channel_members_user_id ON messaging_channel_members(user_id);
		CREATE INDEX IF NOT EXISTS idx_channel_members_role ON messaging_channel_members(role);
		CREATE INDEX IF NOT EXISTS idx_channel_members_is_muted ON messaging_channel_members(is_muted);
		CREATE INDEX IF NOT EXISTS idx_channel_members_is_banned ON messaging_channel_members(is_banned);

		CREATE INDEX IF NOT EXISTS idx_pinned_messages_channel_id ON messaging_pinned_messages(channel_id);
		CREATE INDEX IF NOT EXISTS idx_pinned_messages_message_id ON messaging_pinned_messages(message_id);
		CREATE INDEX IF NOT EXISTS idx_pinned_messages_pinned_by ON messaging_pinned_messages(pinned_by);

		CREATE INDEX IF NOT EXISTS idx_attachments_message_id ON messaging_attachments(message_id);
		CREATE INDEX IF NOT EXISTS idx_attachments_url ON messaging_attachments(url);
		CREATE INDEX IF NOT EXISTS idx_attachments_hash ON messaging_attachments(hash);
		CREATE INDEX IF NOT EXISTS idx_attachments_is_nsfw ON messaging_attachments(is_nsfw);
		CREATE INDEX IF NOT EXISTS idx_attachments_is_spoiler ON messaging_attachments(is_spoiler);

		CREATE INDEX IF NOT EXISTS idx_embeds_message_id ON messaging_embeds(message_id);
		CREATE INDEX IF NOT EXISTS idx_embeds_url ON messaging_embeds(url);
		CREATE INDEX IF NOT EXISTS idx_embeds_type ON messaging_embeds(type);
		CREATE INDEX IF NOT EXISTS idx_embeds_is_nsfw ON messaging_embeds(is_nsfw);
		CREATE INDEX IF NOT EXISTS idx_embeds_is_spoiler ON messaging_embeds(is_spoiler);

		CREATE INDEX IF NOT EXISTS idx_reactions_message_id ON messaging_reactions(message_id);
		CREATE INDEX IF NOT EXISTS idx_reactions_user_id ON messaging_reactions(user_id);
		CREATE INDEX IF NOT EXISTS idx_reactions_emoji_code ON messaging_reactions(emoji_code);

		CREATE INDEX IF NOT EXISTS idx_read_receipts_message_id ON messaging_read_receipts(message_id);
		CREATE INDEX IF NOT EXISTS idx_read_receipts_user_id ON messaging_read_receipts(user_id);
		CREATE INDEX IF NOT EXISTS idx_read_receipts_read_at ON messaging_read_receipts(read_at);
	`).Error
	if err != nil {
		return err
	}

	return nil
}
