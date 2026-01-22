package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/internal/environment"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// Clean creates the clean command
func (c *Commands) Clean() *cli.Command {
	return &cli.Command{
		Name:      "clean",
		Usage:     "Clean up a vulnerability environment",
		ArgsUsage: "[keyword]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "yes",
				Aliases: []string{"y"},
				Usage:   "Skip confirmation prompts",
			},
			&cli.BoolFlag{
				Name:  "volumes",
				Usage: "Remove volumes",
			},
			&cli.BoolFlag{
				Name:  "images",
				Usage: "Remove images",
			},
			&cli.BoolFlag{
				Name:  "files",
				Usage: "Remove local files",
			},
			&cli.BoolFlag{
				Name:  "all",
				Usage: "Remove everything (volumes, images, and files)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			keyword := cmd.Args().First()
			if keyword == "" {
				return fmt.Errorf("please provide a keyword (CVE number, path, or application name)")
			}

			all := cmd.Bool("all")
			return c.runClean(ctx, keyword, cleanOptions{
				yes:     cmd.Bool("yes"),
				volumes: cmd.Bool("volumes") || all,
				images:  cmd.Bool("images") || all,
				files:   cmd.Bool("files") || all,
			})
		},
	}
}

type cleanOptions struct {
	yes     bool
	volumes bool
	images  bool
	files   bool
}

func (c *Commands) runClean(ctx context.Context, keyword string, opts cleanOptions) error {
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
		return fmt.Errorf("no environment found matching '%s'", keyword)
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

	// Confirm if not skipped
	if !opts.yes {
		msg := fmt.Sprintf("Clean environment '%s'", env.Path)
		if opts.volumes {
			msg += " (including volumes)"
		}
		if opts.images {
			msg += " (including images)"
		}
		if opts.files {
			msg += " (including local files)"
		}
		msg += "?"

		confirmed, err := selector.Confirm(msg, false)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("Clean cancelled.")
			return nil
		}
	}

	// Clean the environment
	table.PrintInfo(fmt.Sprintf("Cleaning environment: %s", env.Path))

	cleanOpts := environment.CleanOptions{
		RemoveVolumes: opts.volumes,
		RemoveImages:  opts.images,
		RemoveFiles:   opts.files,
	}

	if err := c.Environment.Clean(ctx, *env, cleanOpts); err != nil {
		return err
	}

	table.PrintSuccess(fmt.Sprintf("Environment '%s' cleaned successfully!", env.Path))

	return nil
}
