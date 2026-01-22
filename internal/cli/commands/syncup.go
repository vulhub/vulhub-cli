package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
)

// Syncup creates the syncup command
func (c *Commands) Syncup() *cli.Command {
	return &cli.Command{
		Name:  "syncup",
		Usage: "Sync environment list from GitHub",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return c.runSyncup(ctx)
		},
	}
}

func (c *Commands) runSyncup(ctx context.Context) error {
	table := ui.NewTable()

	// Check if initialized
	if !c.Config.IsInitialized() {
		return fmt.Errorf("vulhub-cli is not initialized, please run 'vulhub init' first")
	}

	// Use shared sync logic
	result, err := c.performSync(ctx, table)
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
