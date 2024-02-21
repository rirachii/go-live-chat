package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
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

	// e.Static("/js", "js")
	// e.Static("/css", "css")
	e.File("/favicon.ico", "client/public/images/favicon.ico")

	e.GET("/", handleRootPageRequest)
	e.GET("/landing", handleLanding)
	e.GET("/register", handleRegister)
	e.POST("/register", handleRegister)

	e.GET("/random-msgs", getRandomMsgs)

	// Open server
	log.Println("Listening on:", fmt.Sprintf("http://%s", address))
	err := e.Start(":8080")
	if err != nil {
		e.Logger.Fatal("Error Starting the HTTP Server : ", err)
		return
	}

}

func handleRootPageRequest(c echo.Context) error {


	return c.Redirect(http.StatusPermanentRedirect,"/landing")

}

func handleLanding(c echo.Context) error {

	landingTemplate := "landing"

	data := make(map[string]string)
	data["Title"] = "LIVE CHAT SERVERRR!"

	return c.Render(http.StatusOK, landingTemplate, data)

}

func handleRegister(c echo.Context) error {

	log.Println(c.Request().Method, c.Request().URL.Path)

	requestMethod := c.Request().Method

	switch requestMethod {
	case "GET":
		// Serve html
		registerTemplate := "register"

		return c.Render(http.StatusOK, registerTemplate, nil)

	case "POST":

		// Handle request to register
		log.Println("Register POST data received!")

		type RegisterForm struct {
			Username string
			Password string
		}

		username := c.FormValue("username")
		password := c.FormValue("password")

		postData := RegisterForm{
			Username: username,
			Password: password,
		}

		log.Println("Received Username: ", postData.Username)
		log.Println("Received Password: ", postData.Password)

		c.Response().Header().Set("HX-Location", "landing")
		c.Response().WriteHeader(http.StatusFound)

		return c.Redirect(http.StatusFound, "/landing")

	default:

		// c.Response().WriteHeader(http.StatusMethodNotAllowed)
		return c.NoContent(http.StatusMethodNotAllowed)
	}
}

func getRandomMsgs(c echo.Context) error {

	msgs := []string{
		"random 1",
		"welcome to the unknown",
		"im a random messsage",
		"KKB on toppp",
		"akjsdhiuandi",
	}

	randomIndex := rand.Intn(len(msgs))
	randomMsg := msgs[randomIndex]

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, randomMsg)

}
