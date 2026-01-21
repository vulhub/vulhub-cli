package compose

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/vulhub/vulhub-cli/pkg/types"
)

// Client defines the interface for Docker Compose operations
type Client interface {
	// Start starts the Docker Compose environment
	Start(ctx context.Context, workDir string, options StartOptions) error

	// Stop stops the Docker Compose environment
	Stop(ctx context.Context, workDir string, options StopOptions) error

	// Restart restarts the Docker Compose environment
	Restart(ctx context.Context, workDir string, options RestartOptions) error

	// Status returns the status of containers in the environment
	Status(ctx context.Context, workDir string) ([]types.ContainerStatus, error)

	// Logs retrieves logs from containers
	Logs(ctx context.Context, workDir string, options LogsOptions) (io.ReadCloser, error)

	// Down stops and removes containers, networks, and volumes
	Down(ctx context.Context, workDir string, options DownOptions) error

	// Pull pulls images for the environment
	Pull(ctx context.Context, workDir string) error

	// CheckDocker checks if Docker daemon is running
	CheckDocker(ctx context.Context) error
}

// StartOptions defines options for starting an environment
type StartOptions struct {
	// Detach runs containers in the background
	Detach bool
	// Build builds images before starting
	Build bool
	// ForceRecreate forces recreation of containers
	ForceRecreate bool
	// Pull pulls images before starting
	Pull bool
}

// StopOptions defines options for stopping an environment
type StopOptions struct {
	// Timeout is the timeout in seconds
	Timeout int
}

// RestartOptions defines options for restarting an environment
type RestartOptions struct {
	// Timeout is the timeout in seconds
	Timeout int
}

// LogsOptions defines options for retrieving logs
type LogsOptions struct {
	// Follow follows log output
	Follow bool
	// Tail is the number of lines to show from the end
	Tail string
	// Timestamps shows timestamps
	Timestamps bool
	// Service is the specific service to get logs from
	Service string
}

// DownOptions defines options for down command
type DownOptions struct {
	// RemoveVolumes removes named volumes
	RemoveVolumes bool
	// RemoveImages removes images ("all" or "local")
	RemoveImages string
	// Timeout is the timeout in seconds
	Timeout int
}

// ComposeClient implements the Client interface
type ComposeClient struct {
	executor *Executor
	logger   *slog.Logger
}

// NewComposeClient creates a new ComposeClient
func NewComposeClient(composeCommand string, logger *slog.Logger) *ComposeClient {
	if logger == nil {
		logger = slog.Default()
	}
	return &ComposeClient{
		executor: NewExecutor(composeCommand),
		logger:   logger,
	}
}

// Start starts the Docker Compose environment
func (c *ComposeClient) Start(ctx context.Context, workDir string, options StartOptions) error {
	c.logger.Debug("starting environment", "workDir", workDir)

	args := []string{"up"}

	if options.Detach {
		args = append(args, "-d")
	}
	if options.Build {
		args = append(args, "--build")
	}
	if options.ForceRecreate {
		args = append(args, "--force-recreate")
	}
	if options.Pull {
		args = append(args, "--pull", "always")
	}

	result, err := c.executor.Execute(ctx, workDir, args...)
	if err != nil {
		return fmt.Errorf("failed to start environment: %w", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("failed to start environment: %s", result.Stderr)
	}

	return nil
}

// Stop stops the Docker Compose environment
func (c *ComposeClient) Stop(ctx context.Context, workDir string, options StopOptions) error {
	c.logger.Debug("stopping environment", "workDir", workDir)

	args := []string{"stop"}

	if options.Timeout > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", options.Timeout))
	}

	result, err := c.executor.Execute(ctx, workDir, args...)
	if err != nil {
		return fmt.Errorf("failed to stop environment: %w", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("failed to stop environment: %s", result.Stderr)
	}

	return nil
}

// Restart restarts the Docker Compose environment
func (c *ComposeClient) Restart(ctx context.Context, workDir string, options RestartOptions) error {
	c.logger.Debug("restarting environment", "workDir", workDir)

	args := []string{"restart"}

	if options.Timeout > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", options.Timeout))
	}

	result, err := c.executor.Execute(ctx, workDir, args...)
	if err != nil {
		return fmt.Errorf("failed to restart environment: %w", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("failed to restart environment: %s", result.Stderr)
	}

	return nil
}

// psOutput represents the JSON output from docker compose ps
type psOutput struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	Image   string `json:"Image"`
	Status  string `json:"Status"`
	State   string `json:"State"`
	Service string `json:"Service"`
	Publishers []struct {
		URL           string `json:"URL"`
		TargetPort    int    `json:"TargetPort"`
		PublishedPort int    `json:"PublishedPort"`
		Protocol      string `json:"Protocol"`
	} `json:"Publishers"`
}

// Status returns the status of containers in the environment
func (c *ComposeClient) Status(ctx context.Context, workDir string) ([]types.ContainerStatus, error) {
	c.logger.Debug("getting status", "workDir", workDir)

	args := []string{"ps", "--format", "json"}

	result, err := c.executor.Execute(ctx, workDir, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to get status: %s", result.Stderr)
	}

	// Parse JSON output (each line is a separate JSON object)
	var statuses []types.ContainerStatus
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		var ps psOutput
		if err := json.Unmarshal([]byte(line), &ps); err != nil {
			c.logger.Warn("failed to parse container status", "error", err, "line", line)
			continue
		}

		status := types.ContainerStatus{
			ID:     ps.ID,
			Name:   ps.Name,
			Image:  ps.Image,
			Status: ps.Status,
			State:  ps.State,
		}

		for _, pub := range ps.Publishers {
			port := types.PortMapping{
				HostIP:        pub.URL,
				HostPort:      fmt.Sprintf("%d", pub.PublishedPort),
				ContainerPort: fmt.Sprintf("%d", pub.TargetPort),
				Protocol:      pub.Protocol,
			}
			status.Ports = append(status.Ports, port)
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// Logs retrieves logs from containers
func (c *ComposeClient) Logs(ctx context.Context, workDir string, options LogsOptions) (io.ReadCloser, error) {
	c.logger.Debug("getting logs", "workDir", workDir)

	args := []string{"logs"}

	if options.Follow {
		args = append(args, "-f")
	}
	if options.Tail != "" {
		args = append(args, "--tail", options.Tail)
	}
	if options.Timestamps {
		args = append(args, "-t")
	}
	if options.Service != "" {
		args = append(args, options.Service)
	}

	// For logs, we need to stream output
	cmdParts := strings.Fields(c.executor.composeCommand)
	cmdParts = append(cmdParts, args...)

	cmd := exec.CommandContext(ctx, cmdParts[0], cmdParts[1:]...)
	cmd.Dir = workDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start logs command: %w", err)
	}

	return stdout, nil
}

// Down stops and removes containers, networks, and volumes
func (c *ComposeClient) Down(ctx context.Context, workDir string, options DownOptions) error {
	c.logger.Debug("down environment", "workDir", workDir)

	args := []string{"down"}

	if options.RemoveVolumes {
		args = append(args, "-v")
	}
	if options.RemoveImages != "" {
		args = append(args, "--rmi", options.RemoveImages)
	}
	if options.Timeout > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", options.Timeout))
	}

	result, err := c.executor.Execute(ctx, workDir, args...)
	if err != nil {
		return fmt.Errorf("failed to down environment: %w", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("failed to down environment: %s", result.Stderr)
	}

	return nil
}

// Pull pulls images for the environment
func (c *ComposeClient) Pull(ctx context.Context, workDir string) error {
	c.logger.Debug("pulling images", "workDir", workDir)

	result, err := c.executor.Execute(ctx, workDir, "pull")
	if err != nil {
		return fmt.Errorf("failed to pull images: %w", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("failed to pull images: %s", result.Stderr)
	}

	return nil
}

// CheckDocker checks if Docker daemon is running
func (c *ComposeClient) CheckDocker(ctx context.Context) error {
	if err := c.executor.CheckDockerAvailable(ctx); err != nil {
		return err
	}
	return c.executor.CheckComposeAvailable(ctx)
}
