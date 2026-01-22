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

// StartCommand creates the start command
func StartCommand(
	cfgMgr config.Manager,
	envMgr environment.Manager,
	res resolver.Resolver,
	downloader *github.Downloader,
) *cli.Command {
	return &cli.Command{
		Name:      "start",
		Usage:     "Start a vulnerability environment",
		ArgsUsage: "[keyword]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "yes",
				Aliases: []string{"y"},
				Usage:   "Skip confirmation prompts",
			},
			&cli.BoolFlag{
				Name:  "pull",
				Usage: "Pull images before starting",
			},
			&cli.BoolFlag{
				Name:  "build",
				Usage: "Build images before starting",
			},
			&cli.BoolFlag{
				Name:  "force-recreate",
				Usage: "Force recreate containers",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			keyword := cmd.Args().First()
			if keyword == "" {
				return fmt.Errorf("please provide a keyword (CVE number, path, or application name)")
			}

			return runStart(ctx, cfgMgr, envMgr, res, downloader, keyword, startOptions{
				yes:           cmd.Bool("yes"),
				pull:          cmd.Bool("pull"),
				build:         cmd.Bool("build"),
				forceRecreate: cmd.Bool("force-recreate"),
			})
		},
	}
}

type startOptions struct {
	yes           bool
	pull          bool
	build         bool
	forceRecreate bool
}

func runStart(
	ctx context.Context,
	cfgMgr config.Manager,
	envMgr environment.Manager,
	res resolver.Resolver,
	downloader *github.Downloader,
	keyword string,
	opts startOptions,
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
		return fmt.Errorf("no environment found matching '%s'. Try 'vulhub search %s' to find environments", keyword, keyword)
	}

	if result.HasMultipleMatches() {
		if opts.yes {
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

	// Start the environment
	table.PrintInfo(fmt.Sprintf("Starting environment: %s", env.Path))

	startOpts := environment.StartOptions{
		Pull:          opts.pull,
		Build:         opts.build,
		ForceRecreate: opts.forceRecreate,
	}

	if err := envMgr.Start(ctx, *env, startOpts); err != nil {
		return err
	}

	// Get status to show ports
	status, err := envMgr.Status(ctx, *env)
	if err != nil {
		table.PrintWarning(fmt.Sprintf("Failed to get status: %v", err))
	} else if len(status.Containers) > 0 {
		fmt.Println()
		fmt.Println("Containers:")
		table.PrintContainerStatuses(status.Containers)
	}

	table.PrintSuccess(fmt.Sprintf("Environment '%s' started successfully!", env.Path))

	return nil
}
