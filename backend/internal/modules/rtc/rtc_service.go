package rtc

import (
	"encoding/json"
	"fmt"

	"github.com/D3rise/dchat/internal/modules/rtc/enums"
	"github.com/D3rise/dchat/internal/modules/rtc/structs"
	"github.com/gorilla/websocket"
)

type User struct {
	Conn     *websocket.Conn
	Username string
	RoomID   string
}

type RtcService struct {
	// username => user
	users map[string]User

	// room name => []joined users
	rooms map[string][]User
}

func NewRtcService() RtcService {
	return RtcService{
		users: make(map[string]User),
		rooms: make(map[string][]User),
	}
}

func (r *RtcService) HandleMessage(conn *websocket.Conn, _type enums.WsMessageType, data any) error {
	switch _type {
	case enums.WsMessageTypeEcho:
		{
			if data, err := unmarshalWsMessage[structs.EchoClientMessage](data); err == nil {
				return r.Echo(conn, data)
			}
			break
		}

	case enums.WsMessageTypeJoinRoom:
		{
			if data, err := unmarshalWsMessage[structs.JoinRoomClientMessage](data); err == nil {
				return r.JoinRoom(conn, data)
			}
			break
		}

	case enums.WsMessageTypeSignal:
		{
			if data, err := unmarshalWsMessage[structs.SignalClientMessage](data); err == nil {
				return r.Signal(conn, data)
			}
		}

	case enums.WsMessageTypeReturningSignal:
		{
			if data, err := unmarshalWsMessage[structs.ReturningSignalClientMessage](data); err == nil {
				return r.ReturningSignal(conn, data)
			}
		}
	}

	return nil
}

func (r *RtcService) Echo(ws *websocket.Conn, data structs.EchoClientMessage) error {
	result, err := marshalWsMessage(enums.WsMessageTypeEcho, structs.EchoServerMessage{Text: data.Text})
	if err != nil {
		return err
	}
	return ws.WriteMessage(websocket.TextMessage, result)
}

func (r *RtcService) JoinRoom(ws *websocket.Conn, data structs.JoinRoomClientMessage) error {
	user := User{
		Conn:     ws,
		Username: data.Username,
		RoomID:   data.RoomID,
	}

	r.sendRoomUsers(ws, data.RoomID)

	r.users[data.Username] = user
	r.rooms[data.RoomID] = append(r.rooms[data.RoomID], user)

	return nil
}

// TODO: add a check for the room
func (r *RtcService) Signal(ws *websocket.Conn, data structs.SignalClientMessage) error {
	userToSignal, ok := r.users[data.UsernameToSignal]
	if !ok {
		return fmt.Errorf("user does not exist")
	}

	signalMessage := structs.SignalServerMessage{Signal: data.Signal, Username: data.Username}
	signalMessageJson, err := marshalWsMessage(enums.WsMessageTypeSignal, signalMessage)
	if err != nil {
		return err
	}

	return userToSignal.Conn.WriteMessage(websocket.TextMessage, signalMessageJson)
}

func (r *RtcService) ReturningSignal(ws *websocket.Conn, data structs.ReturningSignalClientMessage) error {
	userToSignal, ok := r.users[data.UsernameToSignal]
	if !ok {
		return fmt.Errorf("user does not exist")
	}

	signalMessage := structs.ReturningSignalServerMessage{Signal: data.Signal, Username: data.Username}
	signalMessageJson, err := marshalWsMessage(enums.WsMessageTypeReturningSignal, signalMessage)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return userToSignal.Conn.WriteMessage(websocket.TextMessage, signalMessageJson)
}

func (r *RtcService) sendRoomUsers(ws *websocket.Conn, room string) error {
	users := r.rooms[room]
	usernames := make([]string, len(users))
	for i, user := range users {
		usernames[i] = user.Username
	}

	result, err := marshalWsMessage(enums.WsMessageTypeRoomUserList, structs.AllRoomUsersServerMessage{Users: usernames})
	if err != nil {
		return err
	}

	return ws.WriteMessage(websocket.TextMessage, result)
}

func marshalWsMessage[T any](_type enums.WsMessageType, data T) ([]byte, error) {
	fullMessage := structs.WsMessage{Type: _type, Data: data}
	jsonData, err := json.Marshal(fullMessage)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jsonData, nil
}

func unmarshalWsMessage[T any](data any) (T, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return *new(T), err
	}

	var unmarshalledData T
	err = json.Unmarshal(jsonData, &unmarshalledData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return *new(T), err
	}

	return unmarshalledData, nil
}
