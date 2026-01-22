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

// NewApp creates a new CLI application with all dependencies injected via fx
func NewApp(
	cfgMgr config.Manager,
	envMgr environment.Manager,
	res resolver.Resolver,
	downloader *github.Downloader,
) *cli.Command {
	// Create commands instance with all dependencies
	cmds := commands.New(cfgMgr, envMgr, res, downloader)

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
		Commands: cmds.All(),
	}
}

// Run runs the CLI application
func Run(ctx context.Context, app *cli.Command, args []string) error {
	return app.Run(ctx, args)
}
