package github

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	gh "github.com/google/go-github/v68/github"
	"github.com/vulhub/vulhub-cli/pkg/types"
	"resty.dev/v3"
)

const (
	// RawContentBaseURL is the base URL for raw content downloads
	RawContentBaseURL = "https://raw.githubusercontent.com"
)

// Client defines the interface for GitHub operations
type Client interface {
	// DownloadFile downloads a single file from a repository
	DownloadFile(ctx context.Context, owner, repo, path, ref string) ([]byte, error)

	// DownloadDirectory downloads all files in a directory from a repository
	DownloadDirectory(ctx context.Context, owner, repo, path, ref, destDir string) error

	// GetFileContent gets the content of a file (for small files)
	GetFileContent(ctx context.Context, owner, repo, path, ref string) (string, error)

	// GetLatestRelease gets the latest release information
	GetLatestRelease(ctx context.Context, owner, repo string) (*types.Release, error)

	// ListDirectoryContents lists contents of a directory in the repository
	ListDirectoryContents(ctx context.Context, owner, repo, path, ref string) ([]types.ContentEntry, error)
}

// GitHubClient implements the Client interface
type GitHubClient struct {
	ghClient   *gh.Client
	httpClient *resty.Client
	logger     *slog.Logger
	token      string
}

// ClientConfig holds configuration for the GitHub client
type ClientConfig struct {
	Token  string
	Logger *slog.Logger
}

// NewClient creates a new GitHub client
func NewClient(cfg ClientConfig) *GitHubClient {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	// Create GitHub API client
	var ghClient *gh.Client
	if cfg.Token != "" {
		ghClient = gh.NewClient(nil).WithAuthToken(cfg.Token)
	} else {
		ghClient = gh.NewClient(nil)
	}

	// Create HTTP client for raw downloads
	httpClient := resty.New()
	if cfg.Token != "" {
		httpClient.SetHeader("Authorization", "token "+cfg.Token)
	}
	httpClient.SetHeader("Accept", "application/vnd.github.v3.raw")

	return &GitHubClient{
		ghClient:   ghClient,
		httpClient: httpClient,
		logger:     logger,
		token:      cfg.Token,
	}
}

// DownloadFile downloads a single file from a repository
func (c *GitHubClient) DownloadFile(ctx context.Context, owner, repo, path, ref string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s/%s/%s", RawContentBaseURL, owner, repo, ref, path)

	c.logger.Debug("downloading file", "url", url)

	resp, err := c.httpClient.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status %d", resp.StatusCode())
	}

	return resp.Bytes(), nil
}

// DownloadDirectory downloads all files in a directory from a repository
func (c *GitHubClient) DownloadDirectory(ctx context.Context, owner, repo, path, ref, destDir string) error {
	c.logger.Debug("downloading directory", "path", path, "dest", destDir)

	contents, err := c.ListDirectoryContents(ctx, owner, repo, path, ref)
	if err != nil {
		return fmt.Errorf("failed to list directory contents: %w", err)
	}

	for _, entry := range contents {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		entryPath := entry.Path
		destPath := filepath.Join(destDir, entry.Name)

		if entry.Type == "dir" {
			// Recursively download subdirectory
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			if err := c.DownloadDirectory(ctx, owner, repo, entryPath, ref, destPath); err != nil {
				return err
			}
		} else {
			// Download file
			data, err := c.DownloadFile(ctx, owner, repo, entryPath, ref)
			if err != nil {
				return fmt.Errorf("failed to download file %s: %w", entryPath, err)
			}

			if err := os.WriteFile(destPath, data, 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", destPath, err)
			}

			c.logger.Debug("downloaded file", "path", entryPath, "dest", destPath)
		}
	}

	return nil
}

// GetFileContent gets the content of a file (for small files)
func (c *GitHubClient) GetFileContent(ctx context.Context, owner, repo, path, ref string) (string, error) {
	data, err := c.DownloadFile(ctx, owner, repo, path, ref)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetLatestRelease gets the latest release information
func (c *GitHubClient) GetLatestRelease(ctx context.Context, owner, repo string) (*types.Release, error) {
	release, _, err := c.ghClient.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %w", err)
	}

	return &types.Release{
		TagName:     release.GetTagName(),
		Name:        release.GetName(),
		PublishedAt: release.GetPublishedAt().String(),
		Body:        release.GetBody(),
	}, nil
}

// ListDirectoryContents lists contents of a directory in the repository
func (c *GitHubClient) ListDirectoryContents(ctx context.Context, owner, repo, path, ref string) ([]types.ContentEntry, error) {
	opts := &gh.RepositoryContentGetOptions{Ref: ref}

	_, dirContents, resp, err := c.ghClient.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory contents: %w", err)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}()

	entries := make([]types.ContentEntry, 0, len(dirContents))
	for _, content := range dirContents {
		entry := types.ContentEntry{
			Name: content.GetName(),
			Path: content.GetPath(),
			Type: content.GetType(),
			Size: int64(content.GetSize()),
			SHA:  content.GetSHA(),
		}
		if content.DownloadURL != nil {
			entry.DownloadURL = *content.DownloadURL
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
