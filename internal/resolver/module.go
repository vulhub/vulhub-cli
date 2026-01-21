package resolver

import (
	"go.uber.org/fx"

	"github.com/vulhub/vulhub-cli/internal/config"
)

// Module provides the resolver module for fx
var Module = fx.Module("resolver",
	fx.Provide(
		NewEnvironmentResolver,
		func(r *EnvironmentResolver) Resolver {
			return r
		},
	),
)

// NewResolverWithConfig creates a resolver with a config manager
func NewResolverWithConfig(cfgMgr config.Manager) Resolver {
	return NewEnvironmentResolver(cfgMgr)
}
