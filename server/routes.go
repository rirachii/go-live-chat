package main

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	api "github.com/rirachii/golivechat/server/api"
	chat "github.com/rirachii/golivechat/server/chat"
	handlers "github.com/rirachii/golivechat/server/handlers"
)

func InitializeRoutes(e *echo.Echo, hubHandler *chat.HubHandler) {

	e.GET("*", redirectToLanding)
	e.GET("/landing", handlers.HandleLanding)

	InitializeAPIRoutes(e)
	InitializeUserRoutes(e)
	InitializeHubRoutes(e, hubHandler)

}

func InitializeHubRoutes(e *echo.Echo, hubHandler *chat.HubHandler) {
	e.GET("/hub*", handlers.HandleHubPage)
	e.GET("/hub/get-rooms", hubHandler.HandleGetChatrooms)
	e.POST("/hub/create-room", hubHandler.HandleCreateRoom)
	e.POST("/hub/join/:roomID", hubHandler.HandleUserJoinRequest)
	e.GET("/hub/chatroom/:roomID", hubHandler.HandleChatroomPage)
	e.GET("/hub/chatroom/:roomID/chat-history", hubHandler.HandleFetchChatroomHistory)
	e.GET("/hub/chatroom/:roomID/ws", hubHandler.HandleChatroomWSConnection)

	e.GET("/ws/:roomID", chat.HandleGetChatroomWebsocket)

}

func InitializeUserRoutes(e *echo.Echo){
	
	e.GET("/register", handlers.HandleRegisterPageDisplay)
	e.POST("/register", handlers.HandleRegisterUser)


	
	// e.GET("/login", func(c echo.Context) error { 
	// 	return c.Redirect(http.StatusFound,"/hub") 
	// })

	e.POST("/signup", handlers.HandleCreateUser)
	e.POST("/login", handlers.HandleLogin)
	e.GET("/logout", handlers.HandleLogout)

}

func InitializeAPIRoutes(e *echo.Echo) {
	e.GET("/random-msgs", getRandomMsg)

}


func getRandomMsg(c echo.Context) error {

	randomMsg := api.RandomMsg()

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, randomMsg)

}


func redirectToLanding(c echo.Context) error {

	return c.Redirect(http.StatusPermanentRedirect, "/landing")

}

