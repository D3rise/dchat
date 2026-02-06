package echo

import (
	"github.com/D3rise/dchat/internal/domain/server"
	"github.com/D3rise/dchat/internal/modules/echo/services"
	"go.uber.org/fx"
)

// This module is for testing purposes exclusively
var Module = fx.Module("echo",
	fx.Provide(
		services.NewEchoService,
		NewEchoHandler,
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(server server.Server, handler *EchoHandler) {
	group := server.Gin.Group("/echo")

	group.POST("", handler.EchoTextHandler)
}
