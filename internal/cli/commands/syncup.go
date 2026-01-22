package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/github"
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

	// Use shared sync logic
	result, err := PerformSync(ctx, cfgMgr, downloader, table)
	if err != nil {
		return err
	}

	// Print detailed summary
	table.PrintSuccess("Environment list updated successfully!")
	fmt.Printf("Previous: %d environments\n", result.PreviousCount)
	fmt.Printf("Current:  %d environments\n", result.CurrentCount)

	if result.CurrentCount > result.PreviousCount {
		fmt.Printf("Added:    %d new environments\n", result.CurrentCount-result.PreviousCount)
	} else if result.CurrentCount < result.PreviousCount {
		fmt.Printf("Removed:  %d environments\n", result.PreviousCount-result.CurrentCount)
	}

	return nil
}
