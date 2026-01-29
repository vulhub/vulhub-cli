package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/samber/lo"
	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/api"
	"github.com/vulhub/vulhub-cli/internal/cli/ui"
)

// Serve creates the serve command
func (c *Commands) Serve() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Start the web API server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"H"},
				Usage:   "Host address to bind to",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Port number to listen on",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return c.runServe(ctx, cmd)
		},
	}
}

func (c *Commands) runServe(ctx context.Context, cmd *cli.Command) error {
	table := ui.NewTable()

	// Check if initialized
	if !c.Config.IsInitialized() {
		return fmt.Errorf("vulhub-cli is not initialized, please run 'vulhub init' first")
	}

	// Check if GitHub token is configured
	cfg := c.Config.Get()
	if cfg.GitHub.Token == "" {
		table.PrintWarning("GitHub authentication is required for the web API server.")
		fmt.Println("The web API needs to fetch README content from GitHub.")
		fmt.Println()

		// Prompt user to authenticate
		if !c.PromptTokenSetup(ctx) {
			return fmt.Errorf("GitHub authentication is required to start the web server")
		}

		// Refresh config after authentication
		cfg = c.Config.Get()
	}
	host := lo.Ternary(cfg.Web.Host != "", cfg.Web.Host, "127.0.0.1")
	port := lo.Ternary(cfg.Web.Port > 0, cfg.Web.Port, 8080)

	// CLI flags take precedence over config
	if cmd.IsSet("host") {
		host = cmd.String("host")
	}
	if cmd.IsSet("port") {
		port = cmd.Int("port")
	}

	// Create server using already-injected dependencies from Commands
	server := api.NewServer(
		host,
		port,
		c.WebHandlers,
	)

	// Print startup message
	table.PrintSuccess(fmt.Sprintf("Starting vulhub API server on http://%s:%d", host, port))
	fmt.Println()
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /health                              - Health check")
	fmt.Println("  GET  /api/v1/status                       - System status")
	fmt.Println("  POST /api/v1/syncup                       - Sync environment list")
	fmt.Println("  GET  /api/v1/environments                 - List all environments")
	fmt.Println("  GET  /api/v1/environments/downloaded      - List downloaded environments")
	fmt.Println("  GET  /api/v1/environments/running         - List running environments")
	fmt.Println("  GET  /api/v1/environments/info/*path      - Get environment info")
	fmt.Println("  GET  /api/v1/environments/status/*path    - Get environment status")
	fmt.Println("  POST /api/v1/environments/start/*path     - Start environment")
	fmt.Println("  POST /api/v1/environments/stop/*path      - Stop environment")
	fmt.Println("  POST /api/v1/environments/restart/*path   - Restart environment")
	fmt.Println("  DELETE /api/v1/environments/clean/*path   - Clean environment")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop the server")
	fmt.Println()

	// Start the server
	if err := server.Start(ctx); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	table.PrintInfo("Received shutdown signal, stopping server...")

	// Stop the server gracefully
	if err := server.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	return nil
}
