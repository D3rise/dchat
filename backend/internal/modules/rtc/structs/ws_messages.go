package structs

import "github.com/D3rise/dchat/internal/modules/rtc/enums"

type WsMessage struct {
	Type enums.WsMessageType `json:"_type"`
	Data any                 `json:"data"`
}

type EchoClientMessage struct {
	Text string `json:"text"`
}

type EchoServerMessage struct {
	Text string `json:"text"`
}

type AllRoomUsersServerMessage struct {
	Users []string `json:"users"`
}

type JoinRoomClientMessage struct {
	Username string `json:"username"`
	RoomID   string `json:"room_id"`
}

type SignalClientMessage struct {
	Username         string `json:"username"`
	UsernameToSignal string `json:"username_to_signal"`
	Signal           string `json:"signal"`
}

type SignalServerMessage struct {
	Username string `json:"username"`
	Signal   string `json:"signal"`
}

type ReturningSignalClientMessage struct {
	Username         string `json:"username"`
	UsernameToSignal string `json:"username_to_signal"`
	Signal           string `json:"signal"`
}

type ReturningSignalServerMessage struct {
	Username string `json:"username"`
	Signal   string `json:"signal"`
}
