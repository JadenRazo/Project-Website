package domain

// Permission represents a user permission in the messaging system
type Permission uint64

const (
	// Message permissions
	PermissionSendMessages Permission = 1 << iota
	PermissionEditMessages
	PermissionDeleteMessages
	PermissionPinMessages
	PermissionReactToMessages
	PermissionReadMessageHistory

	// Channel permissions
	PermissionViewChannels
	PermissionManageChannels
	PermissionManageChannelPermissions
	PermissionCreateChannels
	PermissionDeleteChannels
	PermissionManageChannelCategories

	// User permissions
	PermissionManageUsers
	PermissionManageRoles
	PermissionManageUserPermissions
	PermissionBanUsers
	PermissionKickUsers
	PermissionMuteUsers

	// Server permissions
	PermissionManageServer
	PermissionManageServerSettings
	PermissionManageServerEmojis
	PermissionManageServerIntegrations
	PermissionManageServerWebhooks
	PermissionManageServerInvites
)

// HasPermission checks if a permission set includes a specific permission
func (p Permission) HasPermission(permission Permission) bool {
	return p&permission != 0
}

// AddPermission adds a permission to the permission set
func (p *Permission) AddPermission(permission Permission) {
	*p |= permission
}

// RemovePermission removes a permission from the permission set
func (p *Permission) RemovePermission(permission Permission) {
	*p &^= permission
}
