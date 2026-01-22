package commands

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/github"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// CheckDockerEnvironment checks if Docker and Docker Compose are available
// Returns nil if everything is ready, otherwise returns an error with installation instructions
func CheckDockerEnvironment(ctx context.Context) error {
	table := ui.NewTable()

	// Check if docker command exists
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		table.PrintError("Docker is not installed.")
		printDockerInstallGuide()
		return fmt.Errorf("docker is not installed")
	}

	// Check if docker daemon is running
	cmd := exec.CommandContext(ctx, dockerPath, "info")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		table.PrintError("Docker is installed but not running.")
		fmt.Println()
		fmt.Println("Please start Docker:")
		switch runtime.GOOS {
		case "windows":
			fmt.Println("  - Open Docker Desktop from the Start menu")
		case "darwin":
			fmt.Println("  - Open Docker Desktop from Applications")
		default:
			fmt.Println("  - Run: sudo systemctl start docker")
		}
		fmt.Println()
		return fmt.Errorf("docker daemon is not running")
	}

	// Check if docker compose is available
	cmd = exec.CommandContext(ctx, dockerPath, "compose", "version")
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		table.PrintError("Docker Compose is not available.")
		printDockerComposeInstallGuide()
		return fmt.Errorf("docker compose is not available")
	}

	return nil
}

// printDockerInstallGuide prints installation instructions for Docker
func printDockerInstallGuide() {
	fmt.Println()
	fmt.Println("Please install Docker:")
	fmt.Println()

	switch runtime.GOOS {
	case "windows":
		fmt.Println("  Windows:")
		fmt.Println("    1. Download Docker Desktop from: https://docs.docker.com/desktop/install/windows-install/")
		fmt.Println("    2. Run the installer and follow the instructions")
		fmt.Println("    3. Restart your computer if prompted")
	case "darwin":
		fmt.Println("  macOS:")
		fmt.Println("    Option 1 - Docker Desktop:")
		fmt.Println("      Download from: https://docs.docker.com/desktop/install/mac-install/")
		fmt.Println()
		fmt.Println("    Option 2 - Homebrew:")
		fmt.Println("      brew install --cask docker")
	default:
		fmt.Println("  Linux:")
		fmt.Println("    Ubuntu/Debian:")
		fmt.Println("      curl -fsSL https://get.docker.com | sh")
		fmt.Println("      sudo usermod -aG docker $USER")
		fmt.Println()
		fmt.Println("    Or follow the official guide:")
		fmt.Println("      https://docs.docker.com/engine/install/")
	}
	fmt.Println()
}

// printDockerComposeInstallGuide prints installation instructions for Docker Compose
func printDockerComposeInstallGuide() {
	fmt.Println()
	fmt.Println("Docker Compose is required but not available.")
	fmt.Println()

	switch runtime.GOOS {
	case "windows", "darwin":
		fmt.Println("  Docker Compose should be included with Docker Desktop.")
		fmt.Println("  Please update Docker Desktop to the latest version:")
		fmt.Println("    https://docs.docker.com/desktop/")
	default:
		fmt.Println("  Linux:")
		fmt.Println("    Docker Compose V2 (recommended):")
		fmt.Println("      sudo apt-get update")
		fmt.Println("      sudo apt-get install docker-compose-plugin")
		fmt.Println()
		fmt.Println("    Or install manually:")
		fmt.Println("      https://docs.docker.com/compose/install/linux/")
	}
	fmt.Println()
}

// EnsureInitialized checks if vulhub-cli is initialized, and prompts the user to initialize if not
// Returns true if initialized (or just initialized), false if user declined
func EnsureInitialized(ctx context.Context, cfgMgr config.Manager, downloader *github.Downloader) (bool, error) {
	// First, check Docker environment
	if err := CheckDockerEnvironment(ctx); err != nil {
		return false, err
	}

	if cfgMgr.IsInitialized() {
		return true, nil
	}

	selector := ui.NewSelector()
	table := ui.NewTable()

	// Prompt user to initialize
	fmt.Println("vulhub-cli is not initialized.")
	confirmed, err := selector.Confirm("Would you like to initialize now?", true)
	if err != nil {
		return false, err
	}

	if !confirmed {
		fmt.Println("Please run 'vulhub init' to initialize.")
		return false, nil
	}

	// Run initialization
	if err := cfgMgr.Paths().EnsureConfigDir(); err != nil {
		return false, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create default configuration
	table.PrintInfo("Creating default configuration...")
	defaultCfg := types.DefaultConfig()
	cfgMgr.Set(&defaultCfg)
	if err := cfgMgr.Save(ctx); err != nil {
		return false, fmt.Errorf("failed to save configuration: %w", err)
	}

	// Download environments.toml
	table.PrintInfo("Downloading environment list from GitHub...")
	envData, err := downloader.DownloadEnvironmentsList(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to download environments list: %w", err)
	}

	// Parse and save environments
	var envList types.EnvironmentList
	if _, err := toml.Decode(string(envData), &envList); err != nil {
		return false, fmt.Errorf("failed to parse environments list: %w", err)
	}

	if err := cfgMgr.SaveEnvironments(ctx, &envList); err != nil {
		return false, fmt.Errorf("failed to save environments list: %w", err)
	}

	// Update last sync time
	if err := cfgMgr.UpdateLastSyncTime(ctx); err != nil {
		return false, fmt.Errorf("failed to update sync time: %w", err)
	}

	// Ensure environments directory exists
	if err := cfgMgr.Paths().EnsureEnvironmentsDir(); err != nil {
		return false, fmt.Errorf("failed to create environments directory: %w", err)
	}

	table.PrintSuccess(fmt.Sprintf("Initialization complete! Found %d environments.", len(envList.Environment)))
	fmt.Println()

	return true, nil
}

// CheckAndPromptSync checks if environments need to be synced and prompts the user
// Returns true if sync was performed or not needed, false if user declined
func CheckAndPromptSync(ctx context.Context, cfgMgr config.Manager, downloader *github.Downloader) (bool, error) {
	if !cfgMgr.NeedSync() {
		return true, nil
	}

	selector := ui.NewSelector()
	table := ui.NewTable()

	lastSync := cfgMgr.GetLastSyncTime()
	var syncMsg string
	if lastSync.IsZero() {
		syncMsg = "Environment list has never been synced."
	} else {
		days := int(time.Since(lastSync).Hours() / 24)
		syncMsg = fmt.Sprintf("Environment list was last synced %d days ago.", days)
	}

	fmt.Println(syncMsg)
	confirmed, err := selector.Confirm("Would you like to sync now?", true)
	if err != nil {
		return false, err
	}

	if !confirmed {
		return true, nil // Continue without sync
	}

	// Run sync
	result, err := PerformSync(ctx, cfgMgr, downloader, table)
	if err != nil {
		return false, err
	}

	// Print summary
	table.PrintSuccess("Environment list updated!")
	if result.CurrentCount != result.PreviousCount {
		fmt.Printf("Environments: %d -> %d\n", result.PreviousCount, result.CurrentCount)
	}
	fmt.Println()

	return true, nil
}

// SyncResult holds the result of a sync operation
type SyncResult struct {
	PreviousCount int
	CurrentCount  int
}

// PerformSync performs the actual sync operation and returns the result
// This is the core sync logic shared by both CheckAndPromptSync and SyncupCommand
func PerformSync(ctx context.Context, cfgMgr config.Manager, downloader *github.Downloader, table *ui.Table) (*SyncResult, error) {
	// Load current environments
	currentEnvs, err := cfgMgr.LoadEnvironments(ctx)
	if err != nil {
		currentEnvs = &types.EnvironmentList{}
	}

	// Download latest environments.toml
	table.PrintInfo("Downloading latest environment list from GitHub...")
	envData, err := downloader.DownloadEnvironmentsList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to download environments list: %w", err)
	}

	// Parse new environments
	var newEnvs types.EnvironmentList
	if _, err := toml.Decode(string(envData), &newEnvs); err != nil {
		return nil, fmt.Errorf("failed to parse environments list: %w", err)
	}

	// Save environments
	if err := cfgMgr.SaveEnvironments(ctx, &newEnvs); err != nil {
		return nil, fmt.Errorf("failed to save environments list: %w", err)
	}

	// Update last sync time
	if err := cfgMgr.UpdateLastSyncTime(ctx); err != nil {
		return nil, fmt.Errorf("failed to update sync time: %w", err)
	}

	return &SyncResult{
		PreviousCount: len(currentEnvs.Environment),
		CurrentCount:  len(newEnvs.Environment),
	}, nil
}
