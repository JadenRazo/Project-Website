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

type URLValidator interface {
	IsValid(rawURL string) bool
	ValidateWithReason(rawURL string) (bool, error)
	IsSafe(rawURL string) bool
	NormalizeURL(rawURL string) (string, error)
}

type ValidationRules interface {
	Apply(u *url.URL) (bool, error)
	Name() string
}

type StandardURLValidator struct {
	rules         []ValidationRules
	blockedHosts  map[string]bool
	blockedRegexp []*regexp.Regexp
	cache         *ValidationCache
	mu            sync.RWMutex
}

type ValidationCache struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

type CacheItem struct {
	valid      bool
	error      error
	expiration time.Time
}

func NewValidationCache(ttl time.Duration) *ValidationCache {
	return &ValidationCache{
		items: make(map[string]*CacheItem),
		ttl:   ttl,
	}
}

func (c *ValidationCache) Get(key string) (bool, error, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, exists := c.items[key]; exists {
		if time.Now().Before(item.expiration) {
			return item.valid, item.error, true
		}
		delete(c.items, key)
	}
	return false, nil, false
}

func (c *ValidationCache) Set(key string, valid bool, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &CacheItem{
		valid:      valid,
		error:      err,
		expiration: time.Now().Add(c.ttl),
	}

	if len(c.items)%100 == 0 {
		go c.cleanup()
	}
}

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

func NewStandardURLValidator() *StandardURLValidator {
	v := &StandardURLValidator{
		rules: []ValidationRules{
			&SchemeRule{},
			&HostRule{},
			&PathRule{},
			&QueryParamRule{},
			&LengthRule{maxLength: 2048},
		},
		blockedHosts:  make(map[string]bool),
		blockedRegexp: make([]*regexp.Regexp, 0),
		cache:         NewValidationCache(10 * time.Minute),
	}

	v.BlockHost("malware-site.com")
	v.BlockHost("phishing-example.org")
	v.BlockHost("spam-domain.net")

	v.BlockPattern(`\.exe$`)
	v.BlockPattern(`\.zip$`)
	v.BlockPattern(`[0-9a-f]{32}`)

	return v
}

func (v *StandardURLValidator) AddRule(rule ValidationRules) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.rules = append(v.rules, rule)
}

func (v *StandardURLValidator) BlockHost(host string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.blockedHosts[strings.ToLower(host)] = true
}

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

func (v *StandardURLValidator) IsValid(rawURL string) bool {
	isValid, _ := v.ValidateWithReason(rawURL)
	return isValid
}

func (v *StandardURLValidator) ValidateWithReason(rawURL string) (bool, error) {
	if valid, err, found := v.cache.Get(rawURL); found {
		return valid, err
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		v.cache.Set(rawURL, false, fmt.Errorf("failed to parse URL: %w", err))
		return false, fmt.Errorf("failed to parse URL: %w", err)
	}

	host := strings.ToLower(u.Hostname())

	v.mu.RLock()
	if _, blocked := v.blockedHosts[host]; blocked {
		v.mu.RUnlock()
		err := errors.New("URL domain is in blocklist")
		v.cache.Set(rawURL, false, err)
		return false, err
	}
	v.mu.RUnlock()

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

func (v *StandardURLValidator) IsSafe(rawURL string) bool {

	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	host := strings.ToLower(u.Hostname())

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

	if strings.Count(rawURL, "http") > 1 {
		return false
	}

	if len(host) > 50 {
		return false
	}

	if strings.Count(host, ".") > 4 {
		return false
	}

	return true
}

func (v *StandardURLValidator) NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if u.Scheme == "http" {
		supportsHTTPS, _ := domainSupportsHTTPS(u.Hostname())
		if supportsHTTPS {
			u.Scheme = "https"
		}
	}

	q := u.Query()
	trackingParams := []string{"utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content", "fbclid", "gclid"}
	for _, param := range trackingParams {
		q.Del(param)
	}
	u.RawQuery = q.Encode()

	if len(u.Path) > 1 && strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path[:len(u.Path)-1]
	}

	if u.Host != "" {
		u.Host = strings.ToLower(u.Host)
	}

	if u.Scheme == "http" && strings.HasSuffix(u.Host, ":80") {
		u.Host = u.Host[:len(u.Host)-3]
	} else if u.Scheme == "https" && strings.HasSuffix(u.Host, ":443") {
		u.Host = u.Host[:len(u.Host)-4]
	}

	return u.String(), nil
}

type SchemeRule struct{}

func (r *SchemeRule) Apply(u *url.URL) (bool, error) {
	if u.Scheme == "" {
		return false, errors.New("URL scheme is missing")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false, fmt.Errorf("URL scheme must be http or https, got %s", u.Scheme)
	}

	return true, nil
}

func (r *SchemeRule) Name() string {
	return "SchemeRule"
}

type HostRule struct{}

func (r *HostRule) Apply(u *url.URL) (bool, error) {
	if u.Host == "" {
		return false, errors.New("URL host is missing")
	}

	host := u.Hostname()
	if len(host) > 253 {
		return false, errors.New("host name is too long (max 253 characters)")
	}

	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() {
			return false, errors.New("URL with private or loopback IP address is not allowed")
		}

		return true, nil
	}

	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return false, errors.New("invalid hostname format")
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 63 {
			return false, errors.New("each part of hostname must be between 1 and 63 characters")
		}

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

type PathRule struct{}

func (r *PathRule) Apply(u *url.URL) (bool, error) {
	if u.Path == "" {
		return true, nil
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-\_\.\~\!\$\&\'\(\)\*\+\,\;\=\:\@\/\%]+$`, u.Path)
	if !matched {
		return false, errors.New("path contains invalid characters")
	}

	return true, nil
}

func (r *PathRule) Name() string {
	return "PathRule"
}

type QueryParamRule struct{}

func (r *QueryParamRule) Apply(u *url.URL) (bool, error) {
	if u.RawQuery == "" {
		return true, nil
	}

	if len(u.RawQuery) > 1000 {
		return false, errors.New("query string too long")
	}

	params := u.Query()
	if len(params) > 50 {
		return false, errors.New("too many query parameters")
	}

	return true, nil
}

func (r *QueryParamRule) Name() string {
	return "QueryParamRule"
}

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

func domainSupportsHTTPS(domain string) (bool, error) {


	commonDomains := map[string]bool{
		"google.com":    true,
		"facebook.com":  true,
		"twitter.com":   true,
		"github.com":    true,
		"microsoft.com": true,
		"apple.com":     true,
		"amazon.com":    true,
	}

	parts := strings.Split(domain, ".")
	if len(parts) >= 2 {
		baseDomain := parts[len(parts)-2] + "." + parts[len(parts)-1]
		if commonDomains[baseDomain] {
			return true, nil
		}
	}

	return true, nil
}

func IsValidURL(rawURL string) bool {
	validator := NewStandardURLValidator()
	return validator.IsValid(rawURL)
}

func ValidateURLWithReason(rawURL string) (bool, string) {
	validator := NewStandardURLValidator()
	valid, err := validator.ValidateWithReason(rawURL)
	if !valid && err != nil {
		return false, err.Error()
	}
	return valid, ""
}
