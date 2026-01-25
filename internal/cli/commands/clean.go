package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
)

// Clean creates the clean command
func (c *Commands) Clean() *cli.Command {
	return &cli.Command{
		Name:      "clean",
		Usage:     "Completely remove an environment (containers, volumes, and local files)",
		ArgsUsage: "[keyword]",
		Aliases:   []string{"rm", "down"},
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

			return c.runClean(ctx, keyword, cmd.Bool("yes"))
		},
	}
}

func (c *Commands) runClean(ctx context.Context, keyword string, skipConfirm bool) error {
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

	// Clean the environment
	table.PrintInfo(fmt.Sprintf("Stopping and removing environment: %s", env.Path))

	if err := c.Environment.Down(ctx, *env); err != nil {
		return err
	}

	table.PrintSuccess(fmt.Sprintf("Environment '%s' has been stopped and removed!", env.Path))

	return nil
}
