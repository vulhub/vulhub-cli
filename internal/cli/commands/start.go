package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/internal/environment"
	"github.com/vulhub/vulhub-cli/internal/github"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// Start creates the start command
func (c *Commands) Start() *cli.Command {
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

			return c.runStart(ctx, keyword, startOptions{
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

func (c *Commands) runStart(ctx context.Context, keyword string, opts startOptions) error {
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

	// Resolve keyword
	result, err := c.Resolver.Resolve(ctx, keyword)
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

	if err := c.Environment.Start(ctx, *env, startOpts); err != nil {
		// Check for rate limit error and prompt for token setup
		if github.IsRateLimitError(err) {
			cfg := c.Config.Get()
			if cfg.GitHub.Token == "" {
				if c.PromptTokenSetup(ctx) {
					fmt.Println()
					fmt.Println("Token saved! Please run the command again to continue.")
					return nil
				}
			}
		}
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

	table.PrintSuccess(fmt.Sprintf("Environment '%s' started successfully!", env.Path))

	return nil
}
