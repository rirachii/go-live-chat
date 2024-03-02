package chat

import (
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
)



func InitiateHub() (*ChatroomsHub, *HubHandler) {

	hub := &ChatroomsHub{
		ChatRooms:       make(map[RoomID]*Chatroom),
		UserChatrooms:   make(map[UserID]SetOfChatrooms),
		RegisterQueue:   make(chan *User),
		UnregisterQueue: make(chan *User),
	}

	handler := &HubHandler{
		Hub: hub,
	}

	return hub, handler
}



func HandleGetChatroomWebsocket(c echo.Context) error {

	roomID := c.Param("roomID")
	echo.New().Logger.Printf(c.QueryString())


	roomData := map[string]string{
		"ConnectionRoute": fmt.Sprintf("/hub/chatroom/%s/ws", roomID),
		"RoomID": roomID,
	}

	const chatroomConnectionTemplateID = "chatroom-connection"
	return c.Render(http.StatusOK, chatroomConnectionTemplateID, roomData)

}