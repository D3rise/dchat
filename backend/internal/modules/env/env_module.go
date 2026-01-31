package env

import (
	"github.com/D3rise/dchat/internal/interfaces"
	"go.uber.org/fx"
)

var Module = fx.Module("env",
	fx.Provide(
		fx.Annotate(
			NewProcessEnvService,
			fx.As(new(interfaces.EnvService)),
		),
	),
)
