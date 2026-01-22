package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/environment"
	"github.com/vulhub/vulhub-cli/internal/github"
	"github.com/vulhub/vulhub-cli/internal/resolver"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// StatusCommand creates the status command
func StatusCommand(
	cfgMgr config.Manager,
	envMgr environment.Manager,
	res resolver.Resolver,
	downloader *github.Downloader,
) *cli.Command {
	return &cli.Command{
		Name:      "status",
		Usage:     "Show status of vulnerability environments",
		ArgsUsage: "[keyword]",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			keyword := cmd.Args().First()
			return runStatus(ctx, cfgMgr, envMgr, res, downloader, keyword)
		},
	}
}

func runStatus(
	ctx context.Context,
	cfgMgr config.Manager,
	envMgr environment.Manager,
	res resolver.Resolver,
	downloader *github.Downloader,
	keyword string,
) error {
	table := ui.NewTable()
	selector := ui.NewSelector()

	// Check if initialized, prompt to initialize if not
	initialized, err := EnsureInitialized(ctx, cfgMgr, downloader)
	if err != nil {
		return err
	}
	if !initialized {
		return nil
	}

	// Check if sync is needed
	if _, err := CheckAndPromptSync(ctx, cfgMgr, downloader); err != nil {
		return err
	}

	// If no keyword, show all running environments
	if keyword == "" {
		statuses, err := envMgr.ListRunning(ctx)
		if err != nil {
			return err
		}

		if len(statuses) == 0 {
			fmt.Println("No running environments.")
			return nil
		}

		for _, status := range statuses {
			printEnvironmentStatus(table, &status)
			fmt.Println()
		}

		return nil
	}

	// Resolve keyword
	result, err := res.Resolve(ctx, keyword)
	if err != nil {
		return err
	}

	var env *types.Environment

	if result.HasNoMatches() {
		return fmt.Errorf("no environment found matching '%s'", keyword)
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
	status, err := envMgr.Status(ctx, *env)
	if err != nil {
		return err
	}

	printEnvironmentStatus(table, status)

	return nil
}

func printEnvironmentStatus(table *ui.Table, status *types.EnvironmentStatus) {
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
