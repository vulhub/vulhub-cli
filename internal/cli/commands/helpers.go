package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/github"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// EnsureInitialized checks if vulhub-cli is initialized, and prompts the user to initialize if not
// Returns true if initialized (or just initialized), false if user declined
func EnsureInitialized(ctx context.Context, cfgMgr config.Manager, downloader *github.Downloader) (bool, error) {
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
