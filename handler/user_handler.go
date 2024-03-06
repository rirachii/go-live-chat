package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rirachii/golivechat/model"
	"github.com/rirachii/golivechat/service"
	"github.com/rirachii/golivechat/service/db"
)

type UserHandler struct {
	UserService service.UserService
}

func NewHandler(s service.UserService) *UserHandler {
	return &UserHandler{
		UserService: s,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	email := c.FormValue("email")
	username := c.FormValue("username")
	password := c.FormValue("password")
	u := model.CreateUserReq{Username: username, Email: email, Password: password}

	if email == "" || username == "" || password == "" {
		c.JSONBlob(http.StatusBadRequest, []byte("a field is empty"))
		return errors.New("email and password field incorrect")
	}

	_, err := h.UserService.CreateUser(c.Request().Context(), &u)
	if err != nil {
		c.JSONBlob(http.StatusInternalServerError, []byte("create user error: "+err.Error()))
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}

func (h *UserHandler) Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	user := model.LoginUserReq{Email: email, Password: password}

	if email == "" || password == "" {
		c.JSONBlob(http.StatusBadRequest, []byte("a field is empty"))
		return errors.New("email and password field incorrect")
	}

	u, err := h.UserService.Login(c.Request().Context(), &user)
	if err != nil {
		c.JSONBlob(http.StatusInternalServerError, []byte("login user error: "+err.Error()))
		return err
	}

	c.SetCookie(&http.Cookie{Name: "jwt", Value: u.GetAccessToken(), MaxAge: 3600, Domain: "localhost", Secure: false, HttpOnly: true})
	return c.Redirect(http.StatusSeeOther, "/hub")
	
}

func (h *UserHandler) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{Name: "jwt", Value: "", MaxAge: -1, Domain: "localhost", Secure: false, HttpOnly: true})
	return c.JSON(http.StatusOK, "Logout successful")
}



// USER ROUTES HANDLER
func getUserHandler() (*UserHandler, error) {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize postgres db connection: %s", err)
	}

	userRep := service.NewRepository(dbConn.GetDB())
	userSvc := service.NewService(userRep)
	userHandler := NewHandler(userSvc)
	return userHandler, nil
}

func HandleCreateUser(c echo.Context) error {
	userHandler, err := getUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}

	return userHandler.CreateUser(c)
}

func HandleLogin(c echo.Context) error {
	userHandler, err := getUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}
	return userHandler.Login(c)
}

func HandleLogout(c echo.Context) error {
	userHandler, err := getUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}
	
	return userHandler.Logout(c)
}
