package github

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vulhub/vulhub-cli/pkg/types"
)

// Downloader provides high-level download operations
type Downloader struct {
	client Client
	owner  string
	repo   string
	branch string
}

// DownloadEnvironmentsList downloads the environments.toml file
func (d *Downloader) DownloadEnvironmentsList(ctx context.Context) ([]byte, error) {
	return d.client.DownloadFile(ctx, d.owner, d.repo, "environments.toml", d.branch)
}

// DownloadEnvironment downloads all files for a specific environment
func (d *Downloader) DownloadEnvironment(ctx context.Context, env types.Environment, destDir string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Download the environment directory
	return d.client.DownloadDirectory(ctx, d.owner, d.repo, env.Path, d.branch, destDir)
}

// DownloadEnvironmentFile downloads a specific file from an environment
func (d *Downloader) DownloadEnvironmentFile(ctx context.Context, env types.Environment, filename string) ([]byte, error) {
	path := filepath.Join(env.Path, filename)
	// Convert Windows path separators to URL path separators
	path = filepath.ToSlash(path)
	return d.client.DownloadFile(ctx, d.owner, d.repo, path, d.branch)
}

// GetEnvironmentReadme gets the README content for an environment
func (d *Downloader) GetEnvironmentReadme(ctx context.Context, env types.Environment) (string, error) {
	// Try different README file names
	readmeNames := []string{"README.md", "README.zh-cn.md", "readme.md"}

	for _, name := range readmeNames {
		path := filepath.Join(env.Path, name)
		path = filepath.ToSlash(path)
		content, err := d.client.GetFileContent(ctx, d.owner, d.repo, path, d.branch)
		if err == nil {
			return content, nil
		}
	}

	return "", fmt.Errorf("README not found for environment %s", env.Path)
}

// GetEnvironmentComposeFile gets the docker-compose.yml content for an environment
func (d *Downloader) GetEnvironmentComposeFile(ctx context.Context, env types.Environment) (string, error) {
	// Try different compose file names
	composeNames := []string{"docker-compose.yml", "docker-compose.yaml"}

	for _, name := range composeNames {
		path := filepath.Join(env.Path, name)
		path = filepath.ToSlash(path)
		content, err := d.client.GetFileContent(ctx, d.owner, d.repo, path, d.branch)
		if err == nil {
			return content, nil
		}
	}

	return "", fmt.Errorf("docker-compose.yml not found for environment %s", env.Path)
}

// CheckEnvironmentExists checks if an environment exists in the repository
func (d *Downloader) CheckEnvironmentExists(ctx context.Context, envPath string) (bool, error) {
	_, err := d.client.ListDirectoryContents(ctx, d.owner, d.repo, envPath, d.branch)
	if err != nil {
		return false, nil
	}
	return true, nil
}
