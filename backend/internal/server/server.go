package server

import (
	"context"

	"github.com/D3rise/dchat/internal/interfaces"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Server struct {
	Gin *gin.Engine
}

func NewHandler(lc fx.Lifecycle, envService interfaces.EnvService) Server {
	r := gin.Default()
	server := Server{
		Gin: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go r.Run(envService.GetListenAddr())
			return nil
		},
	})

	return server
}
