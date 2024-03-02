package chat

import (
	"fmt"
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
)

type HubHandler struct {
	Hub *ChatroomsHub
}

type hubChatroom struct {
	RoomID   RoomID
	RoomName string
}

func (handler *HubHandler) HandleGetChatrooms(c echo.Context) error {

	// name of the loop in the template
	chatrooms := handler.Hub.ChatRooms

	chatroomsData := map[string][]hubChatroom{}
	const roomsLoopID = "Rooms"

	for roomID, room := range chatrooms {

		roomName := room.getName()
		roomData := hubChatroom{
			RoomName: roomName,
			RoomID:   roomID,
		}

		chatroomsData[roomsLoopID] = append(chatroomsData[roomsLoopID], roomData)
	}

	const chatroomsTemplateID = "hub-chatrooms"
	return c.Render(http.StatusOK, chatroomsTemplateID, &chatroomsData)

}

type createRoomRequest struct {
	// TODO: user id instead of user display name
	UserID   string `json:"display-name"`
	RoomName string `json:"room-name"`
}

func (handler *HubHandler) HandleCreateRoom(c echo.Context) error {

	var newRoomRequest createRoomRequest

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

	newRoom := NewChatroom(UserID(uid), rid, name)
	handler.Hub.AddandOpenRoom(newRoom)

	// return echo.ErrNotImplemented

	templateData := map[string][]map[string]string{
		"Rooms": {
			{"RoomName": name, "RoomID": string(rid)},
		},
	}

	// TODO make sure this is correct template. Unsure why HandleGetChatrooms renders teh same template.
	const hubChatroomsTemplateID = "hub-chatrooms"
	return c.Render(http.StatusOK, hubChatroomsTemplateID, templateData)
}

type registerRoomRequest struct {
	// TODO user ID instaed of user display name
	UserID string `json:"display-name"`
	RoomID string `json:"room-id"`
}

func (handler *HubHandler) HandleUserJoinRequest(c echo.Context) error {

	var registerRequest registerRoomRequest
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

	// set header for htmx to redirect from client-side
	c.Response().Header().Set("HX-Redirect", chatroomRoute)
	return c.NoContent(http.StatusFound)
}

func (handler *HubHandler) HandleChatroomPage(c echo.Context) error {

	// TODO handle unauthorized access to page

	roomID := c.Param("roomID")
	getChatroom := handler.Hub.ChatRooms[RoomID(roomID)]

	return getChatroom.RenderChatroomPage(c)

}

type unRegisterRoomRequest struct {
	// TODO user ID instaed of user display name
	UserID string `json:"display-name"`
	RoomID string `json:"room-id"`
}

func (handler *HubHandler) HandleUserLeave(c echo.Context) error {

	// handler
	var unregisterRequest unRegisterRoomRequest
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

func (handler *HubHandler) HandleFetchChatroomHistory(c echo.Context) error {

	roomID := c.Param("roomID")
	getChatroom := handler.Hub.getChatroom(RoomID(roomID))

	return getChatroom.HandleChatroomLogs(c)
}

func (handler *HubHandler) HandleChatroomMessage(c echo.Context) error {

	roomID := c.Param("roomID")
	getChatroom := handler.Hub.getChatroom(RoomID(roomID))

	return getChatroom.HandleNewMessage(c)
}
