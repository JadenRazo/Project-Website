package visitor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ComplianceService struct {
	db     *gorm.DB
	config *ConfigService
}

func NewComplianceService(db *gorm.DB, config *ConfigService) *ComplianceService {
	return &ComplianceService{
		db:     db,
		config: config,
	}
}

// CheckDNT checks if the Do Not Track header is set
func (cs *ComplianceService) CheckDNT(r *http.Request) bool {
	dnt := r.Header.Get("DNT")
	return dnt == "1" && cs.config.GetConfig().DataCollection.RespectDNT
}

// AnonymizeIP anonymizes IP addresses based on configuration
func (cs *ComplianceService) AnonymizeIP(ip string) string {
	if !cs.config.GetConfig().AnonymizationOptions.AnonymizeIP {
		return ip
	}

	mode := cs.config.GetConfig().AnonymizationOptions.IPAnonymizationMode

	// Parse IP
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ""
	}

	switch mode {
	case "remove_last_octet":
		if parsedIP.To4() != nil {
			// IPv4
			parts := strings.Split(ip, ".")
			if len(parts) == 4 {
				return fmt.Sprintf("%s.%s.%s.0", parts[0], parts[1], parts[2])
			}
		} else {
			// IPv6 - remove last 2 groups
			parts := strings.Split(ip, ":")
			if len(parts) > 2 {
				parts[len(parts)-1] = "0"
				parts[len(parts)-2] = "0"
				return strings.Join(parts, ":")
			}
		}

	case "remove_last_two":
		if parsedIP.To4() != nil {
			parts := strings.Split(ip, ".")
			if len(parts) == 4 {
				return fmt.Sprintf("%s.%s.0.0", parts[0], parts[1])
			}
		}

	case "hash":
		hash := sha256.Sum256([]byte(ip))
		return hex.EncodeToString(hash[:16]) // Return first 16 bytes for brevity
	}

	return ip
}

// ProcessConsentRequest handles consent updates
func (cs *ComplianceService) ProcessConsentRequest(ctx context.Context, req ConsentRequest) error {
	// Validate request
	if req.SessionHash == "" {
		return fmt.Errorf("session hash required")
	}

	// Store consent for each category
	tx := cs.db.Begin()
	defer tx.Rollback()

	for category, granted := range req.Categories {
		consent := PrivacyConsent{
			SessionHash: req.SessionHash,
			ConsentType: category,
			Granted:     granted,
			CreatedAt:   time.Now(),
			ExpiresAt:   time.Now().Add(365 * 24 * time.Hour),
		}

		if err := tx.Create(&consent).Error; err != nil {
			return err
		}
	}

	// Create audit log entry
	auditLog := ConsentAuditLog{
		ID:          uuid.New(),
		SessionHash: req.SessionHash,
		Action:      "consent_updated",
		Details:     req.Categories,
		Timestamp:   time.Now(),
	}

	if err := tx.Create(&auditLog).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

// ExportUserData exports all data for a user (GDPR Article 20 - Data Portability)
func (cs *ComplianceService) ExportUserData(ctx context.Context, sessionHash string) (*UserDataExport, error) {
	export := &UserDataExport{
		ExportDate:  time.Now(),
		SessionHash: sessionHash,
	}

	// Get session data
	var session VisitorSession
	if err := cs.db.Where("session_hash = ?", sessionHash).First(&session).Error; err == nil {
		export.SessionData = session
	}

	// Get page views
	var pageViews []PageView
	cs.db.Joins("JOIN visitor_sessions ON visitor_sessions.id = page_views.session_id").
		Where("visitor_sessions.session_hash = ?", sessionHash).
		Find(&pageViews)
	export.PageViews = pageViews

	// Get consent records
	var consents []PrivacyConsent
	cs.db.Where("session_hash = ?", sessionHash).Find(&consents)
	export.ConsentRecords = consents

	// Get events if collected
	var events []VisitorEvent
	cs.db.Where("session_hash = ?", sessionHash).Find(&events)
	export.Events = events

	return export, nil
}

// DeleteUserData deletes all data for a user (GDPR Article 17 - Right to Erasure)
func (cs *ComplianceService) DeleteUserData(ctx context.Context, sessionHash string) error {
	tx := cs.db.Begin()
	defer tx.Rollback()

	// Get session ID first
	var session VisitorSession
	if err := tx.Where("session_hash = ?", sessionHash).First(&session).Error; err != nil {
		return err
	}

	// Delete page views
	if err := tx.Where("session_id = ?", session.ID).Delete(&PageView{}).Error; err != nil {
		return err
	}

	// Delete events
	if err := tx.Where("session_hash = ?", sessionHash).Delete(&VisitorEvent{}).Error; err != nil {
		return err
	}

	// Delete consent records
	if err := tx.Where("session_hash = ?", sessionHash).Delete(&PrivacyConsent{}).Error; err != nil {
		return err
	}

	// Delete realtime tracking
	if err := tx.Where("session_hash = ?", sessionHash).Delete(&VisitorRealtime{}).Error; err != nil {
		return err
	}

	// Finally delete the session
	if err := tx.Where("session_hash = ?", sessionHash).Delete(&VisitorSession{}).Error; err != nil {
		return err
	}

	// Create deletion audit log
	auditLog := ConsentAuditLog{
		ID:          uuid.New(),
		SessionHash: sessionHash,
		Action:      "data_deleted",
		Details:     map[string]bool{"complete_deletion": true},
		Timestamp:   time.Now(),
	}

	if err := tx.Create(&auditLog).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

// GetDataDisclosure provides information about data collection (CCPA requirement)
func (cs *ComplianceService) GetDataDisclosure(ctx context.Context) *DataDisclosure {
	config := cs.config.GetConfig()

	return &DataDisclosure{
		DataCollected: []string{
			conditionalAdd("Browser information", config.DataCollection.CollectBrowserInfo),
			conditionalAdd("Device information", config.DataCollection.CollectDeviceInfo),
			conditionalAdd("Geographic location (country/region)", config.DataCollection.CollectGeographicData),
			conditionalAdd("Page views and navigation", config.DataCollection.CollectSessionData),
			conditionalAdd("Referrer information", config.DataCollection.CollectReferrers),
		},
		Purpose: []string{
			"Analytics and site improvement",
			"Performance monitoring",
			"User experience optimization",
			"Security and fraud prevention",
		},
		DataSharing: DataSharingInfo{
			ThirdParties: []string{},
			Purpose:      "We do not sell or share your data with third parties",
		},
		RetentionPeriod: fmt.Sprintf("%d days for session data, %d days for aggregated data",
			config.Retention.SessionDataDays,
			config.Retention.AggregatedDataDays),
		UserRights: []string{
			"Right to access your data",
			"Right to delete your data",
			"Right to opt-out of tracking",
			"Right to data portability",
		},
		ContactInfo: "privacy@example.com",
	}
}

// HandleOptOut processes CCPA opt-out requests
func (cs *ComplianceService) HandleOptOut(ctx context.Context, sessionHash string) error {
	// Set all consents to false
	consents := map[string]bool{
		"analytics":  false,
		"functional": false,
		"marketing":  false,
	}

	req := ConsentRequest{
		SessionHash: sessionHash,
		Categories:  consents,
	}

	return cs.ProcessConsentRequest(ctx, req)
}

// ValidateProcessingBasis checks if we have valid basis for processing (GDPR Article 6)
func (cs *ComplianceService) ValidateProcessingBasis(ctx context.Context, sessionHash string, purpose string) bool {
	config := cs.config.GetConfig()

	// Check if GDPR is enabled
	if !config.Compliance.GDPR.Enabled {
		return true
	}

	switch config.Compliance.GDPR.ProcessingBasis {
	case "consent":
		// Check if we have consent
		var consent PrivacyConsent
		err := cs.db.Where("session_hash = ? AND consent_type = ? AND granted = true",
			sessionHash, purpose).First(&consent).Error
		return err == nil

	case "legitimate_interest":
		// For legitimate interest, we allow basic analytics
		return purpose == "analytics" || purpose == "necessary"

	case "contract":
		// For contract basis, we assume it's valid
		return true

	default:
		return false
	}
}

// ApplyRetentionPolicy deletes data based on retention configuration
func (cs *ComplianceService) ApplyRetentionPolicy(ctx context.Context) error {
	config := cs.config.GetConfig()

	if !config.Retention.EnableAutoDelete {
		return nil
	}

	tx := cs.db.Begin()
	defer tx.Rollback()

	now := time.Now()

	// Delete old sessions
	sessionCutoff := now.AddDate(0, 0, -config.Retention.SessionDataDays)
	if err := tx.Where("created_at < ?", sessionCutoff).Delete(&VisitorSession{}).Error; err != nil {
		return err
	}

	// Delete old page views
	pageViewCutoff := now.AddDate(0, 0, -config.Retention.PageViewDataDays)
	if err := tx.Where("created_at < ?", pageViewCutoff).Delete(&PageView{}).Error; err != nil {
		return err
	}

	// Delete old consent records
	consentCutoff := now.AddDate(0, 0, -config.Retention.ConsentRecordDays)
	if err := tx.Where("created_at < ?", consentCutoff).Delete(&PrivacyConsent{}).Error; err != nil {
		return err
	}

	// Delete inactive sessions
	if config.Retention.DeleteInactiveAfter > 0 {
		inactiveCutoff := now.AddDate(0, 0, -config.Retention.DeleteInactiveAfter)
		if err := tx.Where("last_seen_at < ?", inactiveCutoff).Delete(&VisitorSession{}).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

// Helper function
func conditionalAdd(item string, condition bool) string {
	if condition {
		return item
	}
	return ""
}

// Request/Response structures

type ConsentRequest struct {
	SessionHash string          `json:"sessionHash"`
	Categories  map[string]bool `json:"categories"`
	IPAddress   string          `json:"-"`
	UserAgent   string          `json:"-"`
}

type UserDataExport struct {
	ExportDate     time.Time        `json:"exportDate"`
	SessionHash    string           `json:"sessionHash"`
	SessionData    VisitorSession   `json:"sessionData"`
	PageViews      []PageView       `json:"pageViews"`
	Events         []VisitorEvent   `json:"events,omitempty"`
	ConsentRecords []PrivacyConsent `json:"consentRecords"`
}

type DataDisclosure struct {
	DataCollected   []string         `json:"dataCollected"`
	Purpose         []string         `json:"purpose"`
	DataSharing     DataSharingInfo  `json:"dataSharing"`
	RetentionPeriod string           `json:"retentionPeriod"`
	UserRights      []string         `json:"userRights"`
	ContactInfo     string           `json:"contactInfo"`
}

type DataSharingInfo struct {
	ThirdParties []string `json:"thirdParties"`
	Purpose      string   `json:"purpose"`
}

// Additional models

type VisitorEvent struct {
	ID          uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SessionHash string                 `gorm:"type:varchar(64);index"`
	EventType   string                 `gorm:"type:varchar(100)"`
	EventData   map[string]interface{} `gorm:"type:jsonb"`
	CreatedAt   time.Time
}

type ConsentAuditLog struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SessionHash string          `gorm:"type:varchar(64);index"`
	Action      string          `gorm:"type:varchar(100)"`
	Details     map[string]bool `gorm:"type:jsonb"`
	Timestamp   time.Time
	IPAddress   string          `gorm:"type:varchar(64)"`
	UserAgent   string          `gorm:"type:text"`
}

// Extended PrivacyConsent with additional fields
type PrivacyConsentExtended struct {
	PrivacyConsent
	IPAddress string `gorm:"type:varchar(64)"`
	UserAgent string `gorm:"type:text"`
	Timestamp time.Time
}

// Table names
func (VisitorEvent) TableName() string      { return "visitor_events" }
func (ConsentAuditLog) TableName() string   { return "consent_audit_logs" }

// ComplianceHandlers - HTTP handlers for compliance endpoints

func (cs *ComplianceService) RegisterHandlers(router *gin.RouterGroup) {
	compliance := router.Group("/compliance")
	{
		compliance.POST("/consent", cs.handleConsent)
		compliance.GET("/consent/:sessionHash", cs.handleGetConsent)
		compliance.POST("/export/:sessionHash", cs.handleExportData)
		compliance.DELETE("/data/:sessionHash", cs.handleDeleteData)
		compliance.POST("/optout/:sessionHash", cs.handleOptOut)
		compliance.GET("/disclosure", cs.handleDataDisclosure)
		compliance.GET("/privacy-policy", cs.handlePrivacyPolicy)
	}
}

func (cs *ComplianceService) handleConsent(c *gin.Context) {
	var req ConsentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	req.IPAddress = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")

	if err := cs.ProcessConsentRequest(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (cs *ComplianceService) handleGetConsent(c *gin.Context) {
	sessionHash := c.Param("sessionHash")

	var consents []PrivacyConsent
	if err := cs.db.Where("session_hash = ?", sessionHash).Find(&consents).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No consent records found"})
		return
	}

	result := make(map[string]bool)
	for _, consent := range consents {
		result[consent.ConsentType] = consent.Granted
	}

	c.JSON(http.StatusOK, result)
}

func (cs *ComplianceService) handleExportData(c *gin.Context) {
	sessionHash := c.Param("sessionHash")

	export, err := cs.ExportUserData(c.Request.Context(), sessionHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return as downloadable JSON
	data, _ := json.MarshalIndent(export, "", "  ")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"user-data-%s.json\"", sessionHash[:8]))
	c.Data(http.StatusOK, "application/json", data)
}

func (cs *ComplianceService) handleDeleteData(c *gin.Context) {
	sessionHash := c.Param("sessionHash")

	if err := cs.DeleteUserData(c.Request.Context(), sessionHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All data has been permanently deleted",
	})
}

func (cs *ComplianceService) handleOptOut(c *gin.Context) {
	sessionHash := c.Param("sessionHash")

	if err := cs.HandleOptOut(c.Request.Context(), sessionHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "You have been opted out of all tracking",
	})
}

func (cs *ComplianceService) handleDataDisclosure(c *gin.Context) {
	disclosure := cs.GetDataDisclosure(c.Request.Context())
	c.JSON(http.StatusOK, disclosure)
}

func (cs *ComplianceService) handlePrivacyPolicy(c *gin.Context) {
	// This would typically serve your privacy policy
	// For now, return a structured response
	config := cs.config.GetConfig()

	policy := gin.H{
		"lastUpdated": config.UpdatedAt,
		"regulations": gin.H{
			"gdpr":   config.Compliance.GDPR.Enabled,
			"ccpa":   config.Compliance.CCPA.Enabled,
			"lgpd":   config.Compliance.LGPD.Enabled,
			"pipeda": config.Compliance.PIPEDA.Enabled,
		},
		"dataCollection": config.DataCollection,
		"retention":      config.Retention,
		"userRights": []string{
			"Right to access",
			"Right to rectification",
			"Right to erasure",
			"Right to data portability",
			"Right to object",
			"Right to opt-out",
		},
	}

	c.JSON(http.StatusOK, policy)
}