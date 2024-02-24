package handlers

import (
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func HandleLanding(c echo.Context) error {

	landingTemplate := "landing"

	data := make(map[string]string)
	data["Title"] = "LIVE CHAT SERVERRR!"

	return c.Render(http.StatusOK, landingTemplate, data)

}


func HandleRegisterPageDisplay(c echo.Context) error {
	registerTemplate := "register"

	return c.Render(http.StatusOK, registerTemplate, nil)
}

func HandleRegisterUser(c echo.Context) error {


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

}
