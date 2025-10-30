package security

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type IPWhitelist struct {
	allowedIPs    []net.IP
	allowedCIDRs  []*net.IPNet
	localhostOnly bool
	enabled       bool
}

type IPWhitelistConfig struct {
	Enabled       bool     `yaml:"enabled" json:"enabled"`
	LocalhostOnly bool     `yaml:"localhost_only" json:"localhost_only"`
	AllowedIPs    []string `yaml:"allowed_ips" json:"allowed_ips"`
	AllowedCIDRs  []string `yaml:"allowed_cidrs" json:"allowed_cidrs"`
}

func NewIPWhitelist(config IPWhitelistConfig) (*IPWhitelist, error) {
	wl := &IPWhitelist{
		enabled:       config.Enabled,
		localhostOnly: config.LocalhostOnly,
	}

	if !config.Enabled {
		return wl, nil
	}

	// Parse individual IPs
	for _, ipStr := range config.AllowedIPs {
		ip := net.ParseIP(strings.TrimSpace(ipStr))
		if ip == nil {
			return nil, fmt.Errorf("invalid IP address: %s", ipStr)
		}
		wl.allowedIPs = append(wl.allowedIPs, ip)
	}

	// Parse CIDR ranges
	for _, cidrStr := range config.AllowedCIDRs {
		_, ipNet, err := net.ParseCIDR(strings.TrimSpace(cidrStr))
		if err != nil {
			return nil, fmt.Errorf("invalid CIDR range: %s", cidrStr)
		}
		wl.allowedCIDRs = append(wl.allowedCIDRs, ipNet)
	}

	return wl, nil
}

func (wl *IPWhitelist) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !wl.enabled {
			c.Next()
			return
		}

		clientIP := getClientIP(c)
		if !wl.isAllowed(clientIP) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied: IP address not whitelisted",
				"code":  "IP_NOT_WHITELISTED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (wl *IPWhitelist) isAllowed(ipStr string) bool {
	if !wl.enabled {
		return true
	}

	clientIP := net.ParseIP(ipStr)
	if clientIP == nil {
		return false
	}

	// If localhost-only mode is enabled, only allow loopback addresses
	if wl.localhostOnly {
		return clientIP.IsLoopback()
	}

	// Always allow loopback addresses
	if clientIP.IsLoopback() {
		return true
	}

	// Check against individual allowed IPs
	for _, allowedIP := range wl.allowedIPs {
		if clientIP.Equal(allowedIP) {
			return true
		}
	}

	// Check against allowed CIDR ranges
	for _, allowedCIDR := range wl.allowedCIDRs {
		if allowedCIDR.Contains(clientIP) {
			return true
		}
	}

	return false
}

func getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header first (most common proxy header)
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, use the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xrip := c.GetHeader("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}

	// Check CF-Connecting-IP (Cloudflare)
	if cfip := c.GetHeader("CF-Connecting-IP"); cfip != "" {
		return strings.TrimSpace(cfip)
	}

	// Fall back to RemoteAddr
	clientIP, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return clientIP
}

func DefaultLocalhostConfig() IPWhitelistConfig {
	return IPWhitelistConfig{
		Enabled:       true,
		LocalhostOnly: true,
		AllowedIPs:    []string{},
		AllowedCIDRs:  []string{},
	}
}

func DefaultDevelopmentConfig() IPWhitelistConfig {
	return IPWhitelistConfig{
		Enabled:       true,
		LocalhostOnly: false,
		AllowedIPs:    []string{"127.0.0.1", "::1"},
		AllowedCIDRs:  []string{"192.168.0.0/16", "10.0.0.0/8", "172.16.0.0/12"},
	}
}

func (wl *IPWhitelist) IsEnabled() bool {
	return wl.enabled
}

func (wl *IPWhitelist) IsLocalhostOnly() bool {
	return wl.localhostOnly
}

func (wl *IPWhitelist) GetAllowedIPs() []string {
	var ips []string
	for _, ip := range wl.allowedIPs {
		ips = append(ips, ip.String())
	}
	return ips
}

func (wl *IPWhitelist) GetAllowedCIDRs() []string {
	var cidrs []string
	for _, cidr := range wl.allowedCIDRs {
		cidrs = append(cidrs, cidr.String())
	}
	return cidrs
}
