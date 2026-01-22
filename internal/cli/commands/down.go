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

// DownCommand creates the down command
func DownCommand(
	cfgMgr config.Manager,
	envMgr environment.Manager,
	res resolver.Resolver,
	downloader *github.Downloader,
) *cli.Command {
	return &cli.Command{
		Name:      "down",
		Usage:     "Completely remove an environment (containers, volumes, and local files)",
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

			return runDown(ctx, cfgMgr, envMgr, res, downloader, keyword, cmd.Bool("yes"))
		},
	}
}

func runDown(
	ctx context.Context,
	cfgMgr config.Manager,
	envMgr environment.Manager,
	res resolver.Resolver,
	downloader *github.Downloader,
	keyword string,
	skipConfirm bool,
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
		if skipConfirm {
			return fmt.Errorf("multiple environments found matching '%s'. Please provide a more specific keyword", keyword)
		}

		envs := result.GetMatchedEnvironments()
		env, err = selector.SelectEnvironment(envs, fmt.Sprintf("Multiple environments match '%s'. Please select one:", keyword))
		if err != nil {
			return err
		}
	} else {
		env = result.Environment
	}

	// Down the environment
	table.PrintInfo(fmt.Sprintf("Stopping and removing environment: %s", env.Path))

	if err := envMgr.Down(ctx, *env); err != nil {
		return err
	}

	table.PrintSuccess(fmt.Sprintf("Environment '%s' has been stopped and removed!", env.Path))

	return nil
}
