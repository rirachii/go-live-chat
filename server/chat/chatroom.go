package chatroom

import (
	"github.com/labstack/echo/v4"
	users "github.com/rirachii/golivechat/users"
	ws "nhooyr.io/websocket"
)

type HubHandler struct {
	hub *ChatHub
}


func (handler *HubHandler) CreateRoom(c echo.Context){

}

func (handler *HubHandler) RegisterUser(c echo.Context){

}
func (handler *HubHandler) UnregisterUser(c echo.Context){

}

type ChatHub struct {
	chatRooms       map[string]*LiveChatRoom
	registerQueue   chan *users.User
	unregisterQueue chan *users.User
}

func InitiateHub() *ChatHub {
	return &ChatHub{
		chatRooms:       make(map[string]*LiveChatRoom),
		registerQueue:   make(chan *users.User),
		unregisterQueue: make(chan *users.User),
	}
}

type LiveChatRoom struct {
	chatRoomID        string
	chatRoomHistory   []Message
	clientConnections map[*ws.Conn]Chatter
}

type Message struct {
	From    Chatter
	Content string
	roomID  string
}

type Chatter struct {
	user         *users.User
	conn         *ws.Conn
	role         string
	roomID       string
	messageQueue chan *Message
}

func NewChatRoom(id string) *LiveChatRoom {

	instance := &LiveChatRoom{
		clientConnections: make(map[*ws.Conn]Chatter, 5), //max five for now
		chatRoomID:        id,
		chatRoomHistory:   []Message{},
	}

	return instance
}

func (ChatRoom LiveChatRoom) ID() string {

	return ChatRoom.chatRoomID
}

func (ChatRoom LiveChatRoom) ChatHistory() []Message {

	return ChatRoom.chatRoomHistory
}

func (*LiveChatRoom) Open() {

}

func HandleCreateChatRoom(c echo.Context) {

}

func HandleChatSend(c echo.Context) {

	// conn := ws.NetConn()

}
