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

	SetupEchoServer(e)

	_, hubHandler := handler.InitiateHub()
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

	t := service.NewTemplateRenderer("templates/*/*.html")
	e.Renderer = t

	// MIDDLE WARE, TODO middleware for JWT
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "static",
		Browse: false,
	}))

	// TODO middleware for JWT.
	// secretKey := "TODO_change_to_something_better_secret"
	// e.Use(echojwt.WithConfig(echojwt.Config{

	// 	SuccessHandler: func(c echo.Context) {
	// 		c.Logger().Print("Token found!")
	// 	},
	// 	ErrorHandler: func(c echo.Context, err error) error {
	// 		c.Logger().Print("no token found", c.Cookies())
	// 		return nil
	// 	},
	// 	ContextKey:  "jwt",
	// 	SigningKey:  []byte(secretKey),
	// 	TokenLookup: "cookie:jwt",
	// }))

	// e.Use(middleware.JWT())

}

func InitializeRoutes(e *echo.Echo, hubHandler *handler.HubHandler) {
	InitializeAPIRoutes(e)
	InitializeUserAuthRoutes(e)
	InitializeUserDataRoutes(e)
	InitializeHubRoutes(e, hubHandler)
}

func InitializeHubRoutes(e *echo.Echo, hubHandler *handler.HubHandler) {

	hub := e.Group("/hub")

	hub.GET("*", handler.HandleHubPage)
	hub.GET("/get-rooms", hubHandler.HandleGetChatrooms)
	hub.POST("/create-room", hubHandler.HandleCreateRoom)
	hub.POST("/join/:roomID", hubHandler.HandleUserJoinRequest)
	hub.GET("/chatroom/:roomID", hubHandler.HandleChatroomPage)
	hub.GET("/chatroom/:roomID/chat-history", hubHandler.HandleFetchChatroomHistory)
	hub.GET("/chatroom/:roomID/ws", hubHandler.HandleChatroomConnection)
	hub.GET("/ws/:roomID", handler.HandleGetChatroomWebsocket)

	//TODO: instead we should run when user is logged in securly,
	// maybe we can allow guests to join rooms but not create rooms
	go hubHandler.Hub.Run()

}

func InitializeUserAuthRoutes(e *echo.Echo) {
	e.GET("/register", handler.HandleRegisterPage)
	e.POST("/register", handler.HandleUserRegister)

	e.GET("/login", handler.HandleLoginPage)
	e.POST("/login", handler.HandleUserLogin)
	e.GET("/logout", handler.HandleUserLogout)
}

func InitializeUserDataRoutes(e *echo.Echo) {

	e.GET("/user/username", handler.HandleGetUsername)
	e.GET("/user/profile-pic", handler.HandleGetUserProfile)

	e.GET("/user/user-rooms", handler.HandleGetUserRooms)
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
