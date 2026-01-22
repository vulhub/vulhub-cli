package commands

import (
	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/environment"
	"github.com/vulhub/vulhub-cli/internal/github"
	"github.com/vulhub/vulhub-cli/internal/resolver"
)

// Commands holds all dependencies and provides CLI command methods.
// All command logic is implemented as methods on this struct,
// allowing direct access to dependencies without parameter passing.
type Commands struct {
	Config      config.Manager
	Environment environment.Manager
	Resolver    resolver.Resolver
	Downloader  *github.Downloader
}

// New creates a new Commands instance with all dependencies injected via fx.
func New(
	cfgMgr config.Manager,
	envMgr environment.Manager,
	res resolver.Resolver,
	downloader *github.Downloader,
) *Commands {
	return &Commands{
		Config:      cfgMgr,
		Environment: envMgr,
		Resolver:    res,
		Downloader:  downloader,
	}
}

// All returns all CLI commands.
func (c *Commands) All() []*cli.Command {
	return []*cli.Command{
		c.Init(),
		c.Syncup(),
		c.Start(),
		c.Stop(),
		c.Down(),
		c.Restart(),
		c.List(),
		c.ListAvailable(),
		c.Status(),
		c.Search(),
		c.Info(),
		c.Clean(),
	}
}
