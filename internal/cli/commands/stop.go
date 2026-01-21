package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/environment"
	"github.com/vulhub/vulhub-cli/internal/resolver"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// StopCommand creates the stop command
func StopCommand(
	cfgMgr config.Manager,
	envMgr environment.Manager,
	res resolver.Resolver,
) *cli.Command {
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

			return runStop(ctx, cfgMgr, envMgr, res, keyword, cmd.Bool("yes"))
		},
	}
}

func runStop(
	ctx context.Context,
	cfgMgr config.Manager,
	envMgr environment.Manager,
	res resolver.Resolver,
	keyword string,
	skipConfirm bool,
) error {
	table := ui.NewTable()
	selector := ui.NewSelector()

	// Check if initialized
	if !cfgMgr.IsInitialized() {
		return fmt.Errorf("vulhub-cli is not initialized, please run 'vulhub init' first")
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

	// Stop the environment
	table.PrintInfo(fmt.Sprintf("Stopping environment: %s", env.Path))

	if err := envMgr.Stop(ctx, *env); err != nil {
		return err
	}

	table.PrintSuccess(fmt.Sprintf("Environment '%s' stopped successfully!", env.Path))

	return nil
}
