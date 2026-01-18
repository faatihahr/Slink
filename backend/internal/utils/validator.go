package utils

import (
	"net/url"
	"regexp"
	"strings"
)

func IsValidURL(input string) bool {
	if input == "" {
		return false
	}

	// Add scheme if missing for parsing
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		return false
	}

	parsedURL, err := url.Parse(input)
	if err != nil {
		return false
	}

	// Check if scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// Check if host exists
	if parsedURL.Host == "" {
		return false
	}

	// Basic domain validation
	domainRegex := `^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`
	host := parsedURL.Hostname()
	if host == "" || !regexp.MustCompile(domainRegex).MatchString(host) {
		return false
	}

	return true
}

func IsValidCustomAlias(alias string) bool {
	if len(alias) == 0 || len(alias) > 10 {
		return false
	}

	// Only allow alphanumeric characters and hyphens
	validAliasRegex := `^[a-zA-Z0-9\-]+$`
	return regexp.MustCompile(validAliasRegex).MatchString(alias)
}

func GenerateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	result := make([]byte, length)
	for i := range result {
		// Simple pseudo-random generation
		result[i] = charset[i%len(charset)]
	}

	return string(result)
}
