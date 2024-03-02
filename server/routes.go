package main

import (
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
	api "github.com/rirachii/golivechat/server/api"
	chat "github.com/rirachii/golivechat/server/chat"
	handlers "github.com/rirachii/golivechat/server/handlers"
)

func InitializeRoutes(e *echo.Echo, hubHandler *chat.HubHandler) {

	e.GET("/landing", handlers.HandleLanding)

	e.GET("/register", handlers.HandleRegisterPageDisplay)
	e.POST("/register", handlers.HandleRegisterUser)


	
	e.GET("/login", func(c echo.Context) error { 
		return c.Redirect(http.StatusFound,"/hub") 
	})

	e.POST("/signup", handlers.HandleCreateUser)
	e.POST("/login", handlers.HandleLogin)
	e.GET("/logout", handlers.HandleLogout)

	// e.POST("/login", func(c echo.Context) error { return echo.ErrNotImplemented })
	InitializeHubRoutes(e, hubHandler)
	InitializeAPIRoutes(e)

}

func InitializeHubRoutes(e *echo.Echo, hubHandler *chat.HubHandler) {
	e.GET("/hub*", handlers.HandleHubPage)
	e.GET("/hub/get-rooms", hubHandler.HandleGetChatrooms)
	e.POST("/hub/create-room", hubHandler.HandleCreateRoom)
	e.POST("/hub/join/:roomID", hubHandler.HandleUserJoinRequest)
	e.GET("/hub/chatroom/:roomID", hubHandler.HandleChatroomPage)
	e.GET("/hub/chatroom/:roomID/chat-history", hubHandler.HandleFetchChatroomHistory)
	e.GET("/ws/:roomID", ServeChatroomConnection)
	e.GET("/hub/chatroom/:roomID/ws", hubHandler.HandleChatroomWSConnection)

}


func InitializeAPIRoutes(e *echo.Echo) {
	e.GET("/random-msgs", getRandomMsg)

}

func getRandomMsg(c echo.Context) error {

	randomMsg := api.RandomMsg()

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, randomMsg)

}


func ServeChatroomConnection(c echo.Context) error {

	roomID := c.Param("roomID")
	echo.New().Logger.Printf(c.QueryString())


	roomData := map[string]string{
		"ConnectionRoute": fmt.Sprintf("/hub/chatroom/%s/ws", roomID),
		"RoomID": roomID,
	}

	const chatroomConnectionTemplateID = "chatroom-connection"
	return c.Render(http.StatusOK, chatroomConnectionTemplateID, roomData)

}
