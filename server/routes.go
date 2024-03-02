package main

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	api "github.com/rirachii/golivechat/server/api"
	chat "github.com/rirachii/golivechat/server/chat"
	handler "github.com/rirachii/golivechat/server/handlers"
)


func InitializeRoutes(e *echo.Echo, hubHandler *chat.HubHandler) {

	e.GET("/landing", handler.HandleLanding)

	e.GET("/register", handler.HandleRegisterPageDisplay)
	e.POST("/register", handler.HandleRegisterUser)


	
	e.GET("/login", func(c echo.Context) error { 
		return c.Redirect(http.StatusFound,"/hub") 
	})

	e.POST("/signup", handler.HandleCreateUser)
	e.POST("/login", handler.HandleLogin)
	e.GET("/logout", handler.HandleLogout)

	// e.POST("/login", func(c echo.Context) error { return echo.ErrNotImplemented })

	e.GET("/hub", handler.HandleHubPage)
	e.GET("/hub/get-rooms", hubHandler.GetHubRooms)
	e.POST("/hub/create-room", hubHandler.CreateRoom)
	// e.GET("/chat/:roomId", func(c echo.Context) error { return c.String(http.StatusNotImplemented, c.Request().URL.Path) })
	// e.POST("/chat/:roomId/send-chat")

	e.GET("/random-msgs", getRandomMsg)

}

func getRandomMsg(c echo.Context) error {

	randomMsg := api.RandomMsg()

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, randomMsg)

}
