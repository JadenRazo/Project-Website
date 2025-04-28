package database

import (
	"context"
	"database/sql"
	"strings"

	"github.com/JadenRazo/Project-Website/backend/internal/messaging/domain"
	"github.com/google/uuid"
)

type WordFilterRepository struct {
	db *sql.DB
}

func NewWordFilterRepository(db *sql.DB) *WordFilterRepository {
	return &WordFilterRepository{db: db}
}

func (r *WordFilterRepository) CreateFilter(ctx context.Context, filter *domain.WordFilter) error {
	query := `
		INSERT INTO word_filters (server_id, created_by, word, scope, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		filter.ServerID,
		filter.CreatedBy,
		filter.Word,
		filter.Scope,
		filter.Description,
	).Scan(&filter.ID, &filter.CreatedAt, &filter.UpdatedAt)

	if err != nil {
		if isDuplicateError(err) {
			return domain.NewDuplicateEntryError("word filter already exists")
		}
		return domain.NewInternalError("failed to create word filter")
	}

	return nil
}

func (r *WordFilterRepository) UpdateFilter(ctx context.Context, filter *domain.WordFilter) error {
	query := `
		UPDATE word_filters
		SET word = $1, scope = $2, description = $3, is_active = $4
		WHERE id = $5 AND server_id = $6
		RETURNING updated_at
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		filter.Word,
		filter.Scope,
		filter.Description,
		filter.IsActive,
		filter.ID,
		filter.ServerID,
	)

	if err != nil {
		return domain.NewInternalError("failed to update word filter")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError("failed to check rows affected")
	}

	if rowsAffected == 0 {
		return domain.NewNotFoundError("word filter not found")
	}

	return nil
}

func (r *WordFilterRepository) DeleteFilter(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM word_filters
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return domain.NewInternalError("failed to delete word filter")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError("failed to check rows affected")
	}

	if rowsAffected == 0 {
		return domain.NewNotFoundError("word filter not found")
	}

	return nil
}

func (r *WordFilterRepository) GetFilter(ctx context.Context, id uuid.UUID) (*domain.WordFilter, error) {
	query := `
		SELECT id, server_id, created_by, word, scope, description, created_at, updated_at, is_active
		FROM word_filters
		WHERE id = $1
	`

	filter := &domain.WordFilter{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&filter.ID,
		&filter.ServerID,
		&filter.CreatedBy,
		&filter.Word,
		&filter.Scope,
		&filter.Description,
		&filter.CreatedAt,
		&filter.UpdatedAt,
		&filter.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewNotFoundError("word filter not found")
		}
		return nil, domain.NewInternalError("failed to get word filter")
	}

	return filter, nil
}

func (r *WordFilterRepository) ListFilters(ctx context.Context, serverID uuid.UUID) ([]*domain.WordFilter, error) {
	query := `
		SELECT id, server_id, created_by, word, scope, description, created_at, updated_at, is_active
		FROM word_filters
		WHERE server_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, serverID)
	if err != nil {
		return nil, domain.NewInternalError("failed to list word filters")
	}
	defer rows.Close()

	var filters []*domain.WordFilter
	for rows.Next() {
		filter := &domain.WordFilter{}
		err := rows.Scan(
			&filter.ID,
			&filter.ServerID,
			&filter.CreatedBy,
			&filter.Word,
			&filter.Scope,
			&filter.Description,
			&filter.CreatedAt,
			&filter.UpdatedAt,
			&filter.IsActive,
		)
		if err != nil {
			return nil, domain.NewInternalError("failed to scan word filter")
		}
		filters = append(filters, filter)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.NewInternalError("failed to iterate word filters")
	}

	return filters, nil
}

func (r *WordFilterRepository) GetFiltersByScope(ctx context.Context, serverID uuid.UUID, scope domain.WordFilterScope) ([]*domain.WordFilter, error) {
	query := `
		SELECT id, server_id, created_by, word, scope, description, created_at, updated_at, is_active
		FROM word_filters
		WHERE server_id = $1 AND scope = $2 AND is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, serverID, scope)
	if err != nil {
		return nil, domain.NewInternalError("failed to get word filters by scope")
	}
	defer rows.Close()

	var filters []*domain.WordFilter
	for rows.Next() {
		filter := &domain.WordFilter{}
		err := rows.Scan(
			&filter.ID,
			&filter.ServerID,
			&filter.CreatedBy,
			&filter.Word,
			&filter.Scope,
			&filter.Description,
			&filter.CreatedAt,
			&filter.UpdatedAt,
			&filter.IsActive,
		)
		if err != nil {
			return nil, domain.NewInternalError("failed to scan word filter")
		}
		filters = append(filters, filter)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.NewInternalError("failed to iterate word filters")
	}

	return filters, nil
}

func (r *WordFilterRepository) CheckMessage(ctx context.Context, serverID uuid.UUID, message string, scope domain.WordFilterScope) (bool, []*domain.WordFilter, error) {
	query := `
		SELECT id, server_id, created_by, word, scope, description, created_at, updated_at, is_active
		FROM word_filters
		WHERE server_id = $1 AND scope = $2 AND is_active = true
	`

	rows, err := r.db.QueryContext(ctx, query, serverID, scope)
	if err != nil {
		return false, nil, domain.NewInternalError("failed to check message against word filters")
	}
	defer rows.Close()

	var matchingFilters []*domain.WordFilter
	for rows.Next() {
		filter := &domain.WordFilter{}
		err := rows.Scan(
			&filter.ID,
			&filter.ServerID,
			&filter.CreatedBy,
			&filter.Word,
			&filter.Scope,
			&filter.Description,
			&filter.CreatedAt,
			&filter.UpdatedAt,
			&filter.IsActive,
		)
		if err != nil {
			return false, nil, domain.NewInternalError("failed to scan word filter")
		}

		// Check if the word appears in the message
		if containsWord(message, filter.Word) {
			matchingFilters = append(matchingFilters, filter)
		}
	}

	if err = rows.Err(); err != nil {
		return false, nil, domain.NewInternalError("failed to iterate word filters")
	}

	return len(matchingFilters) > 0, matchingFilters, nil
}

// Helper function to check if a word appears in a message
func containsWord(message, word string) bool {
	// Convert both to lowercase for case-insensitive comparison
	message = strings.ToLower(message)
	word = strings.ToLower(word)

	// Check for exact word match
	if strings.Contains(message, word) {
		return true
	}

	// Check for word with punctuation
	punctuation := []string{" ", ".", ",", "!", "?", ";", ":", "'", "\"", "(", ")", "[", "]", "{", "}", "<", ">", "/", "\\", "|", "-", "_", "+", "=", "*", "&", "^", "%", "$", "#", "@", "~", "`"}
	for _, p := range punctuation {
		if strings.Contains(message, word+p) || strings.Contains(message, p+word) {
			return true
		}
	}

	return false
}

// Helper function to check for duplicate entry errors
func isDuplicateError(err error) bool {
	// This is PostgreSQL specific, adjust for your database
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
