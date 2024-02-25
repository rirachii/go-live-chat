package chat

import (
	"log"
	"net/http"

	"strconv"

	"github.com/labstack/echo/v4"
	ws "nhooyr.io/websocket"
)

type Client struct {
	UserID string
	RoomID string
}

type ChatHub struct {
	chatRooms      map[string]*ChatRoomHandler
	joinRoomQueue  chan *Client
	leaveRoomQueue chan *Client
}

func InitiateHub() (*ChatHub, *HubHandler) {

	hub := &ChatHub{
		chatRooms:      make(map[string]*ChatRoomHandler),
		joinRoomQueue:  make(chan *Client),
		leaveRoomQueue: make(chan *Client),
	}

	handler := &HubHandler{
		Hub: hub,
	}

	return hub, handler
}

func (h *ChatHub) Run() {

	for {
		select {
		case client := <-h.joinRoomQueue:
			// register them to the approrpiate chat
			// <client
			log.Println(client)

		case client := <-h.leaveRoomQueue:
			log.Println(client)

		}
	}

}

func (h *ChatHub) AddRoom(chatRoom *LiveChatRoom) error {

	echo.New().Logger.Printf("adding room")

	roomId := chatRoom.ID()

	chatRoomHandler := NewChatRoomHandler(chatRoom)

	h.chatRooms[roomId] = chatRoomHandler

	echo.New().Logger.Print(len(h.chatRooms))

	// TODO check errors

	return nil

}

type HubHandler struct {
	Hub *ChatHub
}

type HubChatRoom struct {
	RoomName string
	RoomID   string
}

func (handler *HubHandler) GetHubRooms(c echo.Context) error {

	echo.New().Logger.Printf("gettin rooms")

	chatroomTemplateName := "hub-chatroom"

	chatrooms := handler.Hub.chatRooms

	templateData := map[string][]HubChatRoom{}
	for id, h := range chatrooms {

		name := h.ChatRoom.roomName

		roomData := HubChatRoom{
			RoomName: name,
			RoomID:   id,
		}

		templateData["Rooms"] = append(templateData["Rooms"], roomData)

		// htmlTemplate := template.
	}

	return c.Render(http.StatusOK, chatroomTemplateName, &templateData)

}

type createRoomRequest struct {
	RoomName string `form:"room-name"`
}

func (handler *HubHandler) CreateRoom(c echo.Context) error {

	var newRoomRequest createRoomRequest

	err := c.Bind(&newRoomRequest)

	if err != nil {
		return echo.ErrBadRequest
	}

	// TODO check if room already exists

	// id := newRoomRequest.RoomID
	id := strconv.Itoa(len(handler.Hub.chatRooms))
	name := newRoomRequest.RoomName

	newRoom := &LiveChatRoom{
		roomID:            id,
		roomName:          name,
		chatHistory:       []Message{},
		clientConnections: make(map[*ws.Conn]Chatter),
		joinQueue:         make(chan *Client),
		leaveQueue:        make(chan *Client),
	}

	handler.Hub.AddRoom(newRoom)

	// return echo.ErrNotImplemented
	return c.NoContent(http.StatusOK)
}

type userRequest struct {
	UserID string `json:"userId"`
	RoomID string `json:"roomId"`
}

func (handler *HubHandler) ForwardUser(c echo.Context) error {

	// handler
	var registerRequest userRequest
	err := c.Bind(&registerRequest)
	if err != nil {
		// TODO handle error
		_ = err
	}

	uid, rid := registerRequest.UserID, registerRequest.RoomID

	client := &Client{
		UserID: uid,
		RoomID: rid,
	}

	handler.Hub.joinRoomQueue <- client

	return nil
}

func (handler *HubHandler) UnregisterUser(c echo.Context) error {

	// handler
	var registerRequest userRequest
	err := c.Bind(&registerRequest)
	if err != nil {
		// TODO handle error
		_ = err

	}

	uid, rid := registerRequest.UserID, registerRequest.RoomID

	client := &Client{
		UserID: uid,
		RoomID: rid,
	}

	handler.Hub.joinRoomQueue <- client

	return nil

}
