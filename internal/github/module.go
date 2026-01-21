package github

import (
	"log/slog"
	"os"

	"go.uber.org/fx"

	"github.com/vulhub/vulhub-cli/internal/config"
)

// Module provides the github module for fx
var Module = fx.Module("github",
	fx.Provide(
		NewClientFromConfig,
		NewDownloaderFromConfig,
		func(c *GitHubClient) Client {
			return c
		},
	),
)

// NewClientFromConfig creates a GitHub client from the config manager
func NewClientFromConfig(cfgMgr config.Manager, logger *slog.Logger) *GitHubClient {
	cfg := cfgMgr.Get()

	// Check for GITHUB_TOKEN in environment
	token := cfg.GitHub.Token
	if envToken := os.Getenv("GITHUB_TOKEN"); envToken != "" {
		token = envToken
	}

	return NewClient(ClientConfig{
		Token:  token,
		Logger: logger,
	})
}

// NewDownloaderFromConfig creates a Downloader from the config manager
func NewDownloaderFromConfig(client Client, cfgMgr config.Manager) *Downloader {
	cfg := cfgMgr.Get()
	return NewDownloader(client, cfg.GitHub.Owner, cfg.GitHub.Repo, cfg.GitHub.Branch)
}
