package handler

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/rirachii/golivechat/service"

	hub_template "github.com/rirachii/golivechat/templates/hub"
	landing_template "github.com/rirachii/golivechat/templates/landing"
	login_template "github.com/rirachii/golivechat/templates/login"
	register_template "github.com/rirachii/golivechat/templates/register"
)

func HandleLanding(c echo.Context) error {

	// jwt, ok := c.Get("jwt").(*jwt.Token)
	// c.Logger().Print(c.Cookies())

	data := landing_template.TemplateLandingPage{
		Title: "LIVE CHAT SERVERRR!",
	}

	landingTemplate := landing_template.LandingPage.TemplateName
	return c.Render(http.StatusOK, landingTemplate, data)
}

func HandleRegisterPage(c echo.Context) error {

	_, err := getJWTCookie(c)
	if err == nil {
		return c.Redirect(http.StatusSeeOther, "/hub")
	}

	registerTemplate := register_template.RegisterPage.TemplateName
	return c.Render(http.StatusOK, registerTemplate, nil)
}

func HandleLoginPage(c echo.Context) error {
	_, err := getJWTCookie(c)
	if err == nil {
		return c.Redirect(http.StatusSeeOther, "/hub")
	}

	loginTemplate := login_template.LoginPage.TemplateName
	return c.Render(http.StatusOK, loginTemplate, nil)
}

func HandleHubPage(c echo.Context) error {
	_, err := getJWTCookie(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	hubTemplate := hub_template.HubPage.TemplateName
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
