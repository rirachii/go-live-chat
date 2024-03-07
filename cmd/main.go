package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rirachii/golivechat/handler"
	"github.com/rirachii/golivechat/service"
)

func main() {
	e := echo.New()

	t := service.NewTemplateRenderer("templates/pages/*.html")
	e.Renderer = t
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "static",
		Browse: false,
	}))

	e.File("/favicon.ico", "static/public/images/favicon.ico")
	e.GET("/", redirectToLanding)
	e.GET("/landing", handler.HandleLanding)

	hub, hubHandler := handler.InitiateHub()
	InitializeRoutes(e, hubHandler)

	//TODO: instead we should run when user is logged in securly,
	go hub.Run()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	e.Logger.Fatal(e.Start(port))
}


func InitializeRoutes(e *echo.Echo, hubHandler *handler.HubHandler) {
	InitializeAPIRoutes(e)
	InitializeUserAuthRoutes(e)
	InitializeHubRoutes(e, hubHandler)
}

func InitializeHubRoutes(e *echo.Echo, hubHandler *handler.HubHandler) {
	e.GET("/hub*", handler.HandleHubPage)
	e.GET("/hub/get-rooms", hubHandler.HandleGetChatrooms)
	e.POST("/hub/create-room", hubHandler.HandleCreateRoom)
	e.POST("/hub/join/:roomID", hubHandler.HandleUserJoinRequest)
	e.GET("/hub/chatroom/:roomID", hubHandler.HandleChatroomPage)
	e.GET("/hub/chatroom/:roomID/chat-history", hubHandler.HandleFetchChatroomHistory)
	e.GET("/hub/chatroom/:roomID/ws", hubHandler.HandleChatroomConnection)
	e.GET("/ws/:roomID", service.HandleGetChatroomWebsocket)

}

func InitializeUserAuthRoutes(e *echo.Echo) {
	e.GET("/register", handler.HandleRegisterPage)
	e.POST("/register", handler.HandleUserRegister)

	e.GET("/login", handler.HandleLoginPage)
	e.POST("/login", handler.HandleUserLogin)
	e.GET("/logout", handler.HandleUserLogout)
}


func InitializeAPIRoutes(e *echo.Echo) {
	e.GET("/random-msgs", getRandomMsg)
}

func getRandomMsg(c echo.Context) error {
	randomMsg := service.RandomMsg()

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, randomMsg)

}

func redirectToLanding(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/landing")

}
