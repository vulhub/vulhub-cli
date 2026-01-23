package github

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	gh "github.com/google/go-github/v68/github"
	"github.com/vulhub/vulhub-cli/pkg/types"
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

// GitHubClient implements the Client interface using go-github
type GitHubClient struct {
	client *gh.Client
	logger *slog.Logger
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
	var client *gh.Client
	if cfg.Token != "" {
		client = gh.NewClient(nil).WithAuthToken(cfg.Token)
	} else {
		client = gh.NewClient(nil)
	}

	return &GitHubClient{
		client: client,
		logger: logger,
	}
}

// SetToken updates the client's authentication token dynamically.
// This allows refreshing the client after OAuth authentication without recreating the client.
func (c *GitHubClient) SetToken(token string) {
	if token != "" {
		c.client = gh.NewClient(nil).WithAuthToken(token)
	} else {
		c.client = gh.NewClient(nil)
	}
}

// DownloadFile downloads a single file from a repository
func (c *GitHubClient) DownloadFile(ctx context.Context, owner, repo, path, ref string) ([]byte, error) {
	c.logger.Debug("downloading file", "owner", owner, "repo", repo, "path", path, "ref", ref)

	opts := &gh.RepositoryContentGetOptions{Ref: ref}

	fileContent, _, resp, err := c.client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get file content: %w", err)
	}
	closeResponse(resp)

	if fileContent == nil {
		return nil, fmt.Errorf("path %s is not a file", path)
	}

	// Get decoded content
	content, err := fileContent.GetContent()
	if err != nil {
		return nil, fmt.Errorf("failed to decode file content: %w", err)
	}

	return []byte(content), nil
}

// DownloadDirectory downloads all files in a directory from a repository
func (c *GitHubClient) DownloadDirectory(ctx context.Context, owner, repo, path, ref, destDir string) error {
	c.logger.Debug("downloading directory", "path", path, "dest", destDir)

	contents, err := c.ListDirectoryContents(ctx, owner, repo, path, ref)
	if err != nil {
		return err // Error already wrapped with rate limit check
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
				return err // Error already wrapped with rate limit check
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
	release, resp, err := c.client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %w", err)
	}
	closeResponse(resp)

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

	_, dirContents, resp, err := c.client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory contents: %w", err)
	}
	closeResponse(resp)

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

// closeResponse safely closes the response body
func closeResponse(resp *gh.Response) {
	if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
