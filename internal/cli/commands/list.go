package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
)

// List creates the list command
func (c *Commands) List() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List all downloaded vulnerability environments",
		Aliases: []string{"ls"},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return c.runList(ctx)
		},
	}
}

func (c *Commands) runList(ctx context.Context) error {
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

	// Get all downloaded environments
	statuses, err := c.Environment.ListDownloaded(ctx)
	if err != nil {
		return err
	}

	table.PrintEnvironmentStatuses(statuses)

	return nil
}

// ListAvailable creates the list-available command
func (c *Commands) ListAvailable() *cli.Command {
	return &cli.Command{
		Name:    "list-available",
		Usage:   "List all available vulnerability environments",
		Aliases: []string{"la"},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"n"},
				Usage:   "Limit number of results",
				Value:   0,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return c.runListAvailable(ctx, int(cmd.Int("limit")))
		},
	}
}

func (c *Commands) runListAvailable(ctx context.Context, limit int) error {
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

	// Load environments
	envList, err := c.Config.LoadEnvironments(ctx)
	if err != nil {
		return err
	}

	envs := envList.Environment
	if limit > 0 && limit < len(envs) {
		envs = envs[:limit]
	}

	browser := ui.NewEnvironmentBrowser()
	table := ui.NewTable()
	pager := ui.NewPager()

	// Loop: browse -> info -> back to browse
	for {
		result, err := browser.Browse(envs)
		if err != nil {
			return err
		}

		// User quit the browser
		if result.Quit || result.Selected == nil {
			return nil
		}

		// User selected an environment, show its info
		info, err := c.Environment.GetInfo(ctx, *result.Selected)
		if err != nil {
			return err
		}

		// Display info using pager (will return when user presses q/esc)
		content := table.FormatEnvironmentInfo(info)
		if err := pager.DisplayWithContent(fmt.Sprintf("Environment: %s", result.Selected.Path), content); err != nil {
			return err
		}

		// After pager exits, loop back to browser
	}
}
