package environment

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/vulhub/vulhub-cli/internal/compose"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/github"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// StartOptions defines options for starting an environment
type StartOptions struct {
	// Pull pulls images before starting
	Pull bool
	// Build builds images before starting
	Build bool
	// ForceRecreate forces recreation of containers
	ForceRecreate bool
}

// CleanOptions defines options for cleaning an environment
type CleanOptions struct {
	// RemoveVolumes removes volumes
	RemoveVolumes bool
	// RemoveImages removes images
	RemoveImages bool
	// RemoveFiles removes local files
	RemoveFiles bool
}

// Manager defines the interface for environment management
type Manager interface {
	// Start starts an environment
	Start(ctx context.Context, env types.Environment, options StartOptions) error

	// Stop stops an environment
	Stop(ctx context.Context, env types.Environment) error

	// Down stops and removes containers, networks, and volumes for an environment
	Down(ctx context.Context, env types.Environment) error

	// Restart restarts an environment
	Restart(ctx context.Context, env types.Environment) error

	// Status gets the status of an environment
	Status(ctx context.Context, env types.Environment) (*types.EnvironmentStatus, error)

	// Clean cleans up an environment
	Clean(ctx context.Context, env types.Environment, options CleanOptions) error

	// ListRunning lists all running environments
	ListRunning(ctx context.Context) ([]types.EnvironmentStatus, error)

	// ListDownloaded lists all downloaded environments with their status
	ListDownloaded(ctx context.Context) ([]types.EnvironmentStatus, error)

	// GetInfo gets detailed information about an environment
	GetInfo(ctx context.Context, env types.Environment) (*types.EnvironmentInfo, error)

	// EnsureDownloaded ensures an environment is downloaded
	EnsureDownloaded(ctx context.Context, env types.Environment) (string, error)

	// IsDownloaded checks if an environment is downloaded
	IsDownloaded(env types.Environment) bool
}

// EnvironmentManager implements the Manager interface
type EnvironmentManager struct {
	configMgr   config.Manager
	composeClient compose.Client
	downloader  *github.Downloader
	logger      *slog.Logger
}

// NewEnvironmentManager creates a new EnvironmentManager
func NewEnvironmentManager(
	configMgr config.Manager,
	composeClient compose.Client,
	downloader *github.Downloader,
	logger *slog.Logger,
) *EnvironmentManager {
	if logger == nil {
		logger = slog.Default()
	}
	return &EnvironmentManager{
		configMgr:   configMgr,
		composeClient: composeClient,
		downloader:  downloader,
		logger:      logger,
	}
}

// Start starts an environment
func (m *EnvironmentManager) Start(ctx context.Context, env types.Environment, options StartOptions) error {
	// Check Docker availability
	if err := m.composeClient.CheckDocker(ctx); err != nil {
		return fmt.Errorf("docker check failed: %w", err)
	}

	// Ensure environment is downloaded
	workDir, err := m.EnsureDownloaded(ctx, env)
	if err != nil {
		return fmt.Errorf("failed to download environment: %w", err)
	}

	m.logger.Info("starting environment", "path", env.Path, "workDir", workDir)

	// Start the environment
	startOpts := compose.StartOptions{
		Detach:        true,
		Pull:          options.Pull,
		Build:         options.Build,
		ForceRecreate: options.ForceRecreate,
	}

	if err := m.composeClient.Start(ctx, workDir, startOpts); err != nil {
		return fmt.Errorf("failed to start environment: %w", err)
	}

	return nil
}

// Stop stops an environment
func (m *EnvironmentManager) Stop(ctx context.Context, env types.Environment) error {
	workDir := m.configMgr.Paths().EnvironmentDir(env.Path)

	if !m.IsDownloaded(env) {
		return fmt.Errorf("environment %s is not downloaded", env.Path)
	}

	m.logger.Info("stopping environment", "path", env.Path)

	if err := m.composeClient.Stop(ctx, workDir, compose.StopOptions{}); err != nil {
		return fmt.Errorf("failed to stop environment: %w", err)
	}

	return nil
}

// Down stops and removes containers, networks, and volumes for an environment,
// and also removes the local environment files
func (m *EnvironmentManager) Down(ctx context.Context, env types.Environment) error {
	workDir := m.configMgr.Paths().EnvironmentDir(env.Path)

	if !m.IsDownloaded(env) {
		return fmt.Errorf("environment %s is not downloaded", env.Path)
	}

	m.logger.Info("downing environment", "path", env.Path)

	downOpts := compose.DownOptions{
		RemoveVolumes: true,
	}

	if err := m.composeClient.Down(ctx, workDir, downOpts); err != nil {
		return fmt.Errorf("failed to down environment: %w", err)
	}

	// Remove local environment files
	if err := os.RemoveAll(workDir); err != nil {
		return fmt.Errorf("failed to remove environment files: %w", err)
	}

	return nil
}

// Restart restarts an environment
func (m *EnvironmentManager) Restart(ctx context.Context, env types.Environment) error {
	workDir := m.configMgr.Paths().EnvironmentDir(env.Path)

	if !m.IsDownloaded(env) {
		return fmt.Errorf("environment %s is not downloaded", env.Path)
	}

	m.logger.Info("restarting environment", "path", env.Path)

	if err := m.composeClient.Restart(ctx, workDir, compose.RestartOptions{}); err != nil {
		return fmt.Errorf("failed to restart environment: %w", err)
	}

	return nil
}

// Status gets the status of an environment
func (m *EnvironmentManager) Status(ctx context.Context, env types.Environment) (*types.EnvironmentStatus, error) {
	status := &types.EnvironmentStatus{
		Environment: env,
		LocalPath:   m.configMgr.Paths().EnvironmentDir(env.Path),
	}

	if !m.IsDownloaded(env) {
		return status, nil
	}

	containers, err := m.composeClient.Status(ctx, status.LocalPath)
	if err != nil {
		// Not an error if no containers exist
		m.logger.Debug("failed to get container status", "error", err)
		return status, nil
	}

	status.Containers = containers
	status.Running = lo.SomeBy(containers, func(c types.ContainerStatus) bool {
		return c.State == "running"
	})

	return status, nil
}

// Clean cleans up an environment
func (m *EnvironmentManager) Clean(ctx context.Context, env types.Environment, options CleanOptions) error {
	workDir := m.configMgr.Paths().EnvironmentDir(env.Path)

	m.logger.Info("cleaning environment", "path", env.Path)

	// Run docker compose down if environment exists
	if m.IsDownloaded(env) {
		downOpts := compose.DownOptions{
			RemoveVolumes: options.RemoveVolumes,
		}
		if options.RemoveImages {
			downOpts.RemoveImages = "local"
		}

		if err := m.composeClient.Down(ctx, workDir, downOpts); err != nil {
			m.logger.Warn("failed to run docker compose down", "error", err)
		}
	}

	// Remove local files if requested
	if options.RemoveFiles {
		if err := os.RemoveAll(workDir); err != nil {
			return fmt.Errorf("failed to remove environment files: %w", err)
		}
	}

	return nil
}

// ListRunning lists all running environments
func (m *EnvironmentManager) ListRunning(ctx context.Context) ([]types.EnvironmentStatus, error) {
	// Load environment list
	envList, err := m.configMgr.LoadEnvironments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load environments: %w", err)
	}

	var statuses []types.EnvironmentStatus

	// Check each downloaded environment
	envsDir := m.configMgr.Paths().EnvironmentsDir()
	if _, err := os.Stat(envsDir); os.IsNotExist(err) {
		return statuses, nil
	}

	for _, env := range envList.Environment {
		if !m.IsDownloaded(env) {
			continue
		}

		status, err := m.Status(ctx, env)
		if err != nil {
			m.logger.Debug("failed to get status", "path", env.Path, "error", err)
			continue
		}

		if status.Running {
			statuses = append(statuses, *status)
		}
	}

	return statuses, nil
}

// ListDownloaded lists all downloaded environments with their status
func (m *EnvironmentManager) ListDownloaded(ctx context.Context) ([]types.EnvironmentStatus, error) {
	// Load environment list
	envList, err := m.configMgr.LoadEnvironments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load environments: %w", err)
	}

	var statuses []types.EnvironmentStatus

	// Check each downloaded environment
	envsDir := m.configMgr.Paths().EnvironmentsDir()
	if _, err := os.Stat(envsDir); os.IsNotExist(err) {
		return statuses, nil
	}

	for _, env := range envList.Environment {
		if !m.IsDownloaded(env) {
			continue
		}

		status, err := m.Status(ctx, env)
		if err != nil {
			m.logger.Debug("failed to get status", "path", env.Path, "error", err)
			continue
		}

		statuses = append(statuses, *status)
	}

	return statuses, nil
}

// GetInfo gets detailed information about an environment
func (m *EnvironmentManager) GetInfo(ctx context.Context, env types.Environment) (*types.EnvironmentInfo, error) {
	info := &types.EnvironmentInfo{
		Environment: env,
		Downloaded:  m.IsDownloaded(env),
		LocalPath:   m.configMgr.Paths().EnvironmentDir(env.Path),
	}

	// Try to get README content
	readme, err := m.downloader.GetEnvironmentReadme(ctx, env)
	if err == nil {
		info.Readme = readme
	}

	// Try to get compose file content
	compose, err := m.downloader.GetEnvironmentComposeFile(ctx, env)
	if err == nil {
		info.ComposeFile = compose
	}

	return info, nil
}

// EnsureDownloaded ensures an environment is downloaded
func (m *EnvironmentManager) EnsureDownloaded(ctx context.Context, env types.Environment) (string, error) {
	destDir := m.configMgr.Paths().EnvironmentDir(env.Path)

	// Check if already downloaded
	composeFile := filepath.Join(destDir, "docker-compose.yml")
	if _, err := os.Stat(composeFile); err == nil {
		m.logger.Debug("environment already downloaded", "path", env.Path)
		return destDir, nil
	}

	m.logger.Info("downloading environment", "path", env.Path)

	// Create directory
	if err := m.configMgr.Paths().EnsureEnvironmentDir(env.Path); err != nil {
		return "", fmt.Errorf("failed to create environment directory: %w", err)
	}

	// Download environment files
	if err := m.downloader.DownloadEnvironment(ctx, env, destDir); err != nil {
		return "", fmt.Errorf("failed to download environment: %w", err)
	}

	return destDir, nil
}

// IsDownloaded checks if an environment is downloaded
func (m *EnvironmentManager) IsDownloaded(env types.Environment) bool {
	destDir := m.configMgr.Paths().EnvironmentDir(env.Path)
	composeFile := filepath.Join(destDir, "docker-compose.yml")
	_, err := os.Stat(composeFile)
	return err == nil
}
