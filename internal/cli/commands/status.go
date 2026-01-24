package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// Status creates the status command
func (c *Commands) Status() *cli.Command {
	return &cli.Command{
		Name:      "status",
		Aliases:   []string{"ls", "list"},
		Usage:     "Show status of vulnerability environments",
		ArgsUsage: "[keyword]",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			keyword := cmd.Args().First()
			return c.runStatus(ctx, keyword)
		},
	}
}

func (c *Commands) runStatus(ctx context.Context, keyword string) error {
	table := ui.NewTable()
	selector := ui.NewSelector()

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

	// If no keyword, show all downloaded environments
	if keyword == "" {
		statuses, err := c.Environment.ListDownloaded(ctx)
		if err != nil {
			return err
		}

		table.PrintEnvironmentStatuses(statuses)
		return nil
	}

	// Resolve keyword
	result, err := c.Resolver.Resolve(ctx, keyword)
	if err != nil {
		return err
	}

	var env *types.Environment

	if result.HasNoMatches() {
		return errNoEnvironmentFound(keyword)
	}

	if result.HasMultipleMatches() {
		envs := result.GetMatchedEnvironments()
		env, err = selector.SelectEnvironment(envs, fmt.Sprintf("Multiple environments match '%s'. Please select one:", keyword))
		if err != nil {
			return err
		}
	} else {
		env = result.Environment
	}

	// Get status
	status, err := c.Environment.Status(ctx, *env)
	if err != nil {
		return err
	}

	c.printEnvironmentStatus(table, status)

	return nil
}

func (c *Commands) printEnvironmentStatus(table *ui.Table, status *types.EnvironmentStatus) {
	env := status.Environment

	fmt.Printf("%s %s\n", ui.MutedStyle.Render("Environment:"), ui.PathStyle.Render(env.Path))
	if len(env.CVE) > 0 {
		fmt.Printf("%s %s\n", ui.MutedStyle.Render("CVE:"), strings.Join(env.CVE, ", "))
	}
	fmt.Printf("%s %s\n", ui.MutedStyle.Render("Application:"), env.App)
	if status.LocalPath != "" {
		fmt.Printf("%s %s\n", ui.MutedStyle.Render("Downloaded:"), ui.SuccessStyle.Render("true"))
	} else {
		fmt.Printf("%s %s\n", ui.MutedStyle.Render("Downloaded:"), ui.MutedStyle.Render("false"))
	}

	if status.Running {
		fmt.Printf("%s %s\n", ui.MutedStyle.Render("Status:"), ui.StatusRunningStyle.Render("running"))
	} else {
		fmt.Printf("%s %s\n", ui.MutedStyle.Render("Status:"), ui.StatusStoppedStyle.Render("stopped"))
	}

	if len(status.Containers) > 0 {
		fmt.Println()
		fmt.Println("Containers:")
		table.PrintContainerStatuses(status.Containers)
	}
}
