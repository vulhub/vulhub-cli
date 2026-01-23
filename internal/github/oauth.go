package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// OAuthClientID is the client ID for the Vulhub CLI OAuth App
	OAuthClientID = "Ov23liDeiHCLOTtZxFY4"

	// GitHub OAuth endpoints
	deviceCodeURL   = "https://github.com/login/device/code"
	accessTokenURL  = "https://github.com/login/oauth/access_token"
	verificationURL = "https://github.com/login/device"

	// OAuth scope needed for public repo access
	oauthScope = "public_repo"
)

// DeviceCodeResponse represents the response from the device code request
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// AccessTokenResponse represents the response from the access token request
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error"`
	ErrorDesc   string `json:"error_description"`
}

// OAuthError represents an OAuth error
type OAuthError struct {
	Code        string
	Description string
}

func (e *OAuthError) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Description)
	}
	return e.Code
}

var (
	// ErrAuthorizationPending indicates the user hasn't completed authorization yet
	ErrAuthorizationPending = errors.New("authorization_pending")

	// ErrSlowDown indicates we should increase the polling interval
	ErrSlowDown = errors.New("slow_down")

	// ErrExpiredToken indicates the device code has expired
	ErrExpiredToken = errors.New("expired_token")

	// ErrAccessDenied indicates the user denied the authorization
	ErrAccessDenied = errors.New("access_denied")
)

// RequestDeviceCode initiates the device flow by requesting a device code
func RequestDeviceCode(ctx context.Context) (*DeviceCodeResponse, error) {
	data := url.Values{}
	data.Set("client_id", OAuthClientID)
	data.Set("scope", oauthScope)

	req, err := http.NewRequestWithContext(ctx, "POST", deviceCodeURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request device code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var deviceCode DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceCode); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Use the standard verification URL if not provided
	if deviceCode.VerificationURI == "" {
		deviceCode.VerificationURI = verificationURL
	}

	// Default interval to 5 seconds if not provided
	if deviceCode.Interval == 0 {
		deviceCode.Interval = 5
	}

	return &deviceCode, nil
}

// PollForAccessToken polls GitHub for the access token
// Returns the access token when authorization is complete
func PollForAccessToken(ctx context.Context, deviceCode string, interval int) (string, error) {
	data := url.Values{}
	data.Set("client_id", OAuthClientID)
	data.Set("device_code", deviceCode)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	req, err := http.NewRequestWithContext(ctx, "POST", accessTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to poll for access token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for errors
	if tokenResp.Error != "" {
		switch tokenResp.Error {
		case "authorization_pending":
			return "", ErrAuthorizationPending
		case "slow_down":
			return "", ErrSlowDown
		case "expired_token":
			return "", ErrExpiredToken
		case "access_denied":
			return "", ErrAccessDenied
		default:
			return "", &OAuthError{Code: tokenResp.Error, Description: tokenResp.ErrorDesc}
		}
	}

	if tokenResp.AccessToken == "" {
		return "", errors.New("no access token in response")
	}

	return tokenResp.AccessToken, nil
}

// DeviceFlowAuth performs the complete device flow authentication
// It returns the access token and calls the provided callbacks for user interaction
type DeviceFlowCallbacks struct {
	// OnUserCode is called when the user code is available
	// The user should be instructed to visit the verification URI and enter the code
	OnUserCode func(userCode, verificationURI string)

	// OnPolling is called each time we poll for the token (for showing progress)
	OnPolling func()
}

// PerformDeviceFlow performs the complete OAuth device flow
func PerformDeviceFlow(ctx context.Context, callbacks DeviceFlowCallbacks) (string, error) {
	// Step 1: Request device code
	deviceCodeResp, err := RequestDeviceCode(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to request device code: %w", err)
	}

	// Step 2: Show user the code
	if callbacks.OnUserCode != nil {
		callbacks.OnUserCode(deviceCodeResp.UserCode, deviceCodeResp.VerificationURI)
	}

	// Step 3: Poll for access token
	interval := time.Duration(deviceCodeResp.Interval) * time.Second
	expiresAt := time.Now().Add(time.Duration(deviceCodeResp.ExpiresIn) * time.Second)

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(interval):
			// Check if expired
			if time.Now().After(expiresAt) {
				return "", ErrExpiredToken
			}

			if callbacks.OnPolling != nil {
				callbacks.OnPolling()
			}

			token, err := PollForAccessToken(ctx, deviceCodeResp.DeviceCode, int(interval.Seconds()))
			if err != nil {
				if errors.Is(err, ErrAuthorizationPending) {
					// Keep polling
					continue
				}
				if errors.Is(err, ErrSlowDown) {
					// Increase interval
					interval += 5 * time.Second
					continue
				}
				return "", err
			}

			return token, nil
		}
	}
}
