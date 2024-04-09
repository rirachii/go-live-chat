package handler

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	chatroom_template "github.com/rirachii/golivechat/templates/chatroom"
)

func HandleGetChatroomWebsocket(c echo.Context) error {

	_, jwtErr := GetJWTClaims(c)
	if jwtErr != nil {
		return c.String(echo.ErrBadRequest.Code, jwtErr.Error())
	}

	roomID := c.Param("roomID")

	connectionData := chatroom_template.TemplateChatroomConnection{
		RoomID: roomID,
	}
	
	c.Logger().Printf("%+v", connectionData)

	connectionTemplate := chatroom_template.ChatroomConnection.TemplateName
	return c.Render(http.StatusOK, connectionTemplate, connectionData)

}
