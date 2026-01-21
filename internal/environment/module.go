package environment

import (
	"go.uber.org/fx"
)

// Module provides the environment module for fx
var Module = fx.Module("environment",
	fx.Provide(
		NewEnvironmentManager,
		func(m *EnvironmentManager) Manager {
			return m
		},
	),
)
