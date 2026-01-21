package compose

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/vulhub/vulhub-cli/internal/config"
)

// Module provides the compose module for fx
var Module = fx.Module("compose",
	fx.Provide(
		NewComposeClientFromConfig,
		func(c *ComposeClient) Client {
			return c
		},
	),
)

// NewComposeClientFromConfig creates a ComposeClient from the config manager
func NewComposeClientFromConfig(cfgMgr config.Manager, logger *slog.Logger) *ComposeClient {
	cfg := cfgMgr.Get()
	return NewComposeClient(cfg.Docker.ComposeCommand, logger)
}
