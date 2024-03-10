package handler

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/rirachii/golivechat/service"
)

func HandleLanding(c echo.Context) error {
	landingTemplate := "landing"

	// jwt, ok := c.Get("jwt").(*jwt.Token)
	// c.Logger().Print(c.Cookies())

	data := make(map[string]string)
	data["Title"] = "LIVE CHAT SERVERRR!"

	return c.Render(http.StatusOK, landingTemplate, data)
}

func HandleRegisterPage(c echo.Context) error {
	template := "register"
	_, err := getJWTCookie(c)
	if err == nil {
		return c.Redirect(http.StatusSeeOther, "/hub")
	}

	return c.Render(http.StatusOK, template, nil)
}

func HandleLoginPage(c echo.Context) error {
	template := "login"
	_, err := getJWTCookie(c)
	if err == nil {
		return c.Redirect(http.StatusSeeOther, "/hub")
	}

	return c.Render(http.StatusOK, template, nil)
}

func HandleHubPage(c echo.Context) error {
	_, err := getJWTCookie(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	hubTemplate := "hub"
	return c.Render(http.StatusOK, hubTemplate, nil)

}

func getJWTCookie(c echo.Context) (*service.MyJWTClaims, error) {

	cookie, err := c.Cookie("jwt")
	if err != nil {
		return nil, err
	}

	tokenString := cookie.Value
	validTokenClaims, validateErr := service.ValidateJWT(tokenString)
	if validateErr != nil {
		echo.New().Logger.Print(validateErr)
		return nil, validateErr
	}

	return validTokenClaims, nil
}
