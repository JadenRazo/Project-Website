package contact

import (
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContactRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
	Website  string `json:"website"`
}

type ContactSubmission struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);not null"`
	Subject   string    `gorm:"type:varchar(500)"`
	Message   string    `gorm:"type:text;not null"`
	IPAddress string    `gorm:"type:inet"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (ContactSubmission) TableName() string {
	return "contact_submissions"
}

type rateLimitEntry struct {
	count     int
	windowEnd time.Time
}

type Handler struct {
	db          *gorm.DB
	emailConfig *EmailConfig
	rateLimits  sync.Map
}

func NewHandler(db *gorm.DB, emailCfg *EmailConfig) *Handler {
	h := &Handler{
		db:          db,
		emailConfig: emailCfg,
	}
	go h.cleanupRateLimits()
	return h
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("", h.HandleContactForm)
}

func (h *Handler) HandleContactForm(c *gin.Context) {
	var req ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	if req.Website != "" {
		logger.Info("Honeypot triggered", "ip", c.ClientIP())
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Message sent successfully",
		})
		return
	}

	req.Name = sanitizeInput(strings.TrimSpace(req.Name))
	req.Email = strings.TrimSpace(req.Email)
	req.Subject = sanitizeInput(strings.TrimSpace(req.Subject))
	req.Message = sanitizeInput(strings.TrimSpace(req.Message))

	if errs := validateContactRequest(&req); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Validation failed",
			"details": errs,
		})
		return
	}

	clientIP := c.ClientIP()
	if !h.checkRateLimit(clientIP) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"success": false,
			"error":   "Too many submissions. Please try again later.",
		})
		return
	}

	if h.db != nil {
		submission := ContactSubmission{
			Name:      req.Name,
			Email:     req.Email,
			Subject:   req.Subject,
			Message:   req.Message,
			IPAddress: clientIP,
		}
		if err := h.db.Create(&submission).Error; err != nil {
			logger.Error("Failed to save contact submission", "error", err)
		}
	}

	if h.emailConfig != nil && h.emailConfig.IsConfigured() {
		emailMsg := &EmailMessage{
			Name:    req.Name,
			Email:   req.Email,
			Subject: req.Subject,
			Message: req.Message,
		}
		if err := SendEmail(h.emailConfig, emailMsg); err != nil {
			logger.Error("Failed to send contact email", "error", err, "from", req.Email)
		}
	} else {
		logger.Info("Contact form submission (email not configured)",
			"name", req.Name,
			"email", req.Email,
			"subject", req.Subject,
			"message_length", len(req.Message),
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Message sent successfully",
	})
}

func validateContactRequest(req *ContactRequest) []string {
	var errs []string

	if req.Name == "" {
		errs = append(errs, "Name is required")
	} else if len(req.Name) < 2 {
		errs = append(errs, "Name must be at least 2 characters")
	} else if len(req.Name) > 255 {
		errs = append(errs, "Name must be 255 characters or fewer")
	}

	if req.Email == "" {
		errs = append(errs, "Email is required")
	} else if !isValidEmail(req.Email) {
		errs = append(errs, "Invalid email address")
	} else if isDisposableEmail(req.Email) {
		errs = append(errs, "Disposable email addresses are not allowed")
	} else if !hasValidMX(req.Email) {
		errs = append(errs, "Email domain does not accept mail")
	}

	if req.Message == "" {
		errs = append(errs, "Message is required")
	} else if len(req.Message) < 10 {
		errs = append(errs, "Message must be at least 10 characters")
	} else if len(req.Message) > 5000 {
		errs = append(errs, "Message must be 5000 characters or fewer")
	}

	if len(req.Subject) > 500 {
		errs = append(errs, "Subject must be 500 characters or fewer")
	}

	if len(errs) == 0 {
		if reason := detectSpam(req); reason != "" {
			errs = append(errs, reason)
		}
	}

	return errs
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	if len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

var disposableDomains = map[string]bool{
	"tempmail.com": true, "throwaway.email": true, "guerrillamail.com": true,
	"guerrillamail.info": true, "mailinator.com": true, "yopmail.com": true,
	"sharklasers.com": true, "guerrillamailblock.com": true, "grr.la": true,
	"dispostable.com": true, "trashmail.com": true, "trashmail.me": true,
	"10minutemail.com": true, "tempail.com": true, "fakeinbox.com": true,
	"mailnesia.com": true, "maildrop.cc": true, "discard.email": true,
	"getnada.com": true, "temp-mail.org": true, "mohmal.com": true,
	"burnermail.io": true, "minutemail.com": true, "tempr.email": true,
	"emailondeck.com": true, "33mail.com": true, "getairmail.com": true,
	"mailsac.com": true, "mytemp.email": true, "tempmailo.com": true,
	"harakirimail.com": true, "spamgourmet.com": true, "trash-mail.com": true,
	"crazymailing.com": true, "mailcatch.com": true, "temp-mail.io": true,
	"mailnull.com": true, "spamfree24.org": true, "jetable.org": true,
}

func isDisposableEmail(email string) bool {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return false
	}
	domain := strings.ToLower(parts[1])
	return disposableDomains[domain]
}

func hasValidMX(email string) bool {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]
	mx, err := net.LookupMX(domain)
	if err == nil && len(mx) > 0 {
		return true
	}
	addrs, err := net.LookupHost(domain)
	return err == nil && len(addrs) > 0
}

var (
	htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
	urlRegex     = regexp.MustCompile(`https?://[^\s]+`)
)

func sanitizeInput(s string) string {
	s = htmlTagRegex.ReplaceAllString(s, "")
	s = strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\r' && r != '\t' {
			return -1
		}
		return r
	}, s)
	return s
}

func detectSpam(req *ContactRequest) string {
	combined := req.Name + " " + req.Subject + " " + req.Message

	urls := urlRegex.FindAllString(combined, -1)
	if len(urls) > 3 {
		return "Message contains too many links"
	}

	upper := 0
	total := 0
	for _, r := range req.Message {
		if unicode.IsLetter(r) {
			total++
			if unicode.IsUpper(r) {
				upper++
			}
		}
	}
	if total > 20 && float64(upper)/float64(total) > 0.7 {
		return "Please avoid excessive use of capital letters"
	}

	lowerMsg := strings.ToLower(combined)
	spamPhrases := []string{
		"buy now", "click here", "free money", "act now",
		"limited time offer", "congratulations you won",
		"earn extra cash", "double your income",
		"no obligation", "risk free", "as seen on",
		"casino online", "online pharmacy",
		"viagra", "cryptocurrency investment",
	}
	for _, phrase := range spamPhrases {
		if strings.Contains(lowerMsg, phrase) {
			return "Message flagged as potential spam"
		}
	}

	if len(req.Message) > 50 {
		words := strings.Fields(req.Message)
		if len(words) > 5 {
			wordCount := make(map[string]int)
			for _, w := range words {
				wordCount[strings.ToLower(w)]++
			}
			for _, count := range wordCount {
				if count > len(words)/2 && count > 3 {
					return "Message appears to contain repetitive content"
				}
			}
		}
	}

	return ""
}

func (h *Handler) checkRateLimit(ip string) bool {
	now := time.Now()
	maxPerDay := 1

	val, loaded := h.rateLimits.Load(ip)
	if !loaded {
		h.rateLimits.Store(ip, &rateLimitEntry{
			count:     1,
			windowEnd: now.Add(24 * time.Hour),
		})
		return true
	}

	entry := val.(*rateLimitEntry)
	if now.After(entry.windowEnd) {
		entry.count = 1
		entry.windowEnd = now.Add(24 * time.Hour)
		return true
	}

	if entry.count >= maxPerDay {
		return false
	}

	entry.count++
	return true
}

func (h *Handler) cleanupRateLimits() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		h.rateLimits.Range(func(key, value interface{}) bool {
			entry := value.(*rateLimitEntry)
			if now.After(entry.windowEnd) {
				h.rateLimits.Delete(key)
			}
			return true
		})
	}
}
