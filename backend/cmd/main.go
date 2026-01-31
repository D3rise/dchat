package main

import (
	"github.com/D3rise/dchat/internal/modules/echo"
	"github.com/D3rise/dchat/internal/modules/env"
	"github.com/D3rise/dchat/internal/modules/rtc"
	"github.com/D3rise/dchat/internal/server"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		env.Module,
		rtc.Module,
		echo.Module,

		fx.Provide(
			server.NewHandler,
		),
		fx.Invoke(func(server.Server) {}),
	).Run()
}
