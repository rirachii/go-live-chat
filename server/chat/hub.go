package chat

import (
	"fmt"
	"html/template"
	"net/http"

	"strconv"

	echo "github.com/labstack/echo/v4"
	// websocket "nhooyr.io/websocket"
)

type UserID string
type RoomID string

type User struct {
	// WebSocket *websocket.Conn
	UserID UserID
	RoomID RoomID
}

type HubHandler struct {
	Hub      *ChatroomsHub
	HTMLTemplate *template.Template
}

type ChatroomsHub struct {
	HTMLTemplate    *template.Template
	ChatRooms       map[RoomID]*ChatRoom
	UserChatrooms   map[UserID]SetOfChatrooms // userid -> their chatrooms (room ids)
	RegisterQueue   chan *User
	UnregisterQueue chan *User
}

type SetOfChatrooms struct {
	Chatrooms map[RoomID]bool
}

func (rooms *SetOfChatrooms) RegisterRoom(roomID RoomID) {
	// TODO be aware does not check if the room is already registered to them
	rooms.Chatrooms[roomID] = true
}
func (rooms *SetOfChatrooms) UnregisterRoom(roomID RoomID) {
	delete(rooms.Chatrooms, roomID)
}

func InitiateHub(t *template.Template) (*ChatroomsHub, *HubHandler) {

	hub := &ChatroomsHub{
		ChatRooms:       make(map[RoomID]*ChatRoom),
		UserChatrooms:   make(map[UserID]SetOfChatrooms),
		RegisterQueue:   make(chan *User),
		UnregisterQueue: make(chan *User),
	}

	handler := &HubHandler{
		Hub: hub,
	}

	return hub, handler
}

func (hub *ChatroomsHub) Run() {

	for {
		select {
		case client := <-hub.RegisterQueue:
			// register them to the approrpiate chat

			clientID, roomID := client.UserID, client.RoomID

			userChatrooms, ok := hub.UserChatrooms[clientID]
			if !ok {
				hub.UserChatrooms[clientID] = SetOfChatrooms{
					Chatrooms: make(map[RoomID]bool),
				}
				userChatrooms = hub.UserChatrooms[clientID]
				// log.Printf("user [%s] not found", clientID)
			}

			userChatrooms.RegisterRoom(roomID)

			// chatRoom := hub.ChatRooms[roomID]
			// // chatRoom.JoinQueue <- client

			echo.New().Logger.Printf("Registering %s to room %s", clientID, roomID)
			echo.New().Logger.Printf("User [%s] rooms: %i", clientID, userChatrooms)

		case client := <-hub.UnregisterQueue:

			clientID, roomID := client.UserID, client.RoomID

			// userChatrooms, ok := hub.UserChatrooms[clientID]
			// if !ok {
			// 	log.Printf("user [%s] not found", clientID)
			// }

			// userChatrooms.UnregisterRoom(roomID)

			// chatRoom := hub.ChatRooms[roomID]

			// chatRoom.LeaveQueue <- client

			echo.New().Logger.Printf("Unregistering %s from room %s", clientID, roomID)

		}
	}

}

func (hub *ChatroomsHub) AddandOpenRoom(newChatRoom *ChatRoom) error {

	roomID := newChatRoom.RoomID
	hub.ChatRooms[roomID] = newChatRoom
	go newChatRoom.Open()

	// TODO check errors
	return nil

}
func (hub ChatroomsHub) GetChatroom(roomID RoomID) *ChatRoom {

	getChatroom, ok := hub.ChatRooms[roomID]

	if !ok {
		return nil
	}

	return getChatroom

}

type HubChatRoom struct {
	RoomID   RoomID
	RoomName string
}

func (handler *HubHandler) HandleGetChatrooms(c echo.Context) error {

	// name of the loop in the template
	chatrooms := handler.Hub.ChatRooms

	templateData := map[string][]HubChatRoom{}
	const roomsContainerID = "Rooms"

	for roomID, room := range chatrooms {

		roomName := room.RoomName

		roomData := HubChatRoom{
			RoomName: roomName,
			RoomID:   roomID,
		}

		templateData[roomsContainerID] = append(templateData[roomsContainerID], roomData)

		// htmlTemplate := template.
	}

	const chatroomTemplateName = "hub-chatrooms"
	return c.Render(http.StatusOK, chatroomTemplateName, &templateData)

}

func (handler *HubHandler) HandleFetchChatroomHistory(c echo.Context) error {

	roomID := c.Param("roomID")
	getChatroom := handler.Hub.GetChatroom(RoomID(roomID))
	return getChatroom.HandleChatroomHistory(c)
}

type CreateRoomRequest struct {
	// TODO: user id instead of user display name
	UserID   string `json:"display-name"`
	RoomName string `json:"room-name"`
}

func (handler *HubHandler) HandleCreateRoom(c echo.Context) error {

	var newRoomRequest CreateRoomRequest

	err := c.Bind(&newRoomRequest)

	if err != nil {
		return echo.ErrBadRequest
	}

	echo.New().Logger.Printf("Create room request received with data: %i", newRoomRequest)

	// TODO check if room already exists

	// id := newRoomRequest.RoomID
	var rid RoomID = RoomID(strconv.Itoa(len(handler.Hub.ChatRooms)))
	uid := newRoomRequest.UserID
	name := newRoomRequest.RoomName

	newRoom := NewChatRoom(UserID(uid), rid, name, handler.HTMLTemplate)

	handler.Hub.AddandOpenRoom(newRoom)

	// return echo.ErrNotImplemented

	templateData := map[string][]map[string]string{
		"Rooms": {
			{"RoomName": name, "RoomID": string(rid)},
		},
	}

	const hubChatroomsTemplateID = "hub-chatrooms"
	return c.Render(http.StatusOK, hubChatroomsTemplateID, templateData)
}

type RegisterRequest struct {
	// TODO user ID instaed of user display name
	UserID string `json:"display-name"`
	RoomID string `json:"room-id"`
}

func (handler *HubHandler) HandleUserJoinRequest(c echo.Context) error {

	//serve it
	// in another handler, accept websocket

	var registerRequest RegisterRequest
	err := c.Bind(&registerRequest)
	if err != nil {
		// TODO handle error
		_ = err
	}

	echo.New().Logger.Printf("user register request received with data: %i", registerRequest)

	uid := registerRequest.UserID
	rid := c.Param("roomID")

	// roomHandler := handler.Hub.chatRooms[rid]

	user := &User{
		// WebSocket: nil,
		UserID: UserID(uid),
		RoomID: RoomID(rid),
	}

	handler.Hub.RegisterQueue <- user

	chatroomRoute := fmt.Sprintf("/hub/chatroom/%s", rid)

	c.Response().Header().Set("HX-Location", chatroomRoute)
	return c.NoContent(http.StatusFound)
}

func (handler *HubHandler) HandleChatroomPage(c echo.Context) error {

	// TODO handle unauthorized access to page, they should be register Queue

	roomID := c.Param("roomID")
	getChatroom := handler.Hub.ChatRooms[RoomID(roomID)]

	return getChatroom.RenderChatroomPage(c)

}

type UnregisterRequest struct {
	// TODO user ID instaed of user display name
	UserID string `json:"display-name"`
	RoomID string `json:"room-id"`
}

func (handler *HubHandler) HandleUserLeave(c echo.Context) error {

	// handler
	var unregisterRequest UnregisterRequest
	err := c.Bind(&unregisterRequest)
	if err != nil {
		// TODO handle error
		_ = err

	}

	uid, rid := unregisterRequest.UserID, unregisterRequest.RoomID

	user := &User{
		UserID: UserID(uid),
		RoomID: RoomID(rid),
	}

	handler.Hub.UnregisterQueue <- user

	return nil

}

func (handler *HubHandler) HandleChatroomWSConnection(c echo.Context) error {
	// Websocket connection, should be
	// c.Echo().Logger.Print(c.Request(), c.Request().Body)

	// userID := c.Param("UserID")
	roomID := c.Param("roomID")

	// check user ID

	getChatroom := handler.Hub.ChatRooms[RoomID(roomID)]

	return getChatroom.HandleNewConnection(c)

}

func (handler *HubHandler) HandleChatroomMessage(c echo.Context) error {

	roomID := c.Param("roomID")
	getChatroom := handler.Hub.GetChatroom(RoomID(roomID))

	return getChatroom.HandleNewMessage(c)
}
