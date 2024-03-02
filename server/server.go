package main

import (
	"fmt"
	"log"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/rirachii/golivechat/server/chat"
	// chat "github.com/rirachii/golivechat/server/chat"
)

const (
	// Host name of the HTTP Server
	Host = "localhost"
	// Port of the HTTP Server
	Port = "8080"
	// File folder for Frontend
	clientFolder = "client"
)

func main() {
	address := fmt.Sprintf("%s:%s", Host, Port)

	e := echo.New()
	SetupEcho(e)

	// Open server
	log.Println("Listening on:", fmt.Sprintf("http://%s", address))
	err := e.Start(":8080")
	if err != nil {
		e.Logger.Fatal("Error Starting the HTTP Server : ", err)
		return
	}

}

func SetupEcho(e *echo.Echo){

	t := NewTemplateRenderer("templates/pages/*.html")

	e.Renderer = t
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   clientFolder,
		Browse: false,
	}))



	e.File("/favicon.ico", clientFolder+"/public/images/favicon.ico")
	e.GET("/", redirectToLanding)
	

	hub, hubHandler := chat.InitiateHub()
	
	InitializeRoutes(e, hubHandler)
	go hub.Run()


}



