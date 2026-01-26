package github

import (
	"os"

	"go.uber.org/fx"

	gh "github.com/google/go-github/v68/github"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/httpclient"
)

// Module provides the github module for fx
var Module = fx.Module("github",
	fx.Provide(
		NewClient,
		NewDownloader,
		func(c *GitHubClient) Client {
			return c
		},
	),
)

// NewClient creates a GitHub client from the config manager
func NewClient(cfgMgr config.Manager, httpClient *httpclient.Client) *GitHubClient {
	cfg := cfgMgr.Get()

	// Check for GITHUB_TOKEN in environment
	token := cfg.GitHub.Token
	if envToken := os.Getenv("GITHUB_TOKEN"); envToken != "" {
		token = envToken
	}

	var client *gh.Client
	if token != "" {
		client = gh.NewClient(httpClient.StandardClient()).WithAuthToken(token)
	} else {
		client = gh.NewClient(httpClient.StandardClient())
	}

	return &GitHubClient{
		client:     client,
		httpClient: httpClient.StandardClient(),
	}
}

// NewDownloader creates a Downloader from the config manager
func NewDownloader(client Client, cfgMgr config.Manager) *Downloader {
	cfg := cfgMgr.Get()
	return &Downloader{
		client: client,
		owner:  cfg.GitHub.Owner,
		repo:   cfg.GitHub.Repo,
		branch: cfg.GitHub.Branch,
	}
}
