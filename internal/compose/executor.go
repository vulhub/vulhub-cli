package compose

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Executor executes docker compose commands
type Executor struct {
	composeCommand string
}

// NewExecutor creates a new Executor
func NewExecutor(composeCommand string) *Executor {
	if composeCommand == "" {
		composeCommand = "docker compose"
	}
	return &Executor{
		composeCommand: composeCommand,
	}
}

// ExecResult represents the result of a command execution
type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Execute executes a docker compose command
func (e *Executor) Execute(ctx context.Context, workDir string, args ...string) (*ExecResult, error) {
	// Split the compose command (e.g., "docker compose" -> ["docker", "compose"])
	cmdParts := strings.Fields(e.composeCommand)
	cmdParts = append(cmdParts, args...)

	cmd := exec.CommandContext(ctx, cmdParts[0], cmdParts[1:]...)
	cmd.Dir = workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &ExecResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
	} else if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}

	return result, nil
}

// ExecuteWithOutput executes a docker compose command and returns combined output
func (e *Executor) ExecuteWithOutput(ctx context.Context, workDir string, args ...string) (string, error) {
	result, err := e.Execute(ctx, workDir, args...)
	if err != nil {
		return "", err
	}

	if result.ExitCode != 0 {
		return "", fmt.Errorf("command failed with exit code %d: %s", result.ExitCode, result.Stderr)
	}

	return result.Stdout, nil
}

// CheckDockerAvailable checks if Docker is available
func (e *Executor) CheckDockerAvailable(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "docker", "info")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker is not available or not running: %s", stderr.String())
	}

	return nil
}

// CheckComposeAvailable checks if docker compose is available
func (e *Executor) CheckComposeAvailable(ctx context.Context) error {
	cmdParts := strings.Fields(e.composeCommand)
	cmdParts = append(cmdParts, "version")

	cmd := exec.CommandContext(ctx, cmdParts[0], cmdParts[1:]...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose is not available: %s", stderr.String())
	}

	return nil
}
