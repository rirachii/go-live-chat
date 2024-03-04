package handler

import (
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

func HandleLoginPageDisplay(c echo.Context) error {
	loginTemplate := "login"

	return c.Render(http.StatusOK, loginTemplate, nil)
}



func HandleHubPage(c echo.Context) error {
	hubTemplate := "hub"

	return c.Render(http.StatusOK, hubTemplate, nil)
}

