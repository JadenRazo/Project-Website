package visitor

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type PrivacyConfig struct {
	ID                   string                `json:"id" gorm:"primary_key"`
	EnableTracking       bool                  `json:"enableTracking"`
	PrivacyMode          string                `json:"privacyMode"` // "strict", "balanced", "minimal"
	DataCollection       DataCollectionConfig  `json:"dataCollection"`
	Retention            RetentionConfig       `json:"retention"`
	Compliance           ComplianceConfig      `json:"compliance"`
	AnonymizationOptions AnonymizationConfig   `json:"anonymization"`
	ConsentRequirements  ConsentConfig         `json:"consent"`
	UpdatedAt            time.Time             `json:"updatedAt"`
	UpdatedBy            string                `json:"updatedBy"`
}

type DataCollectionConfig struct {
	CollectCookies        bool `json:"collectCookies"`
	CollectIPAddresses    bool `json:"collectIPAddresses"`
	CollectUserAgents     bool `json:"collectUserAgents"`
	CollectReferrers      bool `json:"collectReferrers"`
	CollectGeographicData bool `json:"collectGeographicData"`
	CollectSessionData    bool `json:"collectSessionData"`
	CollectEventData      bool `json:"collectEventData"`
	CollectDeviceInfo     bool `json:"collectDeviceInfo"`
	CollectBrowserInfo    bool `json:"collectBrowserInfo"`
	RespectDNT            bool `json:"respectDNT"`
	AnonymousMode         bool `json:"anonymousMode"`
}

type RetentionConfig struct {
	EnableAutoDelete      bool   `json:"enableAutoDelete"`
	SessionDataDays       int    `json:"sessionDataDays"`
	PageViewDataDays      int    `json:"pageViewDataDays"`
	AggregatedDataDays    int    `json:"aggregatedDataDays"`
	ConsentRecordDays     int    `json:"consentRecordDays"`
	DeleteInactiveAfter   int    `json:"deleteInactiveAfter"`
	RetentionPolicy       string `json:"retentionPolicy"` // "minimal", "standard", "extended"
}

type ComplianceConfig struct {
	GDPR struct {
		Enabled           bool `json:"enabled"`
		RequireConsent    bool `json:"requireConsent"`
		AllowPortability  bool `json:"allowPortability"`
		AllowErasure      bool `json:"allowErasure"`
		ProcessingBasis   string `json:"processingBasis"` // "consent", "legitimate_interest", "contract"
	} `json:"gdpr"`
	CCPA struct {
		Enabled                bool `json:"enabled"`
		AllowOptOut            bool `json:"allowOptOut"`
		ProvideDataDisclosure  bool `json:"provideDataDisclosure"`
		DoNotSellData          bool `json:"doNotSellData"`
	} `json:"ccpa"`
	LGPD struct {
		Enabled           bool `json:"enabled"`
		RequireConsent    bool `json:"requireConsent"`
		AllowPortability  bool `json:"allowPortability"`
		AllowErasure      bool `json:"allowErasure"`
		DataProcessingBasis string `json:"dataProcessingBasis"`
	} `json:"lgpd"`
	PIPEDA struct {
		Enabled              bool `json:"enabled"`
		RequireConsent       bool `json:"requireConsent"`
		LimitDataCollection  bool `json:"limitDataCollection"`
		ProvideAccess        bool `json:"provideAccess"`
	} `json:"pipeda"`
}

type AnonymizationConfig struct {
	AnonymizeIP         bool   `json:"anonymizeIP"`
	IPAnonymizationMode string `json:"ipAnonymizationMode"` // "remove_last_octet", "remove_last_two", "hash"
	HashSessionIDs      bool   `json:"hashSessionIDs"`
	RemovePII           bool   `json:"removePII"`
	UseFingerprinting   bool   `json:"useFingerprinting"`
	MaskUserAgents      bool   `json:"maskUserAgents"`
}

type ConsentConfig struct {
	RequireExplicitConsent bool     `json:"requireExplicitConsent"`
	ConsentCategories     []string `json:"consentCategories"` // "necessary", "analytics", "functional", "marketing"
	DefaultConsent        string   `json:"defaultConsent"` // "opt-in", "opt-out"
	ConsentDuration       int      `json:"consentDuration"` // days
	ShowBanner            bool     `json:"showBanner"`
	BannerPosition        string   `json:"bannerPosition"` // "top", "bottom", "center"
	AllowGranularControl  bool     `json:"allowGranularControl"`
	MinimumAge            int      `json:"minimumAge"`
}

type ConfigService struct {
	db     *gorm.DB
	config *PrivacyConfig
}

func NewConfigService(db *gorm.DB) *ConfigService {
	cs := &ConfigService{
		db: db,
	}
	cs.LoadOrCreateDefault()
	return cs
}

func (cs *ConfigService) LoadOrCreateDefault() error {
	var config PrivacyConfig
	err := cs.db.First(&config, "id = ?", "default").Error

	if err == gorm.ErrRecordNotFound {
		config = cs.GetDefaultConfig()
		if err := cs.db.Create(&config).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	cs.config = &config
	return nil
}

func (cs *ConfigService) GetDefaultConfig() PrivacyConfig {
	config := PrivacyConfig{
		ID:             "default",
		EnableTracking: true,
		PrivacyMode:    "balanced",
		DataCollection: DataCollectionConfig{
			CollectCookies:        false,
			CollectIPAddresses:    false,
			CollectUserAgents:     true,
			CollectReferrers:      true,
			CollectGeographicData: true,
			CollectSessionData:    true,
			CollectEventData:      true,
			CollectDeviceInfo:     true,
			CollectBrowserInfo:    true,
			RespectDNT:            true,
			AnonymousMode:         false,
		},
		Retention: RetentionConfig{
			EnableAutoDelete:    true,
			SessionDataDays:     30,
			PageViewDataDays:    90,
			AggregatedDataDays:  365,
			ConsentRecordDays:   365,
			DeleteInactiveAfter: 180,
			RetentionPolicy:     "standard",
		},
		AnonymizationOptions: AnonymizationConfig{
			AnonymizeIP:         true,
			IPAnonymizationMode: "remove_last_octet",
			HashSessionIDs:      true,
			RemovePII:           true,
			UseFingerprinting:   false,
			MaskUserAgents:      false,
		},
		ConsentRequirements: ConsentConfig{
			RequireExplicitConsent: true,
			ConsentCategories:     []string{"necessary", "analytics", "functional", "marketing"},
			DefaultConsent:        "opt-in",
			ConsentDuration:       365,
			ShowBanner:            true,
			BannerPosition:        "bottom",
			AllowGranularControl:  true,
			MinimumAge:            16,
		},
		UpdatedAt: time.Now(),
	}

	// Set default compliance settings
	config.Compliance.GDPR.Enabled = true
	config.Compliance.GDPR.RequireConsent = true
	config.Compliance.GDPR.AllowPortability = true
	config.Compliance.GDPR.AllowErasure = true
	config.Compliance.GDPR.ProcessingBasis = "consent"

	config.Compliance.CCPA.Enabled = true
	config.Compliance.CCPA.AllowOptOut = true
	config.Compliance.CCPA.ProvideDataDisclosure = true
	config.Compliance.CCPA.DoNotSellData = true

	config.Compliance.LGPD.Enabled = true
	config.Compliance.LGPD.RequireConsent = true
	config.Compliance.LGPD.AllowPortability = true
	config.Compliance.LGPD.AllowErasure = true
	config.Compliance.LGPD.DataProcessingBasis = "consent"

	config.Compliance.PIPEDA.Enabled = true
	config.Compliance.PIPEDA.RequireConsent = true
	config.Compliance.PIPEDA.LimitDataCollection = true
	config.Compliance.PIPEDA.ProvideAccess = true

	return config
}

func (cs *ConfigService) GetConfig() *PrivacyConfig {
	return cs.config
}

func (cs *ConfigService) UpdateConfig(config *PrivacyConfig) error {
	config.UpdatedAt = time.Now()
	if err := cs.db.Save(config).Error; err != nil {
		return err
	}
	cs.config = config
	return nil
}

func (cs *ConfigService) ShouldCollect(dataType string) bool {
	if !cs.config.EnableTracking {
		return false
	}

	switch dataType {
	case "cookies":
		return cs.config.DataCollection.CollectCookies
	case "ip":
		return cs.config.DataCollection.CollectIPAddresses
	case "userAgent":
		return cs.config.DataCollection.CollectUserAgents
	case "referrer":
		return cs.config.DataCollection.CollectReferrers
	case "geographic":
		return cs.config.DataCollection.CollectGeographicData
	case "session":
		return cs.config.DataCollection.CollectSessionData
	case "event":
		return cs.config.DataCollection.CollectEventData
	case "device":
		return cs.config.DataCollection.CollectDeviceInfo
	case "browser":
		return cs.config.DataCollection.CollectBrowserInfo
	default:
		return false
	}
}

func (cs *ConfigService) GetRetentionDays(dataType string) int {
	switch dataType {
	case "session":
		return cs.config.Retention.SessionDataDays
	case "pageView":
		return cs.config.Retention.PageViewDataDays
	case "aggregated":
		return cs.config.Retention.AggregatedDataDays
	case "consent":
		return cs.config.Retention.ConsentRecordDays
	default:
		return 30
	}
}

func (cs *ConfigService) IsComplianceEnabled(regulation string) bool {
	switch regulation {
	case "GDPR":
		return cs.config.Compliance.GDPR.Enabled
	case "CCPA":
		return cs.config.Compliance.CCPA.Enabled
	case "LGPD":
		return cs.config.Compliance.LGPD.Enabled
	case "PIPEDA":
		return cs.config.Compliance.PIPEDA.Enabled
	default:
		return false
	}
}

func (cs *ConfigService) ExportConfig() ([]byte, error) {
	return json.MarshalIndent(cs.config, "", "  ")
}

func (cs *ConfigService) ImportConfig(data []byte) error {
	var config PrivacyConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}
	return cs.UpdateConfig(&config)
}

func (PrivacyConfig) TableName() string {
	return "privacy_configs"
}