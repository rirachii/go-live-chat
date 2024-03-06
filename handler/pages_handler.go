package handler

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/rirachii/golivechat/service"
)

func HandleLanding(c echo.Context) error {
	landingTemplate := "landing"

	data := make(map[string]string)
	data["Title"] = "LIVE CHAT SERVERRR!"
	
	return c.Render(http.StatusOK, landingTemplate, data)
}


func HandleRegisterPageDisplay(c echo.Context) error {
	template := "register"
	err := checkCookie(c)
	if err == nil {
		template = "hub"
	}

	return c.Render(http.StatusOK, template, nil)
}

func HandleLoginPageDisplay(c echo.Context) error {
	template := "login"
	err := checkCookie(c)
	if err == nil {
		template = "hub"
	}

	return c.Render(http.StatusOK, template, nil)
}

func HandleHubPage(c echo.Context) error {
	err := checkCookie(c)
	if err != nil {
		return err
	}

	hubTemplate := "hub"
	return c.Render(http.StatusOK, hubTemplate, nil)
}


func checkCookie(c echo.Context) error {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return err
	}

	tokenString := cookie.Value
	err = service.ValidateJWT(tokenString) 
	if err != nil {
		echo.New().Logger.Print(err)
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}
	return nil
}
