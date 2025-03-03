// Package utils provides utility functions for the URL shortener service
package utils

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

// URLValidator defines the interface for URL validation
type URLValidator interface {
	IsValid(rawURL string) bool
	ValidateWithReason(rawURL string) (bool, error)
	IsSafe(rawURL string) bool
	NormalizeURL(rawURL string) (string, error)
}

// ValidationRules defines the interface for validation rules
type ValidationRules interface {
	Apply(u *url.URL) (bool, error)
	Name() string
}

// StandardURLValidator implements the URLValidator interface
type StandardURLValidator struct {
	rules         []ValidationRules
	blockedHosts  map[string]bool
	blockedRegexp []*regexp.Regexp
	cache         *ValidationCache
	mu            sync.RWMutex
}

// ValidationCache provides caching for validation results to improve performance
type ValidationCache struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

// CacheItem represents a cached validation result
type CacheItem struct {
	valid      bool
	error      error
	expiration time.Time
}

// NewValidationCache creates a new validation cache with specified TTL
func NewValidationCache(ttl time.Duration) *ValidationCache {
	return &ValidationCache{
		items: make(map[string]*CacheItem),
		ttl:   ttl,
	}
}

// Get retrieves a validation result from the cache
func (c *ValidationCache) Get(key string) (bool, error, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, exists := c.items[key]; exists {
		if time.Now().Before(item.expiration) {
			return item.valid, item.error, true
		}
		// Expired, will be cleaned up later
		delete(c.items, key)
	}
	return false, nil, false
}

// Set stores a validation result in the cache
func (c *ValidationCache) Set(key string, valid bool, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &CacheItem{
		valid:      valid,
		error:      err,
		expiration: time.Now().Add(c.ttl),
	}

	// Periodically clean up expired items
	// In a production system, you might want a dedicated goroutine for this
	if len(c.items)%100 == 0 {
		go c.cleanup()
	}
}

// cleanup removes expired cache entries
func (c *ValidationCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.expiration) {
			delete(c.items, key)
		}
	}
}

// NewStandardURLValidator creates a new validator with default rules
func NewStandardURLValidator() *StandardURLValidator {
	v := &StandardURLValidator{
		rules: []ValidationRules{
			&SchemeRule{},
			&HostRule{},
			&PathRule{},
			&QueryParamRule{},
			&LengthRule{maxLength: 2048}, // Common max URL length
		},
		blockedHosts:  make(map[string]bool),
		blockedRegexp: make([]*regexp.Regexp, 0),
		cache:         NewValidationCache(10 * time.Minute),
	}

	// Default list of known malicious or unwanted domains
	// In a production system, this would likely be loaded from a database or config file
	v.BlockHost("malware-site.com")
	v.BlockHost("phishing-example.org")
	v.BlockHost("spam-domain.net")
	
	// Block suspicious patterns (regex)
	v.BlockPattern(`\.exe$`)
	v.BlockPattern(`\.zip$`)
	v.BlockPattern(`[0-9a-f]{32}`) // MD5 hash pattern often used in malicious URLs

	return v
}

// AddRule adds a new validation rule to the validator
func (v *StandardURLValidator) AddRule(rule ValidationRules) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.rules = append(v.rules, rule)
}

// BlockHost adds a domain to the blocklist
func (v *StandardURLValidator) BlockHost(host string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.blockedHosts[strings.ToLower(host)] = true
}

// BlockPattern adds a regex pattern to block URLs that match
func (v *StandardURLValidator) BlockPattern(pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern %q: %w", pattern, err)
	}
	
	v.mu.Lock()
	defer v.mu.Unlock()
	v.blockedRegexp = append(v.blockedRegexp, re)
	return nil
}

// IsValid returns true if the URL is valid according to all rules
func (v *StandardURLValidator) IsValid(rawURL string) bool {
	isValid, _ := v.ValidateWithReason(rawURL)
	return isValid
}

// ValidateWithReason checks if a URL is valid and returns an error describing why if not
func (v *StandardURLValidator) ValidateWithReason(rawURL string) (bool, error) {
	// Check cache first
	if valid, err, found := v.cache.Get(rawURL); found {
		return valid, err
	}
	
	// Parse the URL
	u, err := url.Parse(rawURL)
	if err != nil {
		v.cache.Set(rawURL, false, fmt.Errorf("failed to parse URL: %w", err))
		return false, fmt.Errorf("failed to parse URL: %w", err)
	}
	
	// Save a lowercase version of the host for comparison
	host := strings.ToLower(u.Hostname())
	
	// Check if domain is in blocklist
	v.mu.RLock()
	if _, blocked := v.blockedHosts[host]; blocked {
		v.mu.RUnlock()
		err := errors.New("URL domain is in blocklist")
		v.cache.Set(rawURL, false, err)
		return false, err
	}
	v.mu.RUnlock()
	
	// Check if URL matches any blocked patterns
	v.mu.RLock()
	for _, re := range v.blockedRegexp {
		if re.MatchString(rawURL) {
			v.mu.RUnlock()
			err := fmt.Errorf("URL matches blocked pattern: %s", re.String())
			v.cache.Set(rawURL, false, err)
			return false, err
		}
	}
	v.mu.RUnlock()
	
	// Apply all validation rules
	v.mu.RLock()
	defer v.mu.RUnlock()
	
	for _, rule := range v.rules {
		valid, err := rule.Apply(u)
		if !valid {
			v.cache.Set(rawURL, false, err)
			return false, err
		}
	}
	
	v.cache.Set(rawURL, true, nil)
	return true, nil
}

// IsSafe checks if a URL appears to be safe (not malicious or spam)
// This is a more advanced check beyond basic validity
func (v *StandardURLValidator) IsSafe(rawURL string) bool {
	// In a production system, this might call external malware/phishing APIs
	// For now, we're just checking our blocklists
	
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	
	host := strings.ToLower(u.Hostname())
	
	// Check if domain is in blocklist
	v.mu.RLock()
	defer v.mu.RUnlock()
	
	if _, blocked := v.blockedHosts[host]; blocked {
		return false
	}
	
	for _, re := range v.blockedRegexp {
		if re.MatchString(rawURL) {
			return false
		}
	}
	
	// Basic heuristics for suspicious URLs
	if strings.Count(rawURL, "http") > 1 {
		// URL containing multiple http/https is suspicious (common in phishing)
		return false
	}
	
	if len(host) > 50 {
		// Excessively long host names are often suspicious
		return false
	}
	
	// Check for excessive use of subdomains (potential sign of abuse)
	if strings.Count(host, ".") > 4 {
		return false
	}
	
	return true
}

// NormalizeURL standardizes a URL for storage and comparison
func (v *StandardURLValidator) NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	
	// Force HTTPS if scheme is HTTP and the domain supports HTTPS
	// In a real system, you might check if the domain supports HTTPS first
	if u.Scheme == "http" {
		supportsHTTPS, _ := domainSupportsHTTPS(u.Hostname())
		if supportsHTTPS {
			u.Scheme = "https"
		}
	}
	
	// Remove tracking parameters (common in URLs)
	q := u.Query()
	trackingParams := []string{"utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content", "fbclid", "gclid"}
	for _, param := range trackingParams {
		q.Del(param)
	}
	u.RawQuery = q.Encode()
	
	// Remove trailing slash from path if present (unless it's just '/')
	if len(u.Path) > 1 && strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path[:len(u.Path)-1]
	}
	
	// Force lowercase for hostname
	if u.Host != "" {
		u.Host = strings.ToLower(u.Host)
	}
	
	// Remove default ports (80 for HTTP, 443 for HTTPS)
	if u.Scheme == "http" && strings.HasSuffix(u.Host, ":80") {
		u.Host = u.Host[:len(u.Host)-3]
	} else if u.Scheme == "https" && strings.HasSuffix(u.Host, ":443") {
		u.Host = u.Host[:len(u.Host)-4]
	}
	
	return u.String(), nil
}

// Concrete validation rule implementations
// SchemeRule validates the URL scheme (protocol)
type SchemeRule struct{}

func (r *SchemeRule) Apply(u *url.URL) (bool, error) {
	if u.Scheme == "" {
		return false, errors.New("URL scheme is missing")
	}
	
	// Only allow HTTP and HTTPS (could be expanded based on your needs)
	if u.Scheme != "http" && u.Scheme != "https" {
		return false, fmt.Errorf("URL scheme must be http or https, got %s", u.Scheme)
	}
	
	return true, nil
}

func (r *SchemeRule) Name() string {
	return "SchemeRule"
}

// HostRule validates the URL host
type HostRule struct{}

func (r *HostRule) Apply(u *url.URL) (bool, error) {
	if u.Host == "" {
		return false, errors.New("URL host is missing")
	}
	
	// Validate the host format
	host := u.Hostname()
	if len(host) > 253 {
		return false, errors.New("host name is too long (max 253 characters)")
	}
	
	// Handle IP address hosts specially
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() {
			return false, errors.New("URL with private or loopback IP address is not allowed")
		}
		
		// Additional IP-specific checks could be added here
		return true, nil
	}
	
	// Hostname validation - should follow DNS rules
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return false, errors.New("invalid hostname format")
	}
	
	// Check each part of the domain name
	for _, part := range parts {
		if len(part) == 0 || len(part) > 63 {
			return false, errors.New("each part of hostname must be between 1 and 63 characters")
		}
		
		// Check for valid hostname characters (RFC 1123)
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9]([a-zA-Z0-9\-]*[a-zA-Z0-9])?$`, part)
		if !matched {
			return false, fmt.Errorf("invalid hostname part: %s", part)
		}
	}
	
	return true, nil
}

func (r *HostRule) Name() string {
	return "HostRule"
}

// PathRule validates the URL path
type PathRule struct{}

func (r *PathRule) Apply(u *url.URL) (bool, error) {
	// Path is optional, but if present must be valid
	if u.Path == "" {
		return true, nil
	}
	
	// Check for invalid path characters
	// This is a simplified check - in reality, URL encoding allows more characters
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-\_\.\~\!\$\&\'\(\)\*\+\,\;\=\:\@\/\%]+$`, u.Path)
	if !matched {
		return false, errors.New("path contains invalid characters")
	}
	
	return true, nil
}

func (r *PathRule) Name() string {
	return "PathRule"
}

// QueryParamRule validates the query parameters
type QueryParamRule struct{}

func (r *QueryParamRule) Apply(u *url.URL) (bool, error) {
	// Query parameters are optional
	if u.RawQuery == "" {
		return true, nil
	}
	
	// Check for overly complex query string (might be an attack vector)
	if len(u.RawQuery) > 1000 {
		return false, errors.New("query string too long")
	}
	
	// Check for suspiciously many parameters
	params := u.Query()
	if len(params) > 50 {
		return false, errors.New("too many query parameters")
	}
	
	return true, nil
}

func (r *QueryParamRule) Name() string {
	return "QueryParamRule"
}

// LengthRule validates the overall URL length
type LengthRule struct {
	maxLength int
}

func (r *LengthRule) Apply(u *url.URL) (bool, error) {
	urlStr := u.String()
	if len(urlStr) > r.maxLength {
		return false, fmt.Errorf("URL is too long (max %d characters)", r.maxLength)
	}
	
	return true, nil
}

func (r *LengthRule) Name() string {
	return "LengthRule"
}

// Helper function to check if a domain supports HTTPS
func domainSupportsHTTPS(domain string) (bool, error) {
	// In a production system, this might check the domain's TLS support
	// For now, we'll assume most domains support HTTPS
	
	// Optionally, you could do a real check with a HEAD request:
	// client := &http.Client{Timeout: 5 * time.Second}
	// _, err := client.Head("https://" + domain)
	// return err == nil, nil
	
	// For simplicity, assume common domains support HTTPS
	commonDomains := map[string]bool{
		"google.com":    true,
		"facebook.com":  true,
		"twitter.com":   true,
		"github.com":    true,
		"microsoft.com": true,
		"apple.com":     true,
		"amazon.com":    true,
	}
	
	// Extract the base domain for comparison
	parts := strings.Split(domain, ".")
	if len(parts) >= 2 {
		baseDomain := parts[len(parts)-2] + "." + parts[len(parts)-1]
		if commonDomains[baseDomain] {
			return true, nil
		}
	}
	
	// By default, assume HTTPS is supported for most domains
	return true, nil
}

// IsValidURL is a convenience function for quick validation with default settings
func IsValidURL(rawURL string) bool {
	validator := NewStandardURLValidator()
	return validator.IsValid(rawURL)
}

// ValidateURLWithReason is a convenience function for validation with reason
func ValidateURLWithReason(rawURL string) (bool, string) {
	validator := NewStandardURLValidator()
	valid, err := validator.ValidateWithReason(rawURL)
	if !valid && err != nil {
		return false, err.Error()
	}
	return valid, ""
}
