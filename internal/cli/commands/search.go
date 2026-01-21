package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/resolver"
)

// SearchCommand creates the search command
func SearchCommand(cfgMgr config.Manager, res resolver.Resolver) *cli.Command {
	return &cli.Command{
		Name:      "search",
		Usage:     "Search for vulnerability environments",
		ArgsUsage: "[keyword]",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"n"},
				Usage:   "Limit number of results",
				Value:   20,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			keyword := cmd.Args().First()
			if keyword == "" {
				return fmt.Errorf("please provide a search keyword")
			}

			return runSearch(ctx, cfgMgr, res, keyword, int(cmd.Int("limit")))
		},
	}
}

func runSearch(
	ctx context.Context,
	cfgMgr config.Manager,
	res resolver.Resolver,
	keyword string,
	limit int,
) error {
	table := ui.NewTable()

	// Check if initialized
	if !cfgMgr.IsInitialized() {
		return fmt.Errorf("vulhub-cli is not initialized, please run 'vulhub init' first")
	}

	// Resolve keyword
	result, err := res.Resolve(ctx, keyword)
	if err != nil {
		return err
	}

	if result.HasNoMatches() {
		fmt.Printf("No environments found matching '%s'.\n", keyword)
		fmt.Println()
		fmt.Println("Search tips:")
		fmt.Println("  - Try a CVE number: vulhub search CVE-2021-44228")
		fmt.Println("  - Try an application name: vulhub search log4j")
		fmt.Println("  - Try a partial keyword: vulhub search apache")
		return nil
	}

	envs := result.GetMatchedEnvironments()

	// Apply limit
	total := len(envs)
	if limit > 0 && limit < len(envs) {
		envs = envs[:limit]
	}

	fmt.Printf("Found %d environment(s) matching '%s'", total, keyword)
	if total > len(envs) {
		fmt.Printf(" (showing first %d)", len(envs))
	}
	fmt.Println()
	fmt.Println()

	table.PrintEnvironments(envs)

	if total > len(envs) {
		fmt.Printf("\nUse --limit to show more results: vulhub search %s --limit %d\n", keyword, total)
	}

	return nil
}
