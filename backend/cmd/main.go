package main

import (
	"net/http"

	"github.com/D3rise/dchat/internal/handlers"
	"github.com/D3rise/dchat/internal/server"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			fx.Annotate(
				server.NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			server.NewHTTPServer,
			zap.NewExample,
		),
		fx.Provide(
			server.AsRoute(handlers.NewEchoHandler),
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
