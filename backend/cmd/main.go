package main

import (
	"github.com/D3rise/dchat/internal/modules/echo"
	"github.com/D3rise/dchat/internal/modules/env"
	"github.com/D3rise/dchat/internal/modules/rtc"
	"github.com/D3rise/dchat/internal/server"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		env.Module,
		rtc.Module,
		echo.Module,

		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			server.NewHandler,
			zap.NewExample,
		),
		fx.Decorate(
			server.DecorateServerWithLogger,
		),
		fx.Invoke(func(server.Server) {}),
	).Run()
}
