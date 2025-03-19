package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Repository handles all database operations for the messaging platform
type Repository struct {
	db *sql.DB
}

// ChannelMember represents a user's membership in a channel
type ChannelMember struct {
	UserID    uint      `json:"userId"`
	ChannelID uint      `json:"channelId"`
	JoinedAt  time.Time `json:"joinedAt"`
	Role      string    `json:"role"`
}

// Channel represents a messaging channel
type Channel struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uint      `json:"ownerId"`
	IsPrivate   bool      `json:"isPrivate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Message represents a chat message
type Message struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	SenderID  uint      `json:"senderId"`
	ChannelID uint      `json:"channelId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewRepository creates a new SQLite repository for the messaging platform
func NewRepository(dbPath string) (*Repository, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("error creating database directory: %v", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Initialize tables
	if err := initTables(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("error initializing tables: %v", err)
	}

	return &Repository{db: db}, nil
}

// Close closes the database connection
func (r *Repository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// initTables initializes the database tables
func initTables(db *sql.DB) error {
	// Create channels table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS channels (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			owner_id INTEGER NOT NULL,
			is_private BOOLEAN DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create messages table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			sender_id INTEGER NOT NULL,
			channel_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (channel_id) REFERENCES channels (id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	// Create channel_members table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS channel_members (
			user_id INTEGER NOT NULL,
			channel_id INTEGER NOT NULL,
			joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			role TEXT DEFAULT 'member',
			PRIMARY KEY (user_id, channel_id),
			FOREIGN KEY (channel_id) REFERENCES channels (id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	// Create read_receipts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS read_receipts (
			user_id INTEGER NOT NULL,
			message_id INTEGER NOT NULL,
			read_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, message_id),
			FOREIGN KEY (message_id) REFERENCES messages (id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// GetChannel gets a channel by ID
func (r *Repository) GetChannel(id uint) (*Channel, error) {
	var channel Channel
	err := r.db.QueryRow(`
		SELECT id, name, description, owner_id, is_private, created_at, updated_at
		FROM channels
		WHERE id = ?
	`, id).Scan(
		&channel.ID,
		&channel.Name,
		&channel.Description,
		&channel.OwnerID,
		&channel.IsPrivate,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("channel not found: %d", id)
		}
		return nil, err
	}
	return &channel, nil
}

// CreateChannel creates a new channel
func (r *Repository) CreateChannel(channel *Channel) error {
	res, err := r.db.Exec(`
		INSERT INTO channels (name, description, owner_id, is_private, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		channel.Name,
		channel.Description,
		channel.OwnerID,
		channel.IsPrivate,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	channel.ID = uint(id)
	return nil
}

// GetUserChannels gets all channels for a user
func (r *Repository) GetUserChannels(userID uint) ([]Channel, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.name, c.description, c.owner_id, c.is_private, c.created_at, c.updated_at
		FROM channels c
		LEFT JOIN channel_members cm ON c.id = cm.channel_id
		WHERE c.owner_id = ? OR cm.user_id = ?
		ORDER BY c.updated_at DESC
	`, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []Channel
	for rows.Next() {
		var channel Channel
		if err := rows.Scan(
			&channel.ID,
			&channel.Name,
			&channel.Description,
			&channel.OwnerID,
			&channel.IsPrivate,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		); err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}

	return channels, nil
}

// GetChannelMembers gets all members of a channel
func (r *Repository) GetChannelMembers(channelID uint) ([]ChannelMember, error) {
	rows, err := r.db.Query(`
		SELECT user_id, channel_id, joined_at, role
		FROM channel_members
		WHERE channel_id = ?
	`, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []ChannelMember
	for rows.Next() {
		var member ChannelMember
		if err := rows.Scan(
			&member.UserID,
			&member.ChannelID,
			&member.JoinedAt,
			&member.Role,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

// AddChannelMember adds a user to a channel
func (r *Repository) AddChannelMember(userID, channelID uint, role string) error {
	_, err := r.db.Exec(`
		INSERT INTO channel_members (user_id, channel_id, joined_at, role)
		VALUES (?, ?, ?, ?)
	`, userID, channelID, time.Now(), role)
	return err
}

// RemoveChannelMember removes a user from a channel
func (r *Repository) RemoveChannelMember(userID, channelID uint) error {
	_, err := r.db.Exec(`
		DELETE FROM channel_members
		WHERE user_id = ? AND channel_id = ?
	`, userID, channelID)
	return err
}

// CreateMessage creates a new message
func (r *Repository) CreateMessage(message *Message) error {
	res, err := r.db.Exec(`
		INSERT INTO messages (content, sender_id, channel_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`,
		message.Content,
		message.SenderID,
		message.ChannelID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	message.ID = uint(id)
	return nil
}

// GetChannelMessages gets messages from a channel
func (r *Repository) GetChannelMessages(channelID uint, lastMessageID uint, limit int) ([]Message, error) {
	query := `
		SELECT id, content, sender_id, channel_id, created_at, updated_at
		FROM messages
		WHERE channel_id = ?
	`
	args := []interface{}{channelID}

	if lastMessageID > 0 {
		query += " AND id < ?"
		args = append(args, lastMessageID)
	}

	query += " ORDER BY id DESC LIMIT ?"
	args = append(args, limit)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(
			&message.ID,
			&message.Content,
			&message.SenderID,
			&message.ChannelID,
			&message.CreatedAt,
			&message.UpdatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// GetMessage gets a message by ID
func (r *Repository) GetMessage(id uint) (*Message, error) {
	var message Message
	err := r.db.QueryRow(`
		SELECT id, content, sender_id, channel_id, created_at, updated_at
		FROM messages
		WHERE id = ?
	`, id).Scan(
		&message.ID,
		&message.Content,
		&message.SenderID,
		&message.ChannelID,
		&message.CreatedAt,
		&message.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found: %d", id)
		}
		return nil, err
	}
	return &message, nil
}
