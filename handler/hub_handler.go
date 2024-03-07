package handler

import (
	"fmt"
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	model "github.com/rirachii/golivechat/model"
)

type HubHandler struct {
	Hub *model.ChatroomsHub
}

func InitiateHub() (*model.ChatroomsHub, *HubHandler) {
	hub := &model.ChatroomsHub{
		ChatRooms:       make(map[model.RoomID]*model.Chatroom),
		UserChatrooms:   make(map[model.UserID]model.UserSetOfChatrooms),
		RegisterQueue:   make(chan *model.UserRequest),
		UnregisterQueue: make(chan *model.UserRequest),
	}

	handler := &HubHandler{
		Hub: hub,
	}

	return hub, handler
}

// HTMX endpoint
func (handler *HubHandler) HandleGetChatrooms(c echo.Context) error {

	// IDs for template
	const (
		chatroomsTemplateID = "hub-chatrooms"
		roomsLoopID         = "Rooms"
	)

	chatroomsData := map[string][]model.ChatroomData{}

	chatrooms := handler.Hub.ChatRooms
	for roomID, room := range chatrooms {
		roomName := room.GetName()
		roomData := model.ChatroomData{
			RoomName: roomName,
			RoomID:   roomID,
		}

		chatroomsData[roomsLoopID] = append(chatroomsData[roomsLoopID], roomData)
	}

	return c.Render(http.StatusOK, chatroomsTemplateID, &chatroomsData)

}

// HTMX endpoint
func (handler *HubHandler) HandleCreateRoom(c echo.Context) error {

	var newRoomRequest CreateRoomRequest
	err := c.Bind(&newRoomRequest)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// TODO check if room already exists
	echo.New().Logger.Debugf("Create room request received with data: %i", newRoomRequest)

	// IDs for template

	var (
		// TODO
		uid  model.UserID = model.UserID(newRoomRequest.UserID)
		rid  model.RoomID = model.RoomID(strconv.Itoa(len(handler.Hub.ChatRooms)))
		name string       = newRoomRequest.RoomName
	)

	userReq := model.UserRequest{
		UserID: model.UserID(uid),
		RoomID: model.RoomID(rid),
	}

	newRoom := model.NewChatroom(userReq, name)
	handler.Hub.AddandOpenRoom(newRoom)

	// send back to client to render new room

	const (
		chatroomsTemplateID = "hub-chatrooms"
		roomsLoopID         = "Rooms"
	)
	// one room
	chatroomsData := map[string][]model.ChatroomData{
		roomsLoopID: {
			model.ChatroomData{
				RoomID:   rid,
				RoomName: name,
			},
		},
	}

	// TODO make sure this is correct template.
	// Unsure why HandleGetChatrooms renders the same template.
	// ^^ both writes to the same div in htmx:
	// this allowed frontend to show the newly created room quicker,
	// once they receive response, rather than needing to manually refresh or wait to load

	return c.Render(http.StatusOK, chatroomsTemplateID, chatroomsData)
}

// Redirection
func (handler *HubHandler) HandleUserJoinRequest(c echo.Context) error {

	var joinRequest JoinRoomRequest
	err := c.Bind(&joinRequest)
	if err != nil {
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
	}

	// TODO make sure user is invited if room is private

	echo.New().Logger.Debugf("user join request received with data: %i", joinRequest)

	var (
		uid = joinRequest.UserID
		rid = c.Param("roomID")
	)

	userReq := &model.UserRequest{
		UserID: model.UserID(uid),
		RoomID: model.RoomID(rid),
	}

	handler.Hub.RegisterQueue <- userReq

	// set header for htmx to redirect from client-side
	chatroomRoute := fmt.Sprintf("/hub/chatroom/%s", rid)
	c.Response().Header().Set("HX-Redirect", chatroomRoute)
	return c.NoContent(http.StatusFound)
}

func (handler *HubHandler) HandleChatroomPage(c echo.Context) error {
	// TODO handle unauthorized access to page

	roomID := c.Param("roomID")
	getChatroom := handler.Hub.ChatRooms[model.RoomID(roomID)]

	chatroomData := getChatroom.GetChatroomData()

	const chatroomID = "chatroom"
	return c.Render(http.StatusOK, chatroomID, chatroomData)

}

func (handler *HubHandler) HandleUserLeave(c echo.Context) error {
	// handler
	var leaveRequest LeaveRoomRequest
	err := c.Bind(&leaveRequest)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	uid, rid := leaveRequest.UserID, leaveRequest.RoomID

	userReq := &model.UserRequest{
		UserID: model.UserID(uid),
		RoomID: model.RoomID(rid),
	}

	handler.Hub.UnregisterQueue <- userReq

	return nil

}

func (handler *HubHandler) HandleChatroomConnection(c echo.Context) error {
	// Websocket connection, should be
	// c.Echo().Logger.Print(c.Request(), c.Request().Body)

	if !c.IsWebSocket() {
		errMsg := "expected Websocket connection, but was not"
		c.Logger().Debug(errMsg)
		return c.NoContent(http.StatusUpgradeRequired)
	}

	var connReq RoomRequest
	bindErr := c.Bind(&connReq)
	if bindErr != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	uid, rid := connReq.UserID, connReq.RoomID

	userReq := model.UserRequest{
		UserID: model.UserID(uid),
		RoomID: model.RoomID(rid),
	}

	// log.Printf("new user req: [%v]", userReq)

	// check user ID

	getChatroom := handler.Hub.GetChatroom(userReq.RoomID)

	connErr := getChatroom.AcceptConnection(c, userReq)
	if connErr != nil {
		// TODO handle err, tell client what error is maybe
		c.Logger().Print("connection error", connErr)
		return c.NoContent(http.StatusInternalServerError)
	}

	// return nil. Do not write to context after websocket connection is successful.
	// will get some "hijacked connection" warning
	return nil

}

func (handler *HubHandler) HandleFetchChatroomHistory(c echo.Context) error {

	roomID := c.Param("roomID")
	getChatroom := handler.Hub.GetChatroom(model.RoomID(roomID))

	if getChatroom == nil {
		// invalid request
		return c.NoContent(http.StatusBadRequest)

	}

	chatroomHistoryData := getChatroom.GetChatroomHistory(c)

	const msgsTemplateID = "many-messages"
	return c.Render(http.StatusOK, msgsTemplateID, chatroomHistoryData)

}

// i dont think this is used at all
func (handler *HubHandler) HandleChatroomMessage(c echo.Context) error {

	roomID := c.Param("roomID")
	getChatroom := handler.Hub.GetChatroom(model.RoomID(roomID))
	if getChatroom == nil {
		// invalid request
		return c.NoContent(http.StatusBadRequest)

	}

	return getChatroom.ReceiveNewMessage(c)
}
