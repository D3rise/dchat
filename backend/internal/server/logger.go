package server

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"go.uber.org/zap"
)

func DecorateServerWithLogger(server Server, log *zap.Logger) Server {
	server.Gin.Use(ginzap.Ginzap(log, time.RFC3339, true))
	server.Gin.Use(ginzap.RecoveryWithZap(log, true))

	return server
}
