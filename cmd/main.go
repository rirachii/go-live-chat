package main

import (
	"log"
	"net/http"
	"os"

	godotenv "github.com/joho/godotenv"
	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	handler "github.com/rirachii/golivechat/handler"
	templates "github.com/rirachii/golivechat/templates"
	internal "github.com/rirachii/golivechat/internal"

)

func main() {

	// ensure .env exists
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal("Could not load .env file")
	}

	e := echo.New()

	SetupEchoServer(e)

	hubHandler, err := handler.InitiateHubHandler()
	if err != nil {
		log.Print(err.Error())
		return
	}

	InitializeRoutes(e, hubHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	e.Logger.Fatal(e.Start(port))
}

func SetupEchoServer(e *echo.Echo) {

	// STATIC Routes
	e.File("/favicon.ico", "static/public/images/favicon.ico")
	e.GET("/", redirectToLanding)
	e.GET("/landing", handler.HandleLanding)

	t := templates.NewTemplateRenderer("templates/*/*.html")
	e.Renderer = t

	// MIDDLE WARE, TODO middleware for JWT
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "static",
		Browse: false,
	}))

}

func InitializeRoutes(e *echo.Echo, hubHandler *handler.HubHandler) {
	InitializeAPIRoutes(e)
	InitializeUserAuthRoutes(e)
	InitializeUserDataRoutes(e)
	InitializeHubRoutes(e, hubHandler)
	InitializeChatroomRoutes(e, hubHandler)
}

func InitializeHubRoutes(e *echo.Echo, hubHandler *handler.HubHandler) {

	hub := e.Group("/hub")

	hub.GET("*", handler.HandleHubPage)
	hub.POST("/create-room", hubHandler.HandleCreateRoom)
	hub.GET("/get-public-rooms", hubHandler.HandleGetPublicChatrooms)
	hub.POST("/join-room/:roomID", hubHandler.HandleUserJoinRequest)

	//TODO: instead we should run when user is logged in securly,
	// maybe we can allow guests to join rooms but not create rooms

}

func InitializeChatroomRoutes(e *echo.Echo, hubHandler *handler.HubHandler) {

	chatroom := e.Group("/chatroom")

	chatroom.GET("/:roomID", hubHandler.HandleChatroomPage)
	chatroom.GET("/:roomID/chat-history", hubHandler.HandleFetchChatroomHistory)

	// websocket connection
	chatroom.GET("/:roomID/ws", hubHandler.HandleChatroomConnection)

	// gets html for websocket
	chatroom.GET("/:roomID/get-ws", handler.HandleGetChatroomWebsocket)
}

func InitializeUserAuthRoutes(e *echo.Echo) {
	// e.GET("/register", handler.HandleRegisterPage)
	e.POST("/register", handler.HandleUserRegister)

	// e.GET("/login", handler.HandleLoginPage)
	e.POST("/login", handler.HandleUserLogin)
	e.GET("/logout", handler.HandleUserLogout)
}

func InitializeUserDataRoutes(e *echo.Echo) {

	user := e.Group("/user")

	user.GET("/username", handler.HandleGetUsername)
	user.GET("/profile-pic", handler.HandleGetUserProfile)
	user.GET("/user-rooms", handler.HandleGetUserRooms)
}

func InitializeAPIRoutes(e *echo.Echo) {
	e.GET("/random-msgs", getRandomMsg)
	e.POST("/locate-user", handler.HandleLocateUser)
}

func getRandomMsg(c echo.Context) error {
	randomMsg := internal.RandomMsg()

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, randomMsg)

}

func redirectToLanding(c echo.Context) error {

	return c.Redirect(http.StatusPermanentRedirect, "/landing")

}
