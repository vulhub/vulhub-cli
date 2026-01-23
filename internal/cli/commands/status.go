package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// Status creates the status command
func (c *Commands) Status() *cli.Command {
	return &cli.Command{
		Name:      "status",
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

	// If no keyword, show all running environments
	if keyword == "" {
		statuses, err := c.Environment.ListRunning(ctx)
		if err != nil {
			return err
		}

		if len(statuses) == 0 {
			fmt.Println("No running environments.")
			return nil
		}

		for _, status := range statuses {
			c.printEnvironmentStatus(table, &status)
			fmt.Println()
		}

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

	fmt.Printf("Environment: %s\n", env.Path)
	if len(env.CVE) > 0 {
		fmt.Printf("CVE:         %s\n", env.CVE[0])
	}
	fmt.Printf("Application: %s\n", env.App)
	fmt.Printf("Downloaded:  %v\n", status.LocalPath != "")

	if status.Running {
		fmt.Println("Status:      running")
	} else {
		fmt.Println("Status:      stopped")
	}

	if len(status.Containers) > 0 {
		fmt.Println()
		fmt.Println("Containers:")
		table.PrintContainerStatuses(status.Containers)
	}
}
