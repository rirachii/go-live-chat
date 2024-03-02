package chat

import (
	"github.com/labstack/echo/v4"
	users "github.com/rirachii/golivechat/users"
	ws "nhooyr.io/websocket"
)

type ChatRoomHandler struct {
	ChatRoom *LiveChatRoom
}

func NewChatRoomHandler(chatRoom *LiveChatRoom) *ChatRoomHandler{
	return &ChatRoomHandler{
		ChatRoom: chatRoom,
	}
}

type LiveChatRoom struct {
	roomID        string
	roomName	string	
	chatHistory   []Message
	clientConnections map[*ws.Conn]Chatter
	joinQueue         chan *Client
	leaveQueue        chan *Client
}

type Chatter struct {
	user         *users.User
	conn         *ws.Conn
	role         string
	roomID       string
	messageQueue chan *Message
}

type Message struct {
	From    Chatter
	Content string
}

func NewChatRoom(id string) *LiveChatRoom {

	instance := &LiveChatRoom{
		roomID:        id,
		chatHistory:   []Message{},
		clientConnections: make(map[*ws.Conn]Chatter, 5), //max five for now

	}

	return instance
}

func (chatroom *LiveChatRoom) Open() {
	for {
		select {
		case chatter := <- chatroom.joinQueue:
			// TODO add to this chat
			_ = chatter

		
		case chatter := <- chatroom.leaveQueue:
			// TODO leave from this chat
			_ = chatter
		
		}
	}

}

func (ChatRoom LiveChatRoom) ID() string {

	return ChatRoom.roomID
}

func (ChatRoom LiveChatRoom) ChatHistory() []Message {

	return ChatRoom.chatHistory
}


func HandleCreateChatRoom(c echo.Context) {

}

func HandleChatSend(c echo.Context) {

	// conn := ws.NetConn()

}
