package cli

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/commands"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/environment"
	"github.com/vulhub/vulhub-cli/internal/github"
	"github.com/vulhub/vulhub-cli/internal/resolver"
)

// Command is an alias for cli.Command
type Command = cli.Command

// Version information (set by build flags)
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

// AppParams contains the dependencies for the CLI application
type AppParams struct {
	ConfigManager      config.Manager
	EnvironmentManager environment.Manager
	Resolver           resolver.Resolver
	Downloader         *github.Downloader
}

// NewApp creates a new CLI application
func NewApp(params AppParams) *cli.Command {
	return &cli.Command{
		Name:    "vulhub",
		Usage:   "A CLI tool for managing Vulhub vulnerability environments",
		Version: Version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose output",
			},
			&cli.StringFlag{
				Name:  "config",
				Usage: "Path to configuration file",
			},
		},
		Commands: []*cli.Command{
			commands.InitCommand(params.ConfigManager, params.Downloader),
			commands.SyncupCommand(params.ConfigManager, params.Downloader),
			commands.StartCommand(params.ConfigManager, params.EnvironmentManager, params.Resolver),
			commands.StopCommand(params.ConfigManager, params.EnvironmentManager, params.Resolver),
			commands.RestartCommand(params.ConfigManager, params.EnvironmentManager, params.Resolver),
			commands.ListCommand(params.ConfigManager, params.EnvironmentManager),
			commands.ListAvailableCommand(params.ConfigManager),
			commands.StatusCommand(params.ConfigManager, params.EnvironmentManager, params.Resolver),
			commands.SearchCommand(params.ConfigManager, params.Resolver),
			commands.InfoCommand(params.ConfigManager, params.EnvironmentManager, params.Resolver),
			commands.CleanCommand(params.ConfigManager, params.EnvironmentManager, params.Resolver),
		},
	}
}

// Run runs the CLI application
func Run(ctx context.Context, app *cli.Command, args []string) error {
	return app.Run(ctx, args)
}
