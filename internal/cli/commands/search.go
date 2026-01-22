package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
)

// Search creates the search command
func (c *Commands) Search() *cli.Command {
	return &cli.Command{
		Name:      "search",
		Usage:     "Search for vulnerability environments",
		ArgsUsage: "[keyword]",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			keyword := cmd.Args().First()
			return c.runSearch(ctx, keyword)
		},
	}
}

func (c *Commands) runSearch(ctx context.Context, keyword string) error {
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

	// Load all environments
	envList, err := c.Config.LoadEnvironments(ctx)
	if err != nil {
		return err
	}

	envs := envList.Environment

	browser := ui.NewEnvironmentBrowser()
	table := ui.NewTable()
	pager := ui.NewPager()

	// Loop: browse -> info -> back to browse
	for {
		result, err := browser.BrowseWithOptions(envs, ui.BrowseOptions{
			Title:         "Search Environments",
			InitialSearch: keyword,
		})
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
		// Clear the initial search so user can search again
		keyword = ""
	}
}
