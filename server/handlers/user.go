package handlers

import (
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/rirachii/golivechat/db"
	"github.com/rirachii/golivechat/users"
)


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