package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// StorageType represents the type of storage backend
type StorageType string

const (
	// StorageTypeS3 is for AWS S3 or compatible storage
	StorageTypeS3 StorageType = "s3"

	// StorageTypeLocal is for local filesystem storage
	StorageTypeLocal StorageType = "local"
)

// Config contains configuration for the storage provider
type Config struct {
	// Type defines the storage backend type
	Type StorageType `json:"type" mapstructure:"type"`

	// BasePath is the base path for file storage (local filesystem) or prefix (S3)
	BasePath string `json:"base_path" mapstructure:"base_path"`

	// BaseURL is the base URL for generating public URLs
	BaseURL string `json:"base_url" mapstructure:"base_url"`

	// S3 contains S3-specific configuration
	S3 S3Config `json:"s3" mapstructure:"s3"`
}

// S3Config contains configuration for S3 storage
type S3Config struct {
	// Bucket is the S3 bucket name
	Bucket string `json:"bucket" mapstructure:"bucket"`

	// Region is the AWS region
	Region string `json:"region" mapstructure:"region"`

	// Endpoint is the S3 API endpoint (optional, for non-AWS S3-compatible services)
	Endpoint string `json:"endpoint" mapstructure:"endpoint"`

	// AccessKey is the S3 access key
	AccessKey string `json:"access_key" mapstructure:"access_key"`

	// SecretKey is the S3 secret key
	SecretKey string `json:"secret_key" mapstructure:"secret_key"`

	// UseSSL indicates whether to use SSL for S3 connections
	UseSSL bool `json:"use_ssl" mapstructure:"use_ssl"`
}

// DefaultConfig returns the default storage configuration
func DefaultConfig() *Config {
	return &Config{
		Type:     StorageTypeLocal,
		BasePath: "./storage",
		BaseURL:  "http://localhost:8080/files",
		S3: S3Config{
			Bucket: "project-website",
			Region: "us-west-2",
			UseSSL: true,
		},
	}
}

// Provider is a storage provider with support for multiple backends
type Provider struct {
	config       *Config
	s3Client     *s3.S3
	s3Uploader   *s3manager.Uploader
	s3Downloader *s3manager.Downloader
}

// FileInfo contains information about a stored file
type FileInfo struct {
	Key      string    `json:"key"`
	Size     int64     `json:"size"`
	URL      string    `json:"url"`
	MimeType string    `json:"mime_type"`
	ModTime  time.Time `json:"mod_time"`
}

// NewProvider creates a new storage provider
func NewProvider(config *Config) (*Provider, error) {
	if config == nil {
		config = DefaultConfig()
	}

	provider := &Provider{
		config: config,
	}

	// Initialize the appropriate backend
	switch config.Type {
	case StorageTypeS3:
		if err := provider.initS3(); err != nil {
			return nil, fmt.Errorf("failed to initialize S3 storage: %w", err)
		}
	case StorageTypeLocal:
		// No initialization needed for local storage
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}

	return provider, nil
}

// initS3 initializes the S3 client
func (p *Provider) initS3() error {
	s3Config := p.config.S3

	if s3Config.Bucket == "" {
		return errors.New("S3 bucket name is required")
	}

	// Create AWS session
	awsConfig := &aws.Config{
		Region:      aws.String(s3Config.Region),
		Credentials: credentials.NewStaticCredentials(s3Config.AccessKey, s3Config.SecretKey, ""),
	}

	// Use custom endpoint if provided (for non-AWS S3-compatible services)
	if s3Config.Endpoint != "" {
		awsConfig.Endpoint = aws.String(s3Config.Endpoint)
		awsConfig.DisableSSL = aws.Bool(!s3Config.UseSSL)
		awsConfig.S3ForcePathStyle = aws.Bool(true)
	}

	// Create session
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Create S3 client and helpers
	p.s3Client = s3.New(sess)
	p.s3Uploader = s3manager.NewUploader(sess)
	p.s3Downloader = s3manager.NewDownloader(sess)

	return nil
}

// Upload uploads a file to storage
func (p *Provider) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (*FileInfo, error) {
	switch p.config.Type {
	case StorageTypeS3:
		return p.uploadS3(ctx, key, reader, contentType)
	case StorageTypeLocal:
		return p.uploadLocal(ctx, key, reader, contentType)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", p.config.Type)
	}
}

// uploadS3 uploads a file to S3
func (p *Provider) uploadS3(ctx context.Context, key string, reader io.Reader, contentType string) (*FileInfo, error) {
	// Normalize key by removing leading slash and adding prefix if set
	if strings.HasPrefix(key, "/") {
		key = key[1:]
	}

	if p.config.BasePath != "" && !strings.HasPrefix(key, p.config.BasePath) {
		key = path.Join(p.config.BasePath, key)
	}

	// Upload the file
	result, err := p.s3Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(p.config.S3.Bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Get file metadata
	head, err := p.s3Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(p.config.S3.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file metadata from S3: %w", err)
	}

	// Generate URL
	fileURL := result.Location
	if p.config.BaseURL != "" {
		fileURL = p.generateURL(key)
	}

	return &FileInfo{
		Key:      key,
		Size:     *head.ContentLength,
		URL:      fileURL,
		MimeType: contentType,
		ModTime:  *head.LastModified,
	}, nil
}

// uploadLocal uploads a file to local storage
func (p *Provider) uploadLocal(ctx context.Context, key string, reader io.Reader, contentType string) (*FileInfo, error) {
	// TODO: Implement local file storage
	// This is a placeholder for actual implementation
	return nil, errors.New("local storage not implemented")
}

// Download downloads a file from storage
func (p *Provider) Download(ctx context.Context, key string) (io.ReadCloser, *FileInfo, error) {
	switch p.config.Type {
	case StorageTypeS3:
		return p.downloadS3(ctx, key)
	case StorageTypeLocal:
		return p.downloadLocal(ctx, key)
	default:
		return nil, nil, fmt.Errorf("unsupported storage type: %s", p.config.Type)
	}
}

// downloadS3 downloads a file from S3
func (p *Provider) downloadS3(ctx context.Context, key string) (io.ReadCloser, *FileInfo, error) {
	// Normalize key by removing leading slash and adding prefix if set
	if strings.HasPrefix(key, "/") {
		key = key[1:]
	}

	if p.config.BasePath != "" && !strings.HasPrefix(key, p.config.BasePath) {
		key = path.Join(p.config.BasePath, key)
	}

	// Get file metadata
	head, err := p.s3Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(p.config.S3.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file metadata from S3: %w", err)
	}

	// Get the file
	result, err := p.s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.config.S3.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download file from S3: %w", err)
	}

	// Generate URL
	fileURL := p.generateURL(key)

	fileInfo := &FileInfo{
		Key:      key,
		Size:     *head.ContentLength,
		URL:      fileURL,
		MimeType: *head.ContentType,
		ModTime:  *head.LastModified,
	}

	return result.Body, fileInfo, nil
}

// downloadLocal downloads a file from local storage
func (p *Provider) downloadLocal(ctx context.Context, key string) (io.ReadCloser, *FileInfo, error) {
	// TODO: Implement local file storage
	// This is a placeholder for actual implementation
	return nil, nil, errors.New("local storage not implemented")
}

// Delete deletes a file from storage
func (p *Provider) Delete(ctx context.Context, key string) error {
	switch p.config.Type {
	case StorageTypeS3:
		return p.deleteS3(ctx, key)
	case StorageTypeLocal:
		return p.deleteLocal(ctx, key)
	default:
		return fmt.Errorf("unsupported storage type: %s", p.config.Type)
	}
}

// deleteS3 deletes a file from S3
func (p *Provider) deleteS3(ctx context.Context, key string) error {
	// Normalize key by removing leading slash and adding prefix if set
	if strings.HasPrefix(key, "/") {
		key = key[1:]
	}

	if p.config.BasePath != "" && !strings.HasPrefix(key, p.config.BasePath) {
		key = path.Join(p.config.BasePath, key)
	}

	// Delete the file
	_, err := p.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.config.S3.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// deleteLocal deletes a file from local storage
func (p *Provider) deleteLocal(ctx context.Context, key string) error {
	// TODO: Implement local file storage
	// This is a placeholder for actual implementation
	return errors.New("local storage not implemented")
}

// generateURL generates a URL for the given key
func (p *Provider) generateURL(key string) string {
	baseURL := p.config.BaseURL

	// Remove trailing slash from base URL
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Remove leading slash from key
	key = strings.TrimPrefix(key, "/")

	// Build the URL
	fileURL := fmt.Sprintf("%s/%s", baseURL, url.PathEscape(key))

	return fileURL
}
