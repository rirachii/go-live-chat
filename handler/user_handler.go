package handler

import (
	"fmt"
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
	model "github.com/rirachii/golivechat/model"
	service "github.com/rirachii/golivechat/service"
	db "github.com/rirachii/golivechat/service/db"
)

type Handler struct {
	UserService service.UserService
}

func NewHandler(s service.UserService) *Handler {
	return &Handler{
		UserService: s,
	}
}

func (h *Handler) CreateUser(c echo.Context) {
	email := c.FormValue("email")
	username := c.FormValue("username")
	password := c.FormValue("password")
	u := model.CreateUserReq{Username: username, Email: email, Password: password}

	if email == "" || username == "" || password == "" {
		c.JSONBlob(http.StatusBadRequest, []byte("a field is empty"))
		return
	}

	res, err := h.UserService.CreateUser(c.Request().Context(), &u)
	if err != nil {
		c.JSONBlob(http.StatusInternalServerError, []byte("create user error: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, res)
	c.Redirect(http.StatusFound, "/login")

}

func (h *Handler) Login(c echo.Context) {
	email := c.FormValue("email")
	password := c.FormValue("password")
	user := model.LoginUserReq{Email: email, Password: password}

	if email == "" || password == "" {
		c.JSONBlob(http.StatusBadRequest, []byte("a field is empty"))
		return
	}

	u, err := h.UserService.Login(c.Request().Context(), &user)
	if err != nil {
		c.JSONBlob(http.StatusInternalServerError, []byte("login user error: "+err.Error()))
		return
	}

	c.SetCookie(&http.Cookie{Name: "jwt", Value: u.GetAccessToken(), MaxAge: 3600, Path: "/landing", Domain: "localhost", Secure: false, HttpOnly: true})

	// res := &model.LoginUserRes{
	// 	Username: u.Username,
	// 	ID:       u.ID,
	// }

	// c.JSON(http.StatusOK, res)
	c.Redirect(http.StatusFound, "/hub")
}

func (h *Handler) Logout(c echo.Context) {
	c.SetCookie(&http.Cookie{Name: "jwt", Value: "", MaxAge: -1, Path: "/landing", Domain: "localhost", Secure: false, HttpOnly: true})
	c.JSON(http.StatusOK, "Logout successful")
}

// USER ROUTES HANDLER
func getUserHandler() (*Handler, error) {
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
		fmt.Println("FAILED")
		log.Fatalf("Could not get userHandler: %s", err)
	}
	userHandler.CreateUser(c)

	return c.Redirect(http.StatusFound, "/login")
}

func HandleLogin(c echo.Context) error {
	userHandler, err := getUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}
	userHandler.Login(c)
	return c.Redirect(http.StatusFound, "/hub")
}

func HandleLogout(c echo.Context) error {
	userHandler, err := getUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}
	userHandler.Logout(c)
	return c.Redirect(http.StatusFound, "/landing")
}
