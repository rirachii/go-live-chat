package main

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	api "github.com/rirachii/golivechat/server/api"
	handler "github.com/rirachii/golivechat/server/handlers"
)

func InitializeRoutes(e *echo.Echo) {

	e.GET("/landing", handler.HandleLanding)

	e.GET("/register", handler.HandleRegisterPageDisplay)
	e.POST("/register", handler.HandleRegisterUser)


	
	e.POST("/login", func(c echo.Context) error { return echo.ErrNotImplemented })


	// e.GET("/hub", handleHubPage)
	// e.POST("/hub/createRoom") 
	// e.GET("/chat/:roomId", func(c echo.Context) error { return c.String(http.StatusNotImplemented, c.Request().URL.Path) })
	// e.POST("/chat/:roomId/sendChat")

	e.GET("/random-msgs", getRandomMsg)

}

func getRandomMsg(c echo.Context) error {

	randomMsg := api.RandomMsg()

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, randomMsg)

}
