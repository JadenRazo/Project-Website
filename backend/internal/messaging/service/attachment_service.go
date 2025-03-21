package service

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// AttachmentService handles file attachments in messages
type AttachmentService struct {
	attachmentRepo repository.MessagingAttachmentRepository
	fileStorage    FileStorage
}

// FileStorage interface for handling file storage operations
type FileStorage interface {
	UploadFile(file *multipart.FileHeader, path string) (string, error)
	DeleteFile(path string) error
	GetFileURL(path string) string
}

// NewAttachmentService creates a new attachment service
func NewAttachmentService(attachmentRepo repository.MessagingAttachmentRepository, fileStorage FileStorage) *AttachmentService {
	return &AttachmentService{
		attachmentRepo: attachmentRepo,
		fileStorage:    fileStorage,
	}
}

// UploadAttachment uploads a file attachment for a message
func (s *AttachmentService) UploadAttachment(ctx context.Context, file *multipart.FileHeader, messageID uint, userID uint) (*domain.MessagingAttachment, error) {
	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return nil, errors.New("file size exceeds 10MB limit")
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	if !isAllowedFileType(ext) {
		return nil, errors.New("file type not allowed")
	}

	// Generate unique filename
	filename := generateUniqueFilename(file.Filename)

	// Upload file to storage
	fileURL, err := s.fileStorage.UploadFile(file, filename)
	if err != nil {
		return nil, err
	}

	// Create attachment record
	attachment := &domain.MessagingAttachment{
		MessageID: messageID,
		FileName:  file.Filename,
		FileType:  ext,
		FileSize:  file.Size,
		FileURL:   fileURL,
		IsImage:   isImageFile(ext),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save attachment to database
	if err := s.attachmentRepo.AddAttachment(ctx, attachment); err != nil {
		// Clean up uploaded file if database save fails
		s.fileStorage.DeleteFile(filename)
		return nil, err
	}

	return attachment, nil
}

// DeleteAttachment deletes an attachment and its associated file
func (s *AttachmentService) DeleteAttachment(ctx context.Context, attachmentID uint, userID uint) error {
	// Get attachment to verify ownership
	attachment, err := s.attachmentRepo.GetAttachment(ctx, attachmentID)
	if err != nil {
		return err
	}

	// Delete file from storage
	if err := s.fileStorage.DeleteFile(attachment.FileURL); err != nil {
		return err
	}

	// Delete attachment from database
	return s.attachmentRepo.DeleteAttachment(ctx, attachmentID)
}

// GetAttachment retrieves an attachment by ID
func (s *AttachmentService) GetAttachment(ctx context.Context, attachmentID uint) (*domain.MessagingAttachment, error) {
	return s.attachmentRepo.GetAttachment(ctx, attachmentID)
}

// GetMessageAttachments retrieves all attachments for a message
func (s *AttachmentService) GetMessageAttachments(ctx context.Context, messageID uint) ([]domain.MessagingAttachment, error) {
	return s.attachmentRepo.GetMessageAttachments(ctx, messageID)
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
