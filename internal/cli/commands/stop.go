package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
)

// Stop creates the stop command
func (c *Commands) Stop() *cli.Command {
	return &cli.Command{
		Name:      "stop",
		Usage:     "Stop a running vulnerability environment",
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

			return c.runStop(ctx, keyword, cmd.Bool("yes"))
		},
	}
}

func (c *Commands) runStop(ctx context.Context, keyword string, skipConfirm bool) error {
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

	// Resolve keyword within running environments only
	env, err := c.resolveEnvironment(ctx, keyword, ScopeRunning, skipConfirm)
	if err != nil {
		return err
	}
	if env == nil {
		return nil // No environments found, message already printed
	}

	// Stop the environment
	table.PrintInfo(fmt.Sprintf("Stopping environment: %s", env.Path))

	if err := c.Environment.Stop(ctx, *env); err != nil {
		return err
	}

	table.PrintSuccess(fmt.Sprintf("Environment '%s' stopped successfully!", env.Path))

	return nil
}
