package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/environment"
)

// ListCommand creates the list command
func ListCommand(cfgMgr config.Manager, envMgr environment.Manager) *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List running vulnerability environments",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runList(ctx, cfgMgr, envMgr)
		},
	}
}

func runList(ctx context.Context, cfgMgr config.Manager, envMgr environment.Manager) error {
	table := ui.NewTable()

	// Check if initialized
	if !cfgMgr.IsInitialized() {
		return fmt.Errorf("vulhub-cli is not initialized, please run 'vulhub init' first")
	}

	// Get running environments
	statuses, err := envMgr.ListRunning(ctx)
	if err != nil {
		return err
	}

	table.PrintEnvironmentStatuses(statuses)

	return nil
}

// ListAvailableCommand creates the list-available command
func ListAvailableCommand(cfgMgr config.Manager) *cli.Command {
	return &cli.Command{
		Name:  "list-available",
		Usage: "List all available vulnerability environments",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"n"},
				Usage:   "Limit number of results",
				Value:   0,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runListAvailable(ctx, cfgMgr, int(cmd.Int("limit")))
		},
	}
}

func runListAvailable(ctx context.Context, cfgMgr config.Manager, limit int) error {
	table := ui.NewTable()

	// Check if initialized
	if !cfgMgr.IsInitialized() {
		return fmt.Errorf("vulhub-cli is not initialized, please run 'vulhub init' first")
	}

	// Load environments
	envList, err := cfgMgr.LoadEnvironments(ctx)
	if err != nil {
		return err
	}

	envs := envList.Environment
	if limit > 0 && limit < len(envs) {
		envs = envs[:limit]
		fmt.Printf("Showing first %d environments (of %d total)\n\n", limit, len(envList.Environment))
	}

	table.PrintEnvironments(envs)

	return nil
}
