package main

import (
	"github.com/D3rise/dchat/internal/infrastructure/database"
	server2 "github.com/D3rise/dchat/internal/infrastructure/server"
	"github.com/D3rise/dchat/internal/modules/echo"
	"github.com/D3rise/dchat/internal/modules/env"
	"github.com/D3rise/dchat/internal/modules/rtc"
	"github.com/D3rise/dchat/internal/modules/user"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		env.Module,
		rtc.Module,
		echo.Module,
		user.Module,

		// Logger
		fx.Provide(zap.NewExample),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),

		// Database
		fx.Provide(database.NewDatabaseConn),
		fx.Invoke(func(database.Database) {}),

		// Server
		fx.Provide(server2.NewServer),
		fx.Decorate(server2.DecorateServerWithLogger),
		fx.Invoke(func(server2.Server) {}),
	).Run()
}
