package handlers

import (
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/rirachii/golivechat/db"
	"github.com/rirachii/golivechat/users"
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

func getUserHandler() (*users.Handler, error) {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize postgres db connection: %s", err)
	}

	userRep := users.NewRepository(dbConn.GetDB())
	userSvc := users.NewService(userRep)
	userHandler := users.NewHandler(userSvc)
	return userHandler, nil
}

func HandleCreateUser(c echo.Context) error {
	userHandler, err := getUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}
	userHandler.Login(c)

	return c.Redirect(http.StatusFound, "/landing")
}

func HandleLogin(c echo.Context) error {
	userHandler, err := getUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}
	userHandler.Login(c)
	return c.Redirect(http.StatusFound, "/landing")
}

func HandleLogout(c echo.Context) error {
	userHandler, err := getUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}
	userHandler.Logout(c)
	return c.Redirect(http.StatusFound, "/landing")
}
