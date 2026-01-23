package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/internal/github"
)

// GitHubAuth creates the github-auth command for GitHub token configuration
func (c *Commands) GitHubAuth() *cli.Command {
	return &cli.Command{
		Name:    "github-auth",
		Usage:   "Authenticate with GitHub to avoid API rate limits",
		Aliases: []string{"ga"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "status",
				Usage: "Show current authentication status",
			},
			&cli.BoolFlag{
				Name:  "remove",
				Usage: "Remove saved GitHub authentication",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Bool("status") {
				return c.runAuthStatus()
			}
			if cmd.Bool("remove") {
				return c.runAuthRemove(ctx)
			}
			return c.runAuthSetup(ctx)
		},
	}
}

// runAuthStatus shows the current authentication status
func (c *Commands) runAuthStatus() error {
	cfg := c.Config.Get()

	fmt.Println()
	if cfg.GitHub.Token != "" {
		card := &ui.StatusCard{
			Icon:    "●",
			Title:   "Authenticated",
			Details: []string{fmt.Sprintf("Token: %s", maskToken(cfg.GitHub.Token))},
			Style:   ui.StatusSuccess,
		}
		card.Print()
		ui.PrintInfo("API rate limit: 5,000 requests/hour")
	} else {
		card := &ui.StatusCard{
			Icon:  "○",
			Title: "Not authenticated",
			Style: ui.StatusWarning,
		}
		card.Print()
		ui.PrintInfo("API rate limit: 60 requests/hour (unauthenticated)")
		fmt.Println()
		ui.PrintInfo("Run 'vulhub github-auth' to authenticate.")
	}
	fmt.Println()

	return nil
}

// runAuthRemove removes the saved GitHub token
func (c *Commands) runAuthRemove(ctx context.Context) error {
	cfg := c.Config.Get()
	if cfg.GitHub.Token == "" {
		fmt.Println()
		ui.PrintInfo("Not authenticated with GitHub")
		fmt.Println()
		return nil
	}

	confirmed, err := ui.ConfirmPrompt(
		"Remove GitHub authentication?",
		"You will need to re-authenticate to avoid rate limits.",
		"Yes, remove",
		"Cancel",
	)
	if err != nil || !confirmed {
		ui.PrintInfo("Cancelled")
		return nil
	}

	cfg.GitHub.Token = ""
	c.Config.Set(cfg)
	if err := c.Config.Save(ctx); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println()
	ui.PrintSuccess("GitHub authentication removed")
	fmt.Println()
	return nil
}

// runAuthSetup performs GitHub OAuth device flow authentication
func (c *Commands) runAuthSetup(ctx context.Context) error {
	// Check if already authenticated
	cfg := c.Config.Get()
	if cfg.GitHub.Token != "" {
		confirmed, err := ui.ConfirmPrompt(
			"Already authenticated with GitHub",
			"Do you want to re-authenticate?",
			"Yes, re-authenticate",
			"Cancel",
		)
		if err != nil || !confirmed {
			ui.PrintInfo("Cancelled")
			return nil
		}
	}

	// Welcome message
	ui.PrintTitle("GitHub Authentication")
	ui.PrintMuted("Authenticate with GitHub to increase API rate limit from 60 to 5,000 requests/hour.")
	fmt.Println()

	// Perform device flow
	token, err := c.performDeviceFlow(ctx)
	if err != nil {
		if errors.Is(err, github.ErrAccessDenied) {
			fmt.Println()
			ui.PrintWarning("Authorization was denied")
			fmt.Println()
			return nil
		}
		if errors.Is(err, github.ErrExpiredToken) {
			fmt.Println()
			ui.PrintWarning("Authorization timed out")
			ui.PrintInfo("Please try again.")
			fmt.Println()
			return nil
		}
		if errors.Is(err, context.Canceled) {
			fmt.Println()
			ui.PrintInfo("Cancelled")
			fmt.Println()
			return nil
		}
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Save token
	cfg = c.Config.Get()
	cfg.GitHub.Token = token
	c.Config.Set(cfg)

	if err := c.Config.Save(ctx); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// Success message
	fmt.Println()
	msg := &ui.MessageBox{
		Title:   "✓ Authentication successful!",
		Message: "API rate limit increased to 5,000 requests/hour",
		Style:   ui.StatusSuccess,
	}
	msg.Print()
	fmt.Println()

	return nil
}

// performDeviceFlow performs the OAuth device flow and returns the access token
func (c *Commands) performDeviceFlow(ctx context.Context) (string, error) {
	// Request device code with spinner
	var deviceCodeResp *github.DeviceCodeResponse
	var requestErr error

	err := ui.RunWithSpinner("Requesting authorization code...", func() {
		deviceCodeResp, requestErr = github.RequestDeviceCode(ctx)
	})

	if err != nil {
		return "", err
	}
	if requestErr != nil {
		return "", fmt.Errorf("failed to request device code: %w", requestErr)
	}

	// Display user code
	authCard := &ui.AuthCodeCard{
		Instruction: "Please visit the URL below and enter the code:",
		URL:         deviceCodeResp.VerificationURI,
		Code:        deviceCodeResp.UserCode,
	}
	authCard.Print()

	// Try to open browser
	if err := ui.OpenBrowser(deviceCodeResp.VerificationURI); err == nil {
		ui.PrintInfo("Browser opened automatically")
	}
	fmt.Println()

	// Poll for access token with spinner
	var token string
	var pollErr error
	interval := time.Duration(deviceCodeResp.Interval) * time.Second
	expiresAt := time.Now().Add(time.Duration(deviceCodeResp.ExpiresIn) * time.Second)

	err = ui.RunWithSpinner("Waiting for authorization...", func() {
		for {
			select {
			case <-ctx.Done():
				pollErr = ctx.Err()
				return
			case <-time.After(interval):
				if time.Now().After(expiresAt) {
					pollErr = github.ErrExpiredToken
					return
				}

				t, err := github.PollForAccessToken(ctx, deviceCodeResp.DeviceCode, int(interval.Seconds()))
				if err != nil {
					if errors.Is(err, github.ErrAuthorizationPending) {
						continue
					}
					if errors.Is(err, github.ErrSlowDown) {
						interval += 5 * time.Second
						continue
					}
					pollErr = err
					return
				}

				token = t
				return
			}
		}
	})

	if err != nil {
		return "", err
	}
	if pollErr != nil {
		return "", pollErr
	}

	return token, nil
}

// PromptTokenSetup prompts the user to set up GitHub authentication (for use after rate limit errors)
// Returns true if the user set up authentication, false otherwise
func (c *Commands) PromptTokenSetup(ctx context.Context) bool {
	fmt.Println()
	msg := &ui.MessageBox{
		Title:   "⚠ GitHub API rate limit exceeded",
		Message: "Authenticate to increase limit from 60 to 5,000 requests/hour",
		Style:   ui.StatusWarning,
	}
	msg.Print()

	confirmed, err := ui.ConfirmPrompt(
		"Authenticate with GitHub now?",
		"",
		"Yes",
		"Later",
	)
	if err != nil || !confirmed {
		fmt.Println()
		ui.PrintInfo("Run 'vulhub github-auth' to authenticate later.")
		fmt.Println()
		return false
	}

	err = c.runAuthSetup(ctx)
	return err == nil && c.Config.Get().GitHub.Token != ""
}

// maskToken masks a token for display, showing only first and last 4 characters
func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "····" + token[len(token)-4:]
}
