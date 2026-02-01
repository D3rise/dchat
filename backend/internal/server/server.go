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
			addr := envService.GetListenAddr()
			tlsCertPath := envService.GetTLSCertPath()
			tlsKeyPath := envService.GetTLSKeyPath()

			if tlsCertPath != "" && tlsKeyPath != "" {
				go r.RunTLS(addr, envService.GetTLSCertPath(), envService.GetTLSKeyPath())
			} else {
				go r.Run(addr)
			}

			return nil
		},
	})

	return server
}
