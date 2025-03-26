package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// performs URL validation
func ValidateURL(rawURL string) error {
	// Basic checks
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Parse the URL
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}

	// Scheme validation
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("only http and https URLs are allowed")
	}

	// Host validation
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must contain a host")
	}

	// IP address validation (block private/local IPs)
	if isPrivateIP(parsedURL.Hostname()) {
		return fmt.Errorf("URL cannot point to private/local IP addresses")
	}

	// Malicious pattern detection
	if hasMaliciousPattern(parsedURL) {
		return fmt.Errorf("URL contains potentially malicious patterns")
	}

	// Reserved TLDs check
	if isReservedTLD(parsedURL.Hostname()) {
		return fmt.Errorf("URL uses a reserved or special-use domain")
	}

	return nil
}

// checks if the host is a private IP address
func isPrivateIP(host string) bool {
	privateIPBlocks := []string{
		"10.", "172.16.", "172.17.", "172.18.", "172.19.",
		"172.20.", "172.21.", "172.22.", "172.23.", "172.24.",
		"172.25.", "172.26.", "172.27.", "172.28.", "172.29.",
		"172.30.", "172.31.", "192.168.", "127.", "0.", "169.254.",
		"localhost", "::1",
	}

	for _, block := range privateIPBlocks {
		if strings.HasPrefix(host, block) {
			return true
		}
	}
	return false
}

// checks for common malicious URL patterns
func hasMaliciousPattern(u *url.URL) bool {
	maliciousPatterns := []string{
		"javascript:", "data:", "vbscript:", "alert(", "eval(",
		"document.cookie", "window.location", "<script>",
		"%3Cscript%3E", // URL-encoded script tag
	}

	urlString := u.String()
	for _, pattern := range maliciousPatterns {
		if strings.Contains(strings.ToLower(urlString), strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// checks for reserved/special-use domains
func isReservedTLD(host string) bool {
	reservedTLDs := []string{
		".test", ".example", ".invalid", ".localhost",
	}

	for _, tld := range reservedTLDs {
		if strings.HasSuffix(host, tld) {
			return true
		}
	}
	return false
}
