package github

import (
	"errors"
	"strings"

	gh "github.com/google/go-github/v68/github"
)

// ErrRateLimited is returned when the GitHub API rate limit is exceeded
var ErrRateLimited = errors.New("GitHub API rate limit exceeded")

// IsRateLimitError checks if an error is a GitHub rate limit error
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's our custom error
	if errors.Is(err, ErrRateLimited) {
		return true
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

// checkRateLimitError checks if the error is a rate limit error and returns ErrRateLimited if so.
// This function is used to wrap API call errors with a consistent rate limit error.
func checkRateLimitError(resp *gh.Response, err error) error {
	if err == nil {
		return nil
	}

	// Check for go-github rate limit errors directly
	var rateLimitErr *gh.RateLimitError
	if errors.As(err, &rateLimitErr) {
		return ErrRateLimited
	}

	var abuseErr *gh.AbuseRateLimitError
	if errors.As(err, &abuseErr) {
		return ErrRateLimited
	}

	// Check response status code (403 often indicates rate limit)
	if resp != nil && resp.StatusCode == 403 {
		// Check if this is specifically a rate limit error
		if IsRateLimitError(err) {
			return ErrRateLimited
		}
	}

	// Check error message as fallback
	if IsRateLimitError(err) {
		return ErrRateLimited
	}

	// Not a rate limit error, return nil to indicate caller should handle original error
	return nil
}
