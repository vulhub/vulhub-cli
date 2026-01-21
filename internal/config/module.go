package config

import (
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
