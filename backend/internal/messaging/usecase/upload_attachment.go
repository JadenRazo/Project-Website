package usecase

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/errors"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/events"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// UploadAttachmentUseCase handles file attachment uploads
type UploadAttachmentUseCase struct {
	messageRepo    repository.MessagingMessageRepository
	attachmentRepo repository.MessagingAttachmentRepository
	fileStorage    FileStorage
	dispatcher     events.EventDispatcher
}

// FileStorage interface for handling file storage operations
type FileStorage interface {
	UploadFile(file *multipart.FileHeader, path string) (string, error)
	DeleteFile(path string) error
	GetFileURL(path string) string
}

// NewUploadAttachmentUseCase creates a new upload attachment usecase
func NewUploadAttachmentUseCase(
	messageRepo repository.MessagingMessageRepository,
	attachmentRepo repository.MessagingAttachmentRepository,
	fileStorage FileStorage,
	dispatcher events.EventDispatcher,
) *UploadAttachmentUseCase {
	return &UploadAttachmentUseCase{
		messageRepo:    messageRepo,
		attachmentRepo: attachmentRepo,
		fileStorage:    fileStorage,
		dispatcher:     dispatcher,
	}
}

// UploadAttachmentInput represents the input for uploading an attachment
type UploadAttachmentInput struct {
	MessageID uint
	UserID    uint
	File      *multipart.FileHeader
}

// Execute uploads a file attachment with validation
func (uc *UploadAttachmentUseCase) Execute(ctx context.Context, input UploadAttachmentInput) (*domain.MessagingAttachment, error) {
	// Get the message
	message, err := uc.messageRepo.GetMessage(ctx, input.MessageID)
	if err != nil {
		return nil, errors.ErrMessageNotFound
	}

	// Check if user has permission to add attachments
	if message.SenderID != input.UserID {
		return nil, errors.ErrUnauthorized
	}

	// Validate file size (max 10MB)
	if input.File.Size > 10*1024*1024 {
		return nil, errors.ErrFileTooLarge
	}

	// Validate file type
	ext := filepath.Ext(input.File.Filename)
	if !isAllowedFileType(ext) {
		return nil, errors.ErrInvalidFileType
	}

	// Generate unique filename
	filename := generateUniqueFilename(input.File.Filename)

	// Upload file to storage
	fileURL, err := uc.fileStorage.UploadFile(input.File, filename)
	if err != nil {
		return nil, err
	}

	// Create attachment record
	attachment := &domain.MessagingAttachment{
		MessageID: input.MessageID,
		FileName:  input.File.Filename,
		FileType:  ext,
		FileSize:  input.File.Size,
		FileURL:   fileURL,
		IsImage:   isImageFile(ext),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save attachment to database
	if err := uc.attachmentRepo.AddAttachment(ctx, attachment); err != nil {
		// Clean up uploaded file if database save fails
		uc.fileStorage.DeleteFile(filename)
		return nil, err
	}

	// Dispatch attachment created event
	event := events.NewAttachmentEvent(events.MessagingEventAttachmentCreated, attachment)
	uc.dispatcher.Dispatch(ctx, event)

	return attachment, nil
}

// Helper functions

// isAllowedFileType checks if the file extension is allowed
func isAllowedFileType(ext string) bool {
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".txt":  true,
		".zip":  true,
		".rar":  true,
	}
	return allowedTypes[ext]
}

// isImageFile checks if the file is an image
func isImageFile(ext string) bool {
	imageTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	return imageTypes[ext]
}

// generateUniqueFilename generates a unique filename for upload
func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	return time.Now().Format("20060102150405") + ext
}
