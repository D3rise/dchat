package main

import (
	"github.com/D3rise/dchat/internal/domain/database"
	"github.com/D3rise/dchat/internal/domain/server"
	"github.com/D3rise/dchat/internal/modules/echo"
	"github.com/D3rise/dchat/internal/modules/env"
	"github.com/D3rise/dchat/internal/modules/rtc"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		env.Module,
		rtc.Module,
		echo.Module,

		// Logger
		fx.Provide(zap.NewExample),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),

		// Database
		fx.Provide(database.NewDatabaseConn),
		fx.Invoke(func(database.Database) {}),

		// Server
		fx.Provide(server.NewServer),
		fx.Decorate(server.DecorateServerWithLogger),
		fx.Invoke(func(server.Server) {}),
	).Run()
}
