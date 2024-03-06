package service

import (
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
)


func HandleGetChatroomWebsocket(c echo.Context) error {


	// todo use accesstoken
	var (
		userID = c.QueryParam("userID")
		roomID = c.Param("roomID")

	)


	roomData := map[string]string{
		// TODO change userID to be token
		"ConnectionRoute": fmt.Sprintf("/hub/chatroom/%s/ws?userID=%s", roomID, userID),
		"RoomID": roomID,
	}

	const chatroomConnectionTemplateID = "chatroom-connection"
	return c.Render(http.StatusOK, chatroomConnectionTemplateID, roomData)

}