package handler

import (
	"context"
	errors "errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"

	echo "github.com/labstack/echo/v4"

	model "github.com/rirachii/golivechat/model"
	chatroom_model "github.com/rirachii/golivechat/model/chat"
	hub_model "github.com/rirachii/golivechat/model/hub"

	db "github.com/rirachii/golivechat/service/db"
	hub_svc "github.com/rirachii/golivechat/service/hub"

	chatroom_template "github.com/rirachii/golivechat/templates/chatroom"
	hub_template "github.com/rirachii/golivechat/templates/hub"
)

type HubHandler struct {
	ChatroomsHub *hub_model.ChatroomsHub
}

func (h *HubHandler) Hub() *hub_model.ChatroomsHub {
	return h.ChatroomsHub
}

func InitiateHubHandler() (*HubHandler, error) {

	hubChatrooms, err := getChatroomsDB()
	if err != nil {
		return nil, err
	}

	hub := hub_model.CreateChatroomsHub(hubChatrooms)

	handler := &HubHandler{
		ChatroomsHub: hub,
	}

	go handler.Hub().Run()
	go handler.openHubChatrooms()


	return handler, nil
}

// HTMX endpoint
func (handler *HubHandler) HandleGetPublicChatrooms(c echo.Context) error {

	hubSvc, err := createHubService()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	chatroomsReq := hub_model.GetPublicChatroomsRequest{
		IsPublic: true,
		IsActive: true,
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
	for _, room := range chatrooms {

		roomID := room.RoomID
		roomName := room.RoomName
		roomData := hub_template.PrepareChatroomData(roomID, roomName)

		chatroomsData.Rooms = append(chatroomsData.Rooms, roomData)
	}

	return c.Render(http.StatusOK, chatroomsTemplateName, &chatroomsData)

}

// HTMX endpoint
func (handler *HubHandler) HandleCreateRoom(c echo.Context) error {

	userID, err := GetJWTUserID(c)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/login")
		return c.NoContent(http.StatusUnauthorized)
	}

	hubSvc, svcErr := createHubService()
	if svcErr != nil {
		return c.String(http.StatusInternalServerError, svcErr.Error())
	}

	var createRoomReq hub_model.CreateRoomRequest
	err = c.Bind(&createRoomReq)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	createRoomReq.UserID = model.UID(userID)
	createRoomReq.IsPublic = true
	createRoomReq.IsActive = true

	roomData, creationErr := hubSvc.CreateRoom(c.Request().Context(), createRoomReq)
	if creationErr != nil {
		return c.String(http.StatusInternalServerError, creationErr.Error())
	}

	// TODO check if room already exists
	echo.New().Logger.Printf("Create room request received with data: %+v", createRoomReq)

	var (
		uid  model.UserID = createRoomReq.UserID
		rid  model.RoomID = roomData.RoomID
		name string       = roomData.RoomName
	)

	roomInfo := model.ChatroomInfo{
		RoomID:    rid,
		RoomName:  name,
		RoomOwner: uid,
		IsPublic:  roomData.IsPublic,
	}

	chatroom := NewChatroom(roomInfo)

	handler.Hub().RegisterRoom(chatroom)
	go chatroom.Open()

	var (
		chatroomsTemplateID = hub_template.HubChatrooms.TemplateName
	)

	chatroomData := hub_template.PrepareChatroomData(rid, name)

	// one room
	templateData := hub_template.TemplateHubChatrooms{
		Rooms: []hub_template.ChatroomTemplateData{chatroomData},
	}

	return c.Render(http.StatusOK, chatroomsTemplateID, templateData)
}

// Redirection
func (handler *HubHandler) HandleUserJoinRequest(c echo.Context) error {

	rid := c.Param("roomID")

	var joinRequest hub_model.JoinRoomRequest
	err := c.Bind(&joinRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var userID string
	claims, jwtErr := GetJWTClaims(c)
	if jwtErr != nil {
		userID = claims.GetUID()
	} else {
		userID = fmt.Sprint(rand.IntN(1 << 8))
	}
	// TODO be wary of collisions of guest ids and user ids

	joinRequest.UserID = model.UID(userID)

	// TODO make sure user is invited if room is private

	echo.New().Logger.Printf("user join request received with data: %i", joinRequest)

	// set header for htmx to redirect from client-side
	chatroomRoute := fmt.Sprintf("/chatroom/%s", rid)
	c.Response().Header().Set("HX-Redirect", chatroomRoute)
	return c.NoContent(http.StatusFound)
}

func (handler *HubHandler) HandleChatroomPage(c echo.Context) error {
	// TODO handle unauthorized access to page

	userID, uidErr := GetJWTUserID(c)
	if uidErr != nil {
		c.Response().Header().Set("HX-Redirect", "/landing")
		return c.String(http.StatusNoContent, uidErr.Error())
	}

	roomID := c.Param("roomID")

	hubSvc, svcErr := createHubService()
	if svcErr != nil {
		return c.String(http.StatusInternalServerError, svcErr.Error())
	}

	getRoomReq := hub_model.GetChatroomRequest{
		UserID: model.UID(userID),
		RoomID: model.RID(roomID),
	}

	chatroom, dbErr := hubSvc.GetRoomInfo(c.Request().Context(), getRoomReq)
	if dbErr != nil {
		return c.String(http.StatusInternalServerError, dbErr.Error())
	}

	c.Logger().Printf("chatroomDTO: %+v", chatroom)

	chatroomInfo := chatroom_template.TemplateChatroomPage{
		RoomID:   string(chatroom.RoomID),
		RoomName: chatroom.RoomName,
	}

	chatroomID := chatroom_template.ChatroomPage.TemplateName

	return c.Render(http.StatusOK, chatroomID, chatroomInfo)

}

func (handler *HubHandler) HandleChatroomConnection(c echo.Context) error {
	// Websocket connection, should be
	// c.Logger().Print("HandleChatroomConnection()")

	tokenUserID, err := GetJWTUserID(c)
	if err != nil {
		return c.String(http.StatusFailedDependency, "Unable to retrieve JWT user id")
	}

	if !c.IsWebSocket() {
		errMsg := "expected Websocket connection, but was not"
		c.Logger().Print(errMsg)
		return c.NoContent(http.StatusUpgradeRequired)
	}

	var connReq hub_model.RoomRequest
	bindErr := c.Bind(&connReq)
	if bindErr != nil {
		errMsg := "connection request error"
		c.Logger().Print(errMsg)
		return c.NoContent(http.StatusBadRequest)
	}
	connReq.UserID = tokenUserID

	uid, rid := connReq.UserID, connReq.RoomID

	userReq := model.UserRequest{
		UserID: model.UserID(uid),
		RoomID: model.RoomID(rid),
	}


	chatroom := handler.Hub().Chatroom(userReq.RoomID)
	if chatroom == nil {
		chatroom_not_found := errors.New("could not find chatroom").Error()
		return c.String(http.StatusBadRequest, chatroom_not_found)
	}

	connErr := chatroom.AcceptConnection(c, userReq)
	if connErr != nil {
		// TODO handle err, tell client what error is maybe
		c.Logger().Print("connection error", connErr)
		return c.NoContent(http.StatusInternalServerError)
	}

	// Must return nil.
	// Do not write to context after websocket connection is successful.
	// Will get some "hijacked connection" warning.
	return nil

}

func (handler *HubHandler) HandleFetchChatroomHistory(c echo.Context) error {

	roomID := c.Param("roomID")
	chatroom := handler.Hub().Chatroom(model.RoomID(roomID))

	if chatroom == nil {
		// invalid request
		chatroom_not_found := errors.New("could not find chatroom").Error()
		c.Logger().Print(chatroom_not_found)
		return c.NoContent(http.StatusBadRequest)

	}

	chatroomLogs := chatroom.ChatLogs()

	chatLogsData := chatroom_template.TemplateManyMessages{
		ChatMessages: make([]chatroom_template.TemplateSingleMessage, 0),
	}

	for _, chatMessage := range chatroomLogs {

		log.Printf(`msg: "[%s]" by: [%s]`, chatMessage.Content, chatMessage.From)

		singleMsgData := chatroom_template.PrepareMessage(
			chatroom_template.WebsocketDivID,
			false,
			string(chatMessage.From),
			chatMessage.Content,
		)

		chatLogsData.ChatMessages = append(chatLogsData.ChatMessages, singleMsgData)
	}

	msgsTemplate := chatroom_template.ManyMessages.TemplateName
	return c.Render(http.StatusOK, msgsTemplate, chatLogsData)

}

func createHubService() (hub_svc.HubService, error) {

	db, err := db.ConnectDatabase()
	if err != nil {
		return nil, err
	}

	hubRepo := hub_svc.NewHubRepository(db.DB())
	hubSvc := hub_svc.NewHubService(hubRepo)

	return hubSvc, nil

}

func (handler *HubHandler) openHubChatrooms() {

	hub := handler.Hub()

	pubChatrooms := hub.PublicChatrooms()
	privChatrooms := hub.PrivateChatrooms()

	go func() {
		for _, chatroom := range pubChatrooms {
			go chatroom.Open()
		}
	}()

	go func() {
		for _, chatroom := range privChatrooms {
			go chatroom.Open()
		}
	}()

}


func getChatroomsDB() ([]chatroom_model.Chatroom, error) {
	hubSvc, svcErr := createHubService()
	if svcErr != nil {
		return nil, svcErr
	}

	req := hub_model.GetPublicChatroomsRequest{
		IsPublic: true,
		IsActive: true,
	}

	chatroomsRes, resErr := hubSvc.GetRoomsPublic(context.Background(), req)
	if resErr != nil {
		return nil, resErr
	}

	hubChatrooms := []chatroom_model.Chatroom{}
	for _, chatroom := range chatroomsRes {
		chatroomInfo := model.ChatroomInfo{
			RoomID:    chatroom.RoomID,
			RoomName:  chatroom.RoomName,
			RoomOwner: chatroom.OwnerID,
			IsPublic:  chatroom.IsPublic,
		}

		newChatroom := NewChatroom(chatroomInfo)
		hubChatrooms = append(hubChatrooms, newChatroom)

	}


	return hubChatrooms, nil
}