package rtc

import (
	"github.com/D3rise/dchat/internal/server"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("rtc",
	fx.Provide(NewRtcService),
	fx.Invoke(registerHooks),
)

func registerHooks(s server.Server, rtcService RtcService, log *zap.Logger) {
	group := s.Gin.Group("/rtc")

	group.Any("/ws", GetRtcWebsocketHandler(rtcService, log))
}
