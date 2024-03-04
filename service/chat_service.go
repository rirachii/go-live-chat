package service

import (
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
)


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