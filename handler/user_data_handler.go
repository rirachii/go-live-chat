package handler

import (
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
)

func HandleGetUsername(c echo.Context) error {

	const (
		templateID    = "user-username"
		usernameField = "Username"
	)

	// assume failure
	usernameData := map[string]string{
		usernameField: "Guest",
	}

	jwtClaims, jwtError := GetJWTClaims(c)
	if jwtError != nil {
		// do nothing
	} else {
		username := jwtClaims.GetUsername()
		usernameData[usernameField] = username
	}

	return c.Render(http.StatusOK, templateID, usernameData)
}

type userImg struct {
	ImgExists bool
	ImgLink   string
}

func HandleGetUserProfile(c echo.Context) error {
	const (
		templateID = "user-profile"
		imgDataID  = "UserImg"
	)

	// test loading
	time.Sleep(time.Second * 0)

	imgData := userImg{
		ImgExists: false,
	}

	templateData := map[string]userImg{
		imgDataID: imgData,
	}

	return c.Render(http.StatusOK, templateID, templateData)

}

func HandleGetUserRooms(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}
