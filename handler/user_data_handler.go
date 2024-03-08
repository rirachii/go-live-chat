package handler

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func HandleGetUsername(c echo.Context) error {

	c.Logger().Print("get username")

	const (
		templateID    = "user-username"
		usernameField = "Username"
	)

	// assume failure
	usernameData := map[string]string{
		usernameField: "Guest",
	}

	jwt, jwtError := getJWTCookie(c)
	if jwtError != nil {
		// do nothing
	} else {
		username := jwt.GetUsername()
		usernameData[usernameField] = username
	}

	return c.Render(http.StatusOK, templateID, usernameData)
}
