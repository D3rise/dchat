package rtc

import (
	"encoding/json"
	"net/http"

	"github.com/D3rise/dchat/internal/modules/rtc/structs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetRtcWebsocketHandler(rtcService RtcService, log *zap.Logger) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		go handleConn(conn, rtcService, log)
	}
}

func handleConn(conn *websocket.Conn, rtcService RtcService, log *zap.Logger) {
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Error(err.Error())
			return
		}

		if messageType != websocket.TextMessage {
			conn.WriteMessage(websocket.TextMessage, []byte("you aren't supposed to be here"))
			return
		}

		var data structs.WsMessage
		err = json.Unmarshal(message, &data)
		if err != nil {
			log.Error(err.Error())
			return
		}

		err = rtcService.HandleMessage(conn, data.Type, data.Data)
	}
}
