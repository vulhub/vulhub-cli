package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/commands"
	"github.com/vulhub/vulhub-cli/internal/httpclient"
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
	cmds *commands.Commands,
	httpClient *httpclient.Client,
) *cli.Command {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Fprintf(cmd.Root().Writer,
			"Vulhub CLI\n  Version:    %s\n  Commit:     %s\n  Build Time: %s\n",
			Version, Commit, BuildTime,
		)
	}

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
			&cli.StringFlag{
				Name:  "proxy",
				Usage: "Proxy server URL (e.g., http://127.0.0.1:8080 or socks5://127.0.0.1:1080)",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			// Configure log level based on verbose flag
			if cmd.Bool("verbose") {
				handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				})
				slog.SetDefault(slog.New(handler))
				slog.Debug("verbose mode enabled")
			}

			// Configure proxy if specified via CLI flag (highest priority)
			proxyURL := cmd.String("proxy")
			if proxyURL != "" {
				// Validate and set proxy URL
				if err := httpclient.ValidateProxyURL(proxyURL); err != nil {
					return ctx, err
				}

				if err := httpClient.SetProxyURL(proxyURL); err != nil {
					return ctx, fmt.Errorf("failed to configure proxy: %w", err)
				}
			}
			return ctx, nil
		},
		Commands: cmds.All(),
	}
}

// Run runs the CLI application
func Run(ctx context.Context, app *cli.Command, args []string) error {
	return app.Run(ctx, args)
}
