package attachments

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"  // Register GIF format
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/errors"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/events"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/websocket"
	"github.com/google/uuid"
)

// StorageProvider defines the interface for file storage
type StorageProvider interface {
	// StoreFile stores a file and returns its unique identifier and public URL
	StoreFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, string, error)

	// DeleteFile deletes a file by its identifier
	DeleteFile(ctx context.Context, fileID string) error
}

// AttachmentRepository defines the interface for attachment data access
type AttachmentRepository interface {
	CreateAttachment(ctx context.Context, attachment *domain.Attachment) error
	UpdateAttachment(ctx context.Context, attachment *domain.Attachment) error
	GetAttachment(ctx context.Context, attachmentID uint) (*domain.Attachment, error)
	DeleteAttachment(ctx context.Context, attachmentID uint) error
	GetMessageAttachments(ctx context.Context, messageID uint) ([]domain.Attachment, error)
}

// AttachmentService manages message attachments
type AttachmentService struct {
	storage    StorageProvider
	repo       AttachmentRepository
	hub        *websocket.Hub
	dispatcher *events.EventDispatcher
	maxSize    int64
}

// NewAttachmentService creates a new attachment service
func NewAttachmentService(
	storage StorageProvider,
	repo AttachmentRepository,
	hub *websocket.Hub,
	dispatcher *events.EventDispatcher,
	maxSize int64,
) *AttachmentService {
	return &AttachmentService{
		storage:    storage,
		repo:       repo,
		hub:        hub,
		dispatcher: dispatcher,
		maxSize:    maxSize, // Maximum file size in bytes
	}
}

// UploadAttachment handles file upload and creates an attachment record
func (s *AttachmentService) UploadAttachment(
	ctx context.Context,
	file multipart.File,
	header *multipart.FileHeader,
	messageID uint,
	channelID uint,
	userID uint,
) (*domain.Attachment, error) {
	// Validate file size
	if header.Size > s.maxSize {
		return nil, errors.ErrAttachmentTooLarge
	}

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !s.isAllowedFileType(contentType) {
		return nil, errors.ErrInvalidAttachmentType
	}

	// Generate a unique filename to prevent collisions
	filename := s.generateUniqueFilename(header.Filename)

	// Create attachment record
	attachment := &domain.Attachment{
		MessageID: messageID,
		FileName:  header.Filename,
		FileType:  contentType,
		FileSize:  header.Size,
		IsImage:   s.isImageType(contentType),
	}

	// Begin upload - broadcast uploading status
	s.broadcastAttachmentStatus(attachment, "uploading", 0, channelID)

	// Check if it's an image that needs processing
	if attachment.IsImage {
		// Read the file for image processing while keeping a copy for storage
		var buf bytes.Buffer
		tee := io.TeeReader(file, &buf)

		// Get image dimensions (for future use, not stored in current model)
		_, format, err := image.DecodeConfig(tee)
		if err == nil {
			attachment.FileType = "image/" + format
		}

		// Reset file reader position
		file.Seek(0, io.SeekStart)
	}

	// Save the attachment record to get an ID
	if err := s.repo.CreateAttachment(ctx, attachment); err != nil {
		return nil, errors.WrapError(err, errors.ErrorTypeDatabase, "db_error", "Failed to create attachment record")
	}

	// Upload to storage provider
	fileID, fileURL, err := s.storage.StoreFile(ctx, file, filename, contentType)
	if err != nil {
		// Handle upload failure
		s.repo.UpdateAttachment(ctx, attachment)

		s.broadcastAttachmentStatus(attachment, "error", 0, channelID)
		return nil, errors.WrapError(err, errors.ErrorTypeInternal, "upload_failed", "Failed to upload file")
	}

	// Update attachment with storage details
	// Note: fileID is stored separately for internal tracking
	_ = fileID // We'll use this for deletion
	attachment.FileURL = fileURL

	// Save the updated attachment
	if err := s.repo.UpdateAttachment(ctx, attachment); err != nil {
		// Log error but don't fail - file is already uploaded
		// Consider a cleanup process for orphaned files
	}

	// Broadcast attachment complete
	s.broadcastAttachmentStatus(attachment, "complete", 100, channelID)

	return attachment, nil
}

// DeleteAttachment removes an attachment
func (s *AttachmentService) DeleteAttachment(ctx context.Context, attachmentID uint, userID uint) error {
	// Get the attachment
	attachment, err := s.repo.GetAttachment(ctx, attachmentID)
	if err != nil {
		return errors.ErrAttachmentNotFound
	}

	// For now, we'll assume any authenticated user can delete their message attachments
	// This would normally check against the message owner or channel admin
	// TODO: Add proper permission checking when message ownership is available

	// Delete from storage using the FileURL to derive the file ID
	// In a real implementation, we'd store the fileID separately
	if attachment.FileURL != "" {
		// Extract file ID from URL or use URL as identifier
		if err := s.storage.DeleteFile(ctx, attachment.FileURL); err != nil {
			// Log error but continue with deletion from database
		}
	}

	// Delete from database
	if err := s.repo.DeleteAttachment(ctx, attachmentID); err != nil {
		return errors.WrapError(err, errors.ErrorTypeDatabase, "delete_failed", "Failed to delete attachment")
	}

	return nil
}

// GetAttachments retrieves attachments for a message
func (s *AttachmentService) GetAttachments(ctx context.Context, messageID uint) ([]domain.Attachment, error) {
	return s.repo.GetMessageAttachments(ctx, messageID)
}

// Helper functions

// isAllowedFileType checks if the file type is allowed
func (s *AttachmentService) isAllowedFileType(contentType string) bool {
	allowedTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"application/pdf",
		"text/plain",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}

	for _, allowed := range allowedTypes {
		if contentType == allowed {
			return true
		}
	}

	return false
}

// isImageType determines if the content type is an image
func (s *AttachmentService) isImageType(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}

// generateUniqueFilename creates a unique filename with original extension
func (s *AttachmentService) generateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}

// broadcastAttachmentStatus sends attachment status updates to WebSocket clients
func (s *AttachmentService) broadcastAttachmentStatus(
	attachment *domain.Attachment,
	status string,
	progress int,
	channelID uint,
) {
	// Broadcast via WebSocket
	s.hub.BroadcastToChannel(
		channelID,
		websocket.EventTypeAttachment,
		websocket.AttachmentEvent{
			Type:         websocket.EventTypeAttachment,
			AttachmentID: attachment.ID,
			MessageID:    attachment.MessageID,
			ChannelID:    channelID,
			Status:       status,
			Progress:     progress,
			Timestamp:    time.Now().Unix(),
		},
	)
}
