package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
)

// Restart creates the restart command
func (c *Commands) Restart() *cli.Command {
	return &cli.Command{
		Name:      "restart",
		Usage:     "Restart a vulnerability environment",
		ArgsUsage: "[keyword]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "yes",
				Aliases: []string{"y"},
				Usage:   "Skip confirmation prompts",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			keyword := cmd.Args().First()
			if keyword == "" {
				return fmt.Errorf("please provide a keyword (CVE number, path, or application name)")
			}

			return c.runRestart(ctx, keyword, cmd.Bool("yes"))
		},
	}
}

func (c *Commands) runRestart(ctx context.Context, keyword string, skipConfirm bool) error {
	table := ui.NewTable()

	// Check if initialized, prompt to initialize if not
	initialized, err := c.ensureInitialized(ctx)
	if err != nil {
		return err
	}
	if !initialized {
		return nil
	}

	// Check if sync is needed
	if _, err := c.checkAndPromptSync(ctx); err != nil {
		return err
	}

	// Resolve keyword within downloaded environments only
	env, err := c.resolveEnvironment(ctx, keyword, ScopeDownloaded, skipConfirm)
	if err != nil {
		return err
	}
	if env == nil {
		return nil // No environments found, message already printed
	}

	// Restart the environment
	table.PrintInfo(fmt.Sprintf("Restarting environment: %s", env.Path))

	if err := c.Environment.Restart(ctx, *env); err != nil {
		return err
	}

	// Get status to show ports
	status, err := c.Environment.Status(ctx, *env)
	if err != nil {
		table.PrintWarning(fmt.Sprintf("Failed to get status: %v", err))
	} else if len(status.Containers) > 0 {
		fmt.Println()
		fmt.Println("Containers:")
		table.PrintContainerStatuses(status.Containers)
	}

	table.PrintSuccess(fmt.Sprintf("Environment '%s' restarted successfully!", env.Path))

	return nil
}
