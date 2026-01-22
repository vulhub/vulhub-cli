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

// InfoCommand creates the info command
func InfoCommand(
	cfgMgr config.Manager,
	envMgr environment.Manager,
	res resolver.Resolver,
	downloader *github.Downloader,
) *cli.Command {
	return &cli.Command{
		Name:      "info",
		Usage:     "Show detailed information about a vulnerability environment",
		ArgsUsage: "[keyword]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "yes",
				Aliases: []string{"y"},
				Usage:   "Skip confirmation prompts",
			},
			&cli.BoolFlag{
				Name:  "no-readme",
				Usage: "Do not show README content",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			keyword := cmd.Args().First()
			if keyword == "" {
				return fmt.Errorf("please provide a keyword (CVE number, path, or application name)")
			}

			return runInfo(ctx, cfgMgr, envMgr, res, downloader, keyword, infoOptions{
				yes:      cmd.Bool("yes"),
				noReadme: cmd.Bool("no-readme"),
			})
		},
	}
}

type infoOptions struct {
	yes      bool
	noReadme bool
}

func runInfo(
	ctx context.Context,
	cfgMgr config.Manager,
	envMgr environment.Manager,
	res resolver.Resolver,
	downloader *github.Downloader,
	keyword string,
	opts infoOptions,
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

	// Get environment info
	info, err := envMgr.GetInfo(ctx, *env)
	if err != nil {
		return err
	}

	// Clear readme if not wanted
	if opts.noReadme {
		info.Readme = ""
	}

	// Use pager for display
	content := table.FormatEnvironmentInfo(info)
	pager := ui.NewPager()
	return pager.DisplayWithContent(fmt.Sprintf("Environment: %s", env.Path), content)
}
