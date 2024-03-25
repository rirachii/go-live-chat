package handler

import (
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
	model "github.com/rirachii/golivechat/model"
	chat_template "github.com/rirachii/golivechat/templates/chat"
)

func HandleGetChatroomWebsocket(c echo.Context) error {

	// todo use accesstoken
	var (
		userID = c.QueryParam("userID")
		roomID = c.Param("roomID")
	)

	route := fmt.Sprintf("/hub/chatroom/%s/ws?userID=%s", roomID, userID)

	connectionData := chat_template.TemplateChatroomConnection{
		ConnectionRoute: route,
		RoomID:          model.RID(roomID),
	}

	connectionTemplate := chat_template.ChatroomConnection.TemplateName
	return c.Render(http.StatusOK, connectionTemplate, connectionData)

}
