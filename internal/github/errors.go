package github

import (
	"errors"
	"net/http"
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
		strings.Contains(errStr, "api rate limit exceeded") ||
		strings.Contains(errStr, "403") && strings.Contains(errStr, "limit") {
		return true
	}

	return false
}

// wrapRateLimitError checks if the response indicates a rate limit error
// and wraps it with our custom error type
func wrapRateLimitError(resp *http.Response, err error) error {
	if err == nil {
		return nil
	}

	// Check HTTP status code
	if resp != nil && resp.StatusCode == http.StatusForbidden {
		// GitHub returns 403 for rate limit exceeded
		if IsRateLimitError(err) {
			return ErrRateLimited
		}
	}

	// Check the error itself
	if IsRateLimitError(err) {
		return ErrRateLimited
	}

	return err
}
