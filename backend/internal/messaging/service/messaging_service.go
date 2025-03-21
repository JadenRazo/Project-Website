package service

import (
	"context"
	"errors"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/errors"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/events"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/interfaces"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// MessagingService defines the interface for messaging operations
type MessagingService interface {
	// Message operations
	CreateMessage(ctx context.Context, message *domain.MessagingMessage) error
	GetMessage(ctx context.Context, id uint) (*domain.MessagingMessage, error)
	UpdateMessage(ctx context.Context, message *domain.MessagingMessage) error
	DeleteMessage(ctx context.Context, id uint) error
	GetChannelMessages(ctx context.Context, channelID uint, limit, offset int) ([]*domain.MessagingMessage, error)
	GetThreadMessages(ctx context.Context, threadID uint, limit, offset int) ([]*domain.MessagingMessage, error)
	SearchMessages(ctx context.Context, filters interfaces.MessageSearchFilters) ([]*domain.MessagingMessage, error)

	// Channel operations
	CreateChannel(ctx context.Context, channel *domain.MessagingChannel) error
	GetChannel(ctx context.Context, id uint) (*domain.MessagingChannel, error)
	UpdateChannel(ctx context.Context, channel *domain.MessagingChannel) error
	DeleteChannel(ctx context.Context, id uint) error
	GetUserChannels(ctx context.Context, userID uint) ([]*domain.MessagingChannel, error)
	AddChannelMember(ctx context.Context, member *domain.MessagingChannelMember) error
	RemoveChannelMember(ctx context.Context, channelID, userID uint) error
	GetChannelMembers(ctx context.Context, channelID uint) ([]*domain.MessagingChannelMember, error)

	// Category operations
	CreateCategory(ctx context.Context, category *domain.MessagingCategory) error
	GetCategory(ctx context.Context, id uint) (*domain.MessagingCategory, error)
	UpdateCategory(ctx context.Context, category *domain.MessagingCategory) error
	DeleteCategory(ctx context.Context, id uint) error
	GetChannelCategories(ctx context.Context, channelID uint) ([]*domain.MessagingCategory, error)

	// User operations
	GetUser(ctx context.Context, id uint) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint) error

	// Role operations
	CreateRole(ctx context.Context, role *domain.MessagingRole) error
	GetRole(ctx context.Context, id uint) (*domain.MessagingRole, error)
	UpdateRole(ctx context.Context, role *domain.MessagingRole) error
	DeleteRole(ctx context.Context, id uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]*domain.MessagingRole, error)
	AssignRole(ctx context.Context, userID, roleID uint) error
	RemoveRole(ctx context.Context, userID, roleID uint) error

	// Reaction operations
	AddReaction(ctx context.Context, reaction *domain.MessagingReaction) error
	RemoveReaction(ctx context.Context, messageID, userID uint, emoji string) error
	GetMessageReactions(ctx context.Context, messageID uint) ([]*domain.MessagingReaction, error)

	// Read receipt operations
	MarkMessageAsRead(ctx context.Context, receipt *domain.MessagingReadReceipt) error
	GetUnreadCount(ctx context.Context, userID, channelID uint) (int64, error)
	GetReadReceipts(ctx context.Context, messageID uint) ([]*domain.MessagingReadReceipt, error)

	// Permission operations
	HasPermission(ctx context.Context, userID, channelID uint, permission domain.Permission) (bool, error)
	GetUserPermissions(ctx context.Context, userID, channelID uint) ([]domain.Permission, error)
}

// MessagingServiceImpl implements the MessagingService interface
type MessagingServiceImpl struct {
	messageRepo     interfaces.MessageRepository
	channelRepo     interfaces.ChannelRepository
	categoryRepo    interfaces.CategoryRepository
	userRepo        interfaces.UserRepository
	roleRepo        interfaces.RoleRepository
	reactionRepo    interfaces.ReactionRepository
	readReceiptRepo interfaces.ReadReceiptRepository
	attachmentRepo  interfaces.AttachmentRepository
	embedRepo       interfaces.EmbedRepository
	dispatcher      events.EventDispatcher
}

// NewMessagingService creates a new instance of MessagingService
func NewMessagingService(
	messageRepo interfaces.MessageRepository,
	channelRepo interfaces.ChannelRepository,
	categoryRepo interfaces.CategoryRepository,
	userRepo interfaces.UserRepository,
	roleRepo interfaces.RoleRepository,
	reactionRepo interfaces.ReactionRepository,
	readReceiptRepo interfaces.ReadReceiptRepository,
	attachmentRepo interfaces.AttachmentRepository,
	embedRepo interfaces.EmbedRepository,
	dispatcher events.EventDispatcher,
) interfaces.MessagingService {
	return &MessagingServiceImpl{
		messageRepo:     messageRepo,
		channelRepo:     channelRepo,
		categoryRepo:    categoryRepo,
		userRepo:        userRepo,
		roleRepo:        roleRepo,
		reactionRepo:    reactionRepo,
		readReceiptRepo: readReceiptRepo,
		attachmentRepo:  attachmentRepo,
		embedRepo:       embedRepo,
		dispatcher:      dispatcher,
	}
}

// CreateMessage creates a new message
func (s *MessagingServiceImpl) CreateMessage(ctx context.Context, message *domain.MessagingMessage) error {
	// Validate message
	if err := s.validateMessage(message); err != nil {
		return err
	}

	// Create attachments if any
	if len(message.Attachments) > 0 {
		for i := range message.Attachments {
			if err := s.attachmentRepo.Create(ctx, &message.Attachments[i]); err != nil {
				return err
			}
		}
	}

	// Create embeds if any
	if len(message.Embeds) > 0 {
		for i := range message.Embeds {
			if err := s.embedRepo.Create(ctx, &message.Embeds[i]); err != nil {
				return err
			}
		}
	}

	// Create message
	if err := s.messageRepo.Create(ctx, message); err != nil {
		return err
	}

	// Dispatch message created event
	event := &events.MessageCreatedEvent{
		MessageID: message.ID,
		ChannelID: message.ChannelID,
		UserID:    message.SenderID,
	}
	return s.dispatcher.Dispatch(ctx, event)
}

// GetMessage retrieves a message by ID
func (s *MessagingServiceImpl) GetMessage(ctx context.Context, id uint) (*domain.MessagingMessage, error) {
	return s.messageRepo.FindByID(ctx, id)
}

// UpdateMessage updates an existing message
func (s *MessagingServiceImpl) UpdateMessage(ctx context.Context, message *domain.MessagingMessage) error {
	// Validate message
	if err := s.validateMessage(message); err != nil {
		return err
	}

	// Check if message exists
	_, err := s.messageRepo.FindByID(ctx, message.ID)
	if err != nil {
		return err
	}

	// Update message
	if err := s.messageRepo.Update(ctx, message); err != nil {
		return err
	}

	// Dispatch message updated event
	event := &events.MessageUpdatedEvent{
		MessageID: message.ID,
		ChannelID: message.ChannelID,
		UserID:    message.SenderID,
	}
	return s.dispatcher.Dispatch(ctx, event)
}

// DeleteMessage deletes a message
func (s *MessagingServiceImpl) DeleteMessage(ctx context.Context, id uint) error {
	// Check if message exists
	message, err := s.messageRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete message
	if err := s.messageRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Dispatch message deleted event
	event := &events.MessageDeletedEvent{
		MessageID: id,
		ChannelID: message.ChannelID,
		UserID:    message.SenderID,
	}
	return s.dispatcher.Dispatch(ctx, event)
}

// GetChannelMessages retrieves messages from a channel with pagination
func (s *MessagingServiceImpl) GetChannelMessages(ctx context.Context, channelID uint, limit, offset int) ([]*domain.MessagingMessage, error) {
	return s.messageRepo.GetChannelMessages(ctx, channelID, limit, offset)
}

// GetThreadMessages retrieves messages in a thread
func (s *MessagingServiceImpl) GetThreadMessages(ctx context.Context, threadID uint, limit, offset int) ([]*domain.MessagingMessage, error) {
	return s.messageRepo.GetThreadMessages(ctx, threadID, limit, offset)
}

// PinMessage pins a message in a channel
func (s *MessagingServiceImpl) PinMessage(ctx context.Context, messageID uint, userID uint) error {
	// Get the message
	message, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return err
	}

	// Check if user has permission to manage messages
	hasPermission, err := s.HasPermission(ctx, userID, message.ChannelID, domain.Permission(16))
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("user does not have permission to pin messages")
	}

	// Get the channel
	channel, err := s.channelRepo.FindByID(ctx, message.ChannelID)
	if err != nil {
		return err
	}

	// Create pinned message
	pinnedMessage := domain.MessagingPinnedMessage{
		Message:  *message,
		PinnedBy: userID,
	}

	// Update channel's pinned messages
	channel.PinnedMessages = append(channel.PinnedMessages, pinnedMessage)
	if err := s.channelRepo.Update(ctx, channel); err != nil {
		return err
	}

	// Dispatch message pinned event
	s.dispatcher.Dispatch(ctx, events.NewMessageEvent(events.MessagingEventMessagePinned, message))

	return nil
}

// UnpinMessage unpins a message from a channel
func (s *MessagingServiceImpl) UnpinMessage(ctx context.Context, messageID uint, userID uint) error {
	// Get the message
	message, err := s.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return err
	}

	// Check if user has permission to manage messages
	hasPermission, err := s.HasPermission(ctx, userID, message.ChannelID, domain.Permission(16))
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("user does not have permission to unpin messages")
	}

	// Get the channel
	channel, err := s.channelRepo.FindByID(ctx, message.ChannelID)
	if err != nil {
		return err
	}

	// Remove message from pinned messages
	for i, pinnedMessage := range channel.PinnedMessages {
		if pinnedMessage.MessageID == messageID {
			channel.PinnedMessages = append(channel.PinnedMessages[:i], channel.PinnedMessages[i+1:]...)
			break
		}
	}

	// Update channel
	if err := s.channelRepo.Update(ctx, channel); err != nil {
		return err
	}

	// Dispatch message unpinned event
	s.dispatcher.Dispatch(ctx, events.NewMessageEvent(events.MessagingEventMessageUnpinned, message))

	return nil
}

// SearchMessages searches for messages using filters
func (s *MessagingServiceImpl) SearchMessages(ctx context.Context, filters repository.MessageSearchFilters) ([]*domain.MessagingMessage, error) {
	return s.messageRepo.SearchMessages(ctx, interfaces.MessageSearchFilters{
		ChannelID: *filters.ChannelID,
		UserID:    *filters.UserID,
		Query:     filters.Query,
		Limit:     filters.Limit,
		Offset:    filters.Offset,
	})
}

// CreateChannel creates a new channel
func (s *MessagingServiceImpl) CreateChannel(ctx context.Context, channel *domain.MessagingChannel) error {
	return s.channelRepo.Create(ctx, channel)
}

// GetChannel retrieves a channel by ID
func (s *MessagingServiceImpl) GetChannel(ctx context.Context, id uint) (*domain.MessagingChannel, error) {
	return s.channelRepo.FindByID(ctx, id)
}

// UpdateChannel updates an existing channel
func (s *MessagingServiceImpl) UpdateChannel(ctx context.Context, channel *domain.MessagingChannel) error {
	return s.channelRepo.Update(ctx, channel)
}

// DeleteChannel deletes a channel
func (s *MessagingServiceImpl) DeleteChannel(ctx context.Context, id uint) error {
	return s.channelRepo.Delete(ctx, id)
}

// GetUserChannels retrieves channels for a user
func (s *MessagingServiceImpl) GetUserChannels(ctx context.Context, userID uint) ([]*domain.MessagingChannel, error) {
	return s.channelRepo.GetUserChannels(ctx, userID)
}

// AddChannelMember adds a member to a channel
func (s *MessagingServiceImpl) AddChannelMember(ctx context.Context, member *domain.MessagingChannelMember) error {
	return s.channelRepo.AddMember(ctx, member.ChannelID, member.UserID)
}

// RemoveChannelMember removes a member from a channel
func (s *MessagingServiceImpl) RemoveChannelMember(ctx context.Context, channelID, userID uint) error {
	return s.channelRepo.RemoveMember(ctx, channelID, userID)
}

// GetChannelMembers retrieves members of a channel
func (s *MessagingServiceImpl) GetChannelMembers(ctx context.Context, channelID uint) ([]*domain.MessagingChannelMember, error) {
	users, err := s.channelRepo.GetMembers(ctx, channelID)
	if err != nil {
		return nil, err
	}

	members := make([]*domain.MessagingChannelMember, len(users))
	for i, user := range users {
		members[i] = &domain.MessagingChannelMember{
			ChannelID: channelID,
			UserID:    user.ID,
			User:      *user,
		}
	}
	return members, nil
}

// SetChannelSlowMode sets slow mode for a channel
func (s *MessagingServiceImpl) SetChannelSlowMode(ctx context.Context, channelID uint, seconds int, userID uint) error {
	// Check if user has permission to manage channel
	hasPermission, err := s.HasPermission(ctx, userID, channelID, domain.Permission(4))
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("user does not have permission to manage channel")
	}

	// Get the channel
	channel, err := s.channelRepo.FindByID(ctx, channelID)
	if err != nil {
		return err
	}

	// Update channel's slow mode
	channel.RateLimit = seconds
	if err := s.channelRepo.Update(ctx, channel); err != nil {
		return err
	}

	// Dispatch channel updated event
	s.dispatcher.Dispatch(ctx, events.NewChannelEvent(events.MessagingEventChannelUpdated, channel, channelID, userID))

	return nil
}

// SetChannelNSFW sets the NSFW flag for a channel
func (s *MessagingServiceImpl) SetChannelNSFW(ctx context.Context, channelID uint, isNSFW bool, userID uint) error {
	// Check if user has permission to manage the channel
	hasPermission, err := s.HasPermission(ctx, userID, channelID, domain.PermissionManageChannel)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("user does not have permission to manage this channel")
	}

	// Get the channel
	channel, err := s.channelRepo.GetChannel(ctx, channelID)
	if err != nil {
		return err
	}

	// Update NSFW flag
	channel.IsNSFW = isNSFW
	channel.UpdatedAt = time.Now()

	if err := s.channelRepo.UpdateChannel(ctx, channel); err != nil {
		return err
	}

	// Dispatch channel updated event
	s.dispatcher.Dispatch(ctx, events.NewChannelEvent(events.MessagingEventChannelUpdated, channel, channelID, userID))

	return nil
}

// CreateCategory creates a new category
func (s *MessagingServiceImpl) CreateCategory(ctx context.Context, category *domain.MessagingCategory) error {
	return s.categoryRepo.Create(ctx, category)
}

// GetCategory retrieves a category by ID
func (s *MessagingServiceImpl) GetCategory(ctx context.Context, id uint) (*domain.MessagingCategory, error) {
	return s.categoryRepo.FindByID(ctx, id)
}

// UpdateCategory updates an existing category
func (s *MessagingServiceImpl) UpdateCategory(ctx context.Context, category *domain.MessagingCategory) error {
	return s.categoryRepo.Update(ctx, category)
}

// DeleteCategory deletes a category
func (s *MessagingServiceImpl) DeleteCategory(ctx context.Context, id uint) error {
	return s.categoryRepo.Delete(ctx, id)
}

// GetChannelCategories retrieves categories for a channel
func (s *MessagingServiceImpl) GetChannelCategories(ctx context.Context, channelID uint) ([]*domain.MessagingCategory, error) {
	return s.categoryRepo.GetChannelCategories(ctx, channelID)
}

// UpdateUserStatus updates a user's status
func (s *MessagingServiceImpl) UpdateUserStatus(ctx context.Context, userID uint, status, statusMessage string) error {
	// Get the user
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	// Update status
	user.Status = status
	user.StatusMessage = statusMessage
	user.LastSeen = time.Now()

	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}

	// Dispatch status updated event
	s.dispatcher.Dispatch(ctx, events.NewStatusEvent(userID, status, statusMessage))

	return nil
}

// GetUserStatus retrieves a user's status
func (s *MessagingServiceImpl) GetUserStatus(ctx context.Context, userID uint) (*domain.MessagingUser, error) {
	return s.userRepo.GetUser(ctx, userID)
}

// BlockUser blocks a user
func (s *MessagingServiceImpl) BlockUser(ctx context.Context, userID, blockedID uint) error {
	if err := s.userRepo.BlockUser(ctx, userID, blockedID); err != nil {
		return err
	}

	// Dispatch user blocked event
	s.dispatcher.Dispatch(ctx, events.NewUserEvent(events.MessagingEventUserBlocked, userID, blockedID, 0))

	return nil
}

// UnblockUser unblocks a user
func (s *MessagingServiceImpl) UnblockUser(ctx context.Context, userID, blockedID uint) error {
	if err := s.userRepo.UnblockUser(ctx, userID, blockedID); err != nil {
		return err
	}

	// Dispatch user unblocked event
	s.dispatcher.Dispatch(ctx, events.NewUserEvent(events.MessagingEventUserUnblocked, userID, blockedID, 0))

	return nil
}

// MuteUser mutes a user
func (s *MessagingServiceImpl) MuteUser(ctx context.Context, userID, mutedID uint, duration time.Duration) error {
	if err := s.userRepo.MuteUser(ctx, userID, mutedID); err != nil {
		return err
	}

	// Dispatch user muted event
	s.dispatcher.Dispatch(ctx, events.NewUserEvent(events.MessagingEventUserMuted, userID, mutedID, duration))

	return nil
}

// UnmuteUser unmutes a user
func (s *MessagingServiceImpl) UnmuteUser(ctx context.Context, userID, mutedID uint) error {
	if err := s.userRepo.UnmuteUser(ctx, userID, mutedID); err != nil {
		return err
	}

	// Dispatch user unmuted event
	s.dispatcher.Dispatch(ctx, events.NewUserEvent(events.MessagingEventUserUnmuted, userID, mutedID, 0))

	return nil
}

// GetBlockedUsers retrieves all users blocked by a user
func (s *MessagingServiceImpl) GetBlockedUsers(ctx context.Context, userID uint) ([]domain.MessagingUser, error) {
	return s.userRepo.GetBlockedUsers(ctx, userID)
}

// GetMutedUsers retrieves all users muted by a user
func (s *MessagingServiceImpl) GetMutedUsers(ctx context.Context, userID uint) ([]domain.MessagingUser, error) {
	return s.userRepo.GetMutedUsers(ctx, userID)
}

// CreateRole creates a new role
func (s *MessagingServiceImpl) CreateRole(ctx context.Context, role *domain.MessagingRole) error {
	// Create the role
	role := &domain.MessagingRole{
		Name:        role.Name,
		Color:       role.Color,
		Permissions: role.Permissions,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.roleRepo.CreateRole(ctx, role); err != nil {
		return err
	}

	// Dispatch role created event
	s.dispatcher.Dispatch(ctx, events.NewRoleEvent(events.MessagingEventRoleCreated, role, 0, role.ID, role.UserID))

	return nil
}

// GetRole retrieves a role by ID
func (s *MessagingServiceImpl) GetRole(ctx context.Context, roleID uint) (*domain.MessagingRole, error) {
	return s.roleRepo.GetRole(ctx, roleID)
}

// UpdateRole updates an existing role
func (s *MessagingServiceImpl) UpdateRole(ctx context.Context, role *domain.MessagingRole) error {
	return s.roleRepo.UpdateRole(ctx, role)
}

// DeleteRole deletes a role
func (s *MessagingServiceImpl) DeleteRole(ctx context.Context, id uint) error {
	return s.roleRepo.Delete(ctx, id)
}

// AddUserRole adds a role to a user
func (s *MessagingServiceImpl) AddUserRole(ctx context.Context, userID, roleID uint, addedByID uint) error {
	if err := s.roleRepo.AddUserRole(ctx, userID, roleID); err != nil {
		return err
	}

	// Get the role
	role, err := s.roleRepo.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	// Dispatch role added event
	s.dispatcher.Dispatch(ctx, events.NewRoleEvent(events.MessagingEventRoleAdded, role, userID, roleID, addedByID))

	return nil
}

// RemoveUserRole removes a role from a user
func (s *MessagingServiceImpl) RemoveUserRole(ctx context.Context, userID, roleID uint, removedByID uint) error {
	if err := s.roleRepo.RemoveUserRole(ctx, userID, roleID); err != nil {
		return err
	}

	// Get the role
	role, err := s.roleRepo.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	// Dispatch role removed event
	s.dispatcher.Dispatch(ctx, events.NewRoleEvent(events.MessagingEventRoleRemoved, role, userID, roleID, removedByID))

	return nil
}

// GetUserRoles retrieves all roles of a user
func (s *MessagingServiceImpl) GetUserRoles(ctx context.Context, userID uint) ([]*domain.MessagingRole, error) {
	return s.roleRepo.GetUserRoles(ctx, userID)
}

// AddReaction adds a reaction to a message
func (s *MessagingServiceImpl) AddReaction(ctx context.Context, reaction *domain.MessagingReaction) error {
	return s.reactionRepo.Create(ctx, reaction)
}

// RemoveReaction removes a reaction from a message
func (s *MessagingServiceImpl) RemoveReaction(ctx context.Context, messageID, userID uint, emoji string) error {
	return s.reactionRepo.Delete(ctx, messageID)
}

// GetMessageReactions retrieves reactions for a message
func (s *MessagingServiceImpl) GetMessageReactions(ctx context.Context, messageID uint) ([]*domain.MessagingReaction, error) {
	return s.reactionRepo.GetMessageReactions(ctx, messageID)
}

// MarkMessageAsRead marks a message as read
func (s *MessagingServiceImpl) MarkMessageAsRead(ctx context.Context, receipt *domain.MessagingReadReceipt) error {
	return s.readReceiptRepo.Create(ctx, receipt)
}

// GetUnreadCount retrieves the number of unread messages
func (s *MessagingServiceImpl) GetUnreadCount(ctx context.Context, userID, channelID uint) (int64, error) {
	return s.readReceiptRepo.GetUnreadCount(ctx, userID, channelID)
}

// GetReadReceipts retrieves read receipts for a message
func (s *MessagingServiceImpl) GetReadReceipts(ctx context.Context, messageID uint) ([]*domain.MessagingReadReceipt, error) {
	return s.readReceiptRepo.GetMessageReceipts(ctx, messageID)
}

// HasPermission checks if a user has a specific permission
func (s *MessagingServiceImpl) HasPermission(ctx context.Context, userID uint, channelID uint, permission domain.Permission) (bool, error) {
	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role.Permissions.HasPermission(permission) {
			return true, nil
		}
	}

	return false, nil
}

// GetUserPermissions retrieves permissions for a user
func (s *MessagingServiceImpl) GetUserPermissions(ctx context.Context, userID uint, channelID uint) ([]domain.Permission, error) {
	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	var permissions []domain.Permission
	for _, role := range roles {
		for i := 0; i < 64; i++ {
			permission := domain.Permission(1 << i)
			if role.Permissions.HasPermission(permission) {
				permissions = append(permissions, permission)
			}
		}
	}

	return permissions, nil
}

// Helper function to get user ID from context
func getUserIDFromContext(ctx context.Context) uint {
	userID, ok := ctx.Value("userID").(uint)
	if !ok {
		return 0
	}
	return userID
}

// validateMessage validates a message before creation or update
func (s *MessagingServiceImpl) validateMessage(message *domain.MessagingMessage) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}
	if message.Content == "" && len(message.Attachments) == 0 && len(message.Embeds) == 0 {
		return errors.New("message must have content, attachments, or embeds")
	}
	if message.SenderID == 0 {
		return errors.New("message must have a sender")
	}
	if message.ChannelID == 0 {
		return errors.New("message must have a channel")
	}
	return nil
}

// GetUser retrieves a user by ID
func (s *MessagingServiceImpl) GetUser(ctx context.Context, id uint) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// GetUserByUsername retrieves a user by username
func (s *MessagingServiceImpl) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

// GetUserByEmail retrieves a user by email
func (s *MessagingServiceImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// UpdateUser updates an existing user
func (s *MessagingServiceImpl) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.userRepo.Update(ctx, user)
}

// DeleteUser deletes a user
func (s *MessagingServiceImpl) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

// AssignRole assigns a role to a user
func (s *MessagingServiceImpl) AssignRole(ctx context.Context, userID, roleID uint) error {
	return s.roleRepo.AssignRole(ctx, userID, roleID)
}

// RemoveRole removes a role from a user
func (s *MessagingServiceImpl) RemoveRole(ctx context.Context, userID, roleID uint) error {
	return s.roleRepo.RemoveRole(ctx, userID, roleID)
}
