package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/fx"

	"github.com/vulhub/vulhub-cli/internal/cli"
	"github.com/vulhub/vulhub-cli/internal/compose"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/environment"
	"github.com/vulhub/vulhub-cli/internal/github"
	"github.com/vulhub/vulhub-cli/internal/resolver"
)

// Version information (set by build flags)
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	// Set version info
	cli.Version = Version
	cli.Commit = Commit
	cli.BuildTime = BuildTime

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	// Run the application
	if err := run(ctx); err != nil {
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	var cliApp *cli.Command

	app := fx.New(
		fx.NopLogger,

		// Provide logger
		fx.Provide(func() *slog.Logger {
			return slog.Default()
		}),

		// Load modules
		config.Module,
		github.Module,
		compose.Module,
		resolver.Module,
		environment.Module,

		// Provide CLI app (fx automatically injects all dependencies)
		fx.Provide(cli.NewApp),

		// Populate CLI app
		fx.Populate(&cliApp),
	)

	// Start fx (initializes dependencies and loads config)
	if err := app.Start(ctx); err != nil {
		return err
	}
	defer app.Stop(ctx)

	// Run CLI
	return cli.Run(ctx, cliApp, os.Args)
}
