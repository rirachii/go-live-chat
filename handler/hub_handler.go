package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	model "github.com/rirachii/golivechat/model"
	chat "github.com/rirachii/golivechat/model/chat"
	hub_model "github.com/rirachii/golivechat/model/hub"
	hub_svc "github.com/rirachii/golivechat/service/hub"
	"github.com/rirachii/golivechat/service/db"

	chat_template "github.com/rirachii/golivechat/templates/chat"
	hub_template "github.com/rirachii/golivechat/templates/hub"
)

type HubHandler struct {
	Hub *hub_model.ChatroomsHub
}

func InitiateHub() (*hub_model.ChatroomsHub, *HubHandler) {

	hub := hub_model.NewChatroomsHub()

	handler := &HubHandler{
		Hub: hub,
	}

	return hub, handler
}
func createHubService() (*hub_svc.HubService, error){

	db, err := db.ConnectDatabase(); if err != nil{
		return nil, err
	}

	hubRepo := hub_svc.NewHubRepository(db.DB())
	hubSvc := hub_svc.NewHubService(hubRepo)

	return &hubSvc, nil


}

// HTMX endpoint
func (handler *HubHandler) HandleGetChatrooms(c echo.Context) error {

	db, err := db.ConnectDatabase()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	hubRepo := hub_svc.NewHubRepository(db.DB())
	hubSvc := hub_svc.NewHubService(hubRepo)

	var userID model.UserID = ""

	jwtClaims, err := GetJWTClaims(c)
	if err == nil {
		userID = model.UID(jwtClaims.GetUID())
	}

	chatroomsReq := &hub_model.GetChatroomsRequest{
		UserID: userID,
	}

	res, err := hubSvc.GetRoomsPublic(c.Request().Context(), chatroomsReq)
	if err != nil {
		log.Println("error with service")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// IDs for template
	var (
		chatroomsTemplateName = hub_template.HubChatrooms.TemplateName
	)

	chatroomsData := hub_template.TemplateHubChatrooms{
		Rooms: make([]hub_template.ChatroomTemplateData, 0),
	}
	chatrooms := res
	for _, room := range chatrooms.Chatrooms {

		roomID := strconv.Itoa(room.RoomID)
		roomName :=  room.RoomName

		roomData := hub_template.PrepareChatroomData(model.RID(roomID), roomName)

		chatroomsData.Rooms = append(chatroomsData.Rooms, roomData)
	}

	return c.Render(http.StatusOK, chatroomsTemplateName, &chatroomsData)

}

// HTMX endpoint
func (handler *HubHandler) HandleCreateRoom(c echo.Context) error {

	jwtClaims, err := GetJWTClaims(c)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/login")
		return c.NoContent(http.StatusUnauthorized)
	}

	userID := jwtClaims.GetUID()

	var createRoomReq hub_model.CreateRoomRequest
	err = c.Bind(&createRoomReq)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	createRoomReq.UserID = model.UID(userID)

	// TODO check if room already exists
	echo.New().Logger.Debugf("Create room request received with data: %+v", createRoomReq)

	var (
		uid  model.UserID = createRoomReq.UserID
		rid  model.RoomID = model.RID(strconv.Itoa(len(handler.Hub.Chatrooms())))
		name string       = createRoomReq.RoomName
	)

	userReq := model.UserRequest{
		UserID: uid,
		RoomID: rid,
	}

	newRoom := chat.NewChatroom(userReq, name)
	handler.Hub.AddandOpenRoom(newRoom)

	// send back to client to render new room

	var (
		chatroomsTemplateID = hub_template.HubChatrooms.TemplateName
	)

	chatroomData := hub_template.PrepareChatroomData(rid, name)

	// one room
	templateData := hub_template.TemplateHubChatrooms{
		Rooms: []hub_template.ChatroomTemplateData{chatroomData},
	}

	// TODO make sure this is correct template.
	// Unsure why HandleGetChatrooms renders the same template.
	// ^^ both writes to the same div in htmx:
	// this allowed frontend to show the newly created room quicker,
	// once they receive response, rather than needing to manually refresh or wait to load

	return c.Render(http.StatusOK, chatroomsTemplateID, templateData)
}

// Redirection
func (handler *HubHandler) HandleUserJoinRequest(c echo.Context) error {

	var joinRequest hub_model.JoinRoomRequest
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

	handler.Hub.Register(userReq)

	// set header for htmx to redirect from client-side
	chatroomRoute := fmt.Sprintf("/hub/chatroom/%s", rid)
	c.Response().Header().Set("HX-Redirect", chatroomRoute)
	return c.NoContent(http.StatusFound)
}

func (handler *HubHandler) HandleChatroomPage(c echo.Context) error {
	// TODO handle unauthorized access to page

	roomID := c.Param("roomID")
	chatroom := handler.Hub.Chatroom(model.RID(roomID))

	chatroomInfo := chatroom.Info()

	const chatroomID = "chatroom"
	return c.Render(http.StatusOK, chatroomID, chatroomInfo)

}

func (handler *HubHandler) HandleUserLeave(c echo.Context) error {
	// handler
	var leaveRequest hub_model.LeaveRoomRequest
	err := c.Bind(&leaveRequest)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	uid, rid := leaveRequest.UserID, leaveRequest.RoomID

	userReq := &model.UserRequest{
		UserID: model.UserID(uid),
		RoomID: model.RoomID(rid),
	}

	handler.Hub.Unregister(userReq)

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

	var connReq hub_model.RoomRequest
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

	chatroom := handler.Hub.Chatroom(userReq.RoomID)

	connErr := chatroom.AcceptConnection(c, userReq)
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
	chatroom := handler.Hub.Chatroom(model.RoomID(roomID))

	if chatroom == nil {
		// invalid request
		return c.NoContent(http.StatusBadRequest)

	}

	chatroomLogs := chatroom.ChatLogs()

	chatLogsData := chat_template.TemplateManyMessages{
		ChatMessages: make([]chat_template.TemplateSingleMessage, 0),
	}

	for _, chatMessage := range chatroomLogs {

		log.Printf(`msg: "[%s]" by: [%s]`, chatMessage.Content, chatMessage.From)

		singleMsgData := chat_template.PrepareMessage(
			chat_template.WebsocketDivID,
			false,
			string(chatMessage.From),
			chatMessage.Content,
		)

		chatLogsData.ChatMessages = append(chatLogsData.ChatMessages, singleMsgData)
	}

	msgsTemplate := chat_template.ManyMessages.TemplateName
	return c.Render(http.StatusOK, msgsTemplate, chatLogsData)

}
