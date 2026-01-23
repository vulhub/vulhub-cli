package github

import (
	"errors"
	"strings"

	gh "github.com/google/go-github/v68/github"
)

// IsRateLimitError checks if an error is a GitHub rate limit error
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}

	// Check for go-github rate limit error
	var rateLimitErr *gh.RateLimitError
	if errors.As(err, &rateLimitErr) {
		return true
	}

	// Check for go-github abuse rate limit error
	var abuseErr *gh.AbuseRateLimitError
	if errors.As(err, &abuseErr) {
		return true
	}

	// Check error message for rate limit indicators
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "rate limit") ||
		strings.Contains(errStr, "api rate limit exceeded") {
		return true
	}

	return false
}
