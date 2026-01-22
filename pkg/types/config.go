package types

import "time"

// Config represents the main configuration for vulhub-cli
type Config struct {
	// Version is the configuration version
	Version string `toml:"version"`

	// GitHub contains GitHub-related configuration
	GitHub GitHubConfig `toml:"github"`

	// Paths contains path-related configuration
	Paths PathsConfig `toml:"paths"`

	// Docker contains Docker-related configuration
	Docker DockerConfig `toml:"docker"`

	// Sync contains sync-related configuration
	Sync SyncConfig `toml:"sync"`
}

// SyncConfig contains sync-related configuration
type SyncConfig struct {
	// LastSyncTime is the last time environments were synced
	LastSyncTime time.Time `toml:"last_sync_time,omitempty"`

	// AutoSyncDays is the number of days after which to prompt for sync (default: 7)
	AutoSyncDays int `toml:"auto_sync_days,omitempty"`
}

// GitHubConfig contains GitHub-related configuration
type GitHubConfig struct {
	// Owner is the GitHub repository owner
	Owner string `toml:"owner"`

	// Repo is the GitHub repository name
	Repo string `toml:"repo"`

	// Branch is the default branch to use
	Branch string `toml:"branch"`

	// Token is the optional GitHub token for authentication
	// Can be overridden by GITHUB_TOKEN environment variable
	Token string `toml:"token,omitempty"`
}

// PathsConfig contains path-related configuration
type PathsConfig struct {
	// ConfigDir is the configuration directory (default: ~/.vulhub)
	ConfigDir string `toml:"config_dir,omitempty"`

	// EnvironmentsDir is the directory for downloaded environments
	EnvironmentsDir string `toml:"environments_dir,omitempty"`
}

// DockerConfig contains Docker-related configuration
type DockerConfig struct {
	// ComposeCommand is the docker compose command to use (default: "docker compose")
	ComposeCommand string `toml:"compose_command,omitempty"`

	// Timeout is the default timeout for Docker operations in seconds
	Timeout int `toml:"timeout,omitempty"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Version: "1.0",
		GitHub: GitHubConfig{
			Owner:  "vulhub",
			Repo:   "vulhub",
			Branch: "master",
		},
		Docker: DockerConfig{
			ComposeCommand: "docker compose",
			Timeout:        300,
		},
		Sync: SyncConfig{
			AutoSyncDays: 7,
		},
	}
}

// Release represents a GitHub release
type Release struct {
	// TagName is the tag name of the release
	TagName string

	// Name is the name of the release
	Name string

	// PublishedAt is the publication date
	PublishedAt string

	// Body is the release description
	Body string
}

// ContentEntry represents a file or directory entry in a GitHub repository
type ContentEntry struct {
	// Name is the name of the entry
	Name string

	// Path is the full path of the entry
	Path string

	// Type is the type of the entry ("file" or "dir")
	Type string

	// Size is the size of the file (0 for directories)
	Size int64

	// SHA is the SHA hash of the content
	SHA string

	// DownloadURL is the URL to download the file
	DownloadURL string
}
