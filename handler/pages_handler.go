package handler

import (
	"net/http"

	echo "github.com/labstack/echo/v4"

	hub_template "github.com/rirachii/golivechat/templates/hub"
	landing_template "github.com/rirachii/golivechat/templates/landing"
	login_template "github.com/rirachii/golivechat/templates/login"
	register_template "github.com/rirachii/golivechat/templates/register"
)

func HandleLanding(c echo.Context) error {

	// jwt, ok := c.Get("jwt").(*jwt.Token)
	// c.Logger().Print(c.Cookies())

	_, err := GetJWTCookie(c)
	if err == nil {
		return c.Redirect(http.StatusSeeOther, "/hub")
	}

	data := landing_template.TemplateLandingPage{
		Title: "Welcome to GoChat!",
	}

	landingTemplate := landing_template.LandingPage.TemplateName
	return c.Render(http.StatusOK, landingTemplate, data)
}

func HandleRegisterPage(c echo.Context) error {
	

	_, err := GetJWTCookie(c)
	if err == nil {
		return c.Redirect(http.StatusSeeOther, "/hub")
	}

	data := register_template.TemplateRegisterPage{
		Title: "Welcome to GoChat!",
	}


	registerTemplate := register_template.RegisterPage.TemplateName
	return c.Render(http.StatusOK, registerTemplate, data)
}

func HandleLoginPage(c echo.Context) error {
	_, err := GetJWTCookie(c)
	if err == nil {
		return c.Redirect(http.StatusSeeOther, "/hub")
	}

	loginTemplate := login_template.LoginPage.TemplateName
	return c.Render(http.StatusOK, loginTemplate, nil)
}

func HandleHubPage(c echo.Context) error {
	_, err := GetJWTCookie(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/landing")
	}

	hubTemplate := hub_template.HubPage.TemplateName
	return c.Render(http.StatusOK, hubTemplate, nil)

}