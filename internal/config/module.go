package config

import (
	"context"

	"go.uber.org/fx"
)

// Module provides the config module for fx
var Module = fx.Module("config",
	fx.Provide(
		NewPaths,
		NewConfigManager,
		func(m *ConfigManager) Manager {
			return m
		},
	),
	fx.Invoke(func(lc fx.Lifecycle, m *ConfigManager) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				// Load config on startup (ignore errors for uninitialized state)
				_ = m.Load(ctx)
				return nil
			},
		})
	}),
)

// ModuleWithPath provides the config module with a custom config path
func ModuleWithPath(configDir string) fx.Option {
	return fx.Module("config",
		fx.Provide(
			func() *Paths {
				return NewPathsWithDir(configDir)
			},
			NewConfigManager,
			func(m *ConfigManager) Manager {
				return m
			},
		),
	)
}
