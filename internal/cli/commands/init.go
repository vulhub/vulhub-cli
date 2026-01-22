package commands

import (
	"context"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// Init creates the init command
func (c *Commands) Init() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize vulhub-cli configuration",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force overwrite existing configuration",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return c.runInit(ctx, cmd.Bool("force"))
		},
	}
}

func (c *Commands) runInit(ctx context.Context, force bool) error {
	table := ui.NewTable()
	selector := ui.NewSelector()

	// Check if already initialized
	if c.Config.IsInitialized() && !force {
		confirmed, err := selector.Confirm("Configuration already exists. Overwrite?", false)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("Initialization cancelled.")
			return nil
		}
	}

	// Ensure config directory exists
	if err := c.Config.Paths().EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create default configuration
	table.PrintInfo("Creating default configuration...")
	defaultCfg := types.DefaultConfig()
	c.Config.Set(&defaultCfg)
	if err := c.Config.Save(ctx); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// Download environments.toml
	table.PrintInfo("Downloading environment list from GitHub...")
	envData, err := c.Downloader.DownloadEnvironmentsList(ctx)
	if err != nil {
		return fmt.Errorf("failed to download environments list: %w", err)
	}

	// Parse and save environments
	var envList types.EnvironmentList
	if _, err := toml.Decode(string(envData), &envList); err != nil {
		return fmt.Errorf("failed to parse environments list: %w", err)
	}

	if err := c.Config.SaveEnvironments(ctx, &envList); err != nil {
		return fmt.Errorf("failed to save environments list: %w", err)
	}

	// Update last sync time
	if err := c.Config.UpdateLastSyncTime(ctx); err != nil {
		return fmt.Errorf("failed to update sync time: %w", err)
	}

	// Ensure environments directory exists
	if err := c.Config.Paths().EnsureEnvironmentsDir(); err != nil {
		return fmt.Errorf("failed to create environments directory: %w", err)
	}

	table.PrintSuccess(fmt.Sprintf("Initialization complete! Found %d environments.", len(envList.Environment)))
	fmt.Println()
	fmt.Println("Configuration directory:", c.Config.Paths().ConfigDir())
	fmt.Println()
	fmt.Println("Quick start:")
	fmt.Println("  vulhub search log4j         # Search for environments")
	fmt.Println("  vulhub info CVE-2021-44228  # View environment details")
	fmt.Println("  vulhub start CVE-2021-44228 # Start an environment")

	return nil
}
