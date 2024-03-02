package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

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

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	address := fmt.Sprintf("%s:%s", Host, Port)

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/pages/*.html")),
	}

	e := echo.New()

	e.Renderer = t
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   clientFolder,
		Browse: false,
	}))


	e.File("/favicon.ico", clientFolder+"/public/images/favicon.ico")
	e.GET("/", redirectToLanding)

	hub, hubHandler := chat.InitiateHub()
	go hub.Run()
	
	InitializeRoutes(e, hubHandler)

	// Open server
	log.Println("Listening on:", fmt.Sprintf("http://%s", address))
	err := e.Start(":8080")
	if err != nil {
		e.Logger.Fatal("Error Starting the HTTP Server : ", err)
		return
	}

}

// func handle (w,r)

func redirectToLanding(c echo.Context) error {

	return c.Redirect(http.StatusPermanentRedirect, "/landing")

}


