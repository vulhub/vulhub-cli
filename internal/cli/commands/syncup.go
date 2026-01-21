package commands

import (
	"context"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/github"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// SyncupCommand creates the syncup command
func SyncupCommand(cfgMgr config.Manager, downloader *github.Downloader) *cli.Command {
	return &cli.Command{
		Name:  "syncup",
		Usage: "Sync environment list from GitHub",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runSyncup(ctx, cfgMgr, downloader)
		},
	}
}

func runSyncup(ctx context.Context, cfgMgr config.Manager, downloader *github.Downloader) error {
	table := ui.NewTable()

	// Check if initialized
	if !cfgMgr.IsInitialized() {
		return fmt.Errorf("vulhub-cli is not initialized, please run 'vulhub init' first")
	}

	// Load current environments
	currentEnvs, err := cfgMgr.LoadEnvironments(ctx)
	if err != nil {
		currentEnvs = &types.EnvironmentList{}
	}
	currentCount := len(currentEnvs.Environment)

	// Download latest environments.toml
	table.PrintInfo("Downloading latest environment list from GitHub...")
	envData, err := downloader.DownloadEnvironmentsList(ctx)
	if err != nil {
		return fmt.Errorf("failed to download environments list: %w", err)
	}

	// Parse new environments
	var newEnvs types.EnvironmentList
	if _, err := toml.Decode(string(envData), &newEnvs); err != nil {
		return fmt.Errorf("failed to parse environments list: %w", err)
	}

	// Save environments
	if err := cfgMgr.SaveEnvironments(ctx, &newEnvs); err != nil {
		return fmt.Errorf("failed to save environments list: %w", err)
	}

	newCount := len(newEnvs.Environment)

	// Print summary
	table.PrintSuccess("Environment list updated successfully!")
	fmt.Printf("Previous: %d environments\n", currentCount)
	fmt.Printf("Current:  %d environments\n", newCount)

	if newCount > currentCount {
		fmt.Printf("Added:    %d new environments\n", newCount-currentCount)
	} else if newCount < currentCount {
		fmt.Printf("Removed:  %d environments\n", currentCount-newCount)
	}

	return nil
}
