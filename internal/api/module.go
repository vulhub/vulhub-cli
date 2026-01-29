package api

import (
	"go.uber.org/fx"

	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/environment"
	"github.com/vulhub/vulhub-cli/internal/github"
	"github.com/vulhub/vulhub-cli/internal/resolver"
)

// Module provides the API handlers for fx
var Module = fx.Module("api",
	fx.Provide(
		func(
			cfg config.Manager,
			env environment.Manager,
			res resolver.Resolver,
			downloader *github.Downloader,
		) *Handlers {
			return NewHandlers(cfg, env, res, downloader)
		},
	),
)
