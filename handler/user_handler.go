package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
	model "github.com/rirachii/golivechat/model"
	service "github.com/rirachii/golivechat/service"
	db "github.com/rirachii/golivechat/service/db"
)

type UserHandler struct {
	UserService service.UserService
}

func NewHandler(s service.UserService) *UserHandler {
	return &UserHandler{
		UserService: s,
	}
}


func (h *UserHandler) CreateUser(c echo.Context) (*model.CreateUserRes, *echo.HTTPError) {

	var createUserReq model.CreateUserReq
	err := c.Bind(&createUserReq); if err != nil {
		errorText := fmt.Sprintf("Bad request: %s", err.Error())
		err := echo.NewHTTPError(http.StatusBadRequest, errorText)
		return nil, err
	}

	if createUserReq.Email == "" ||
		createUserReq.Username == "" ||
		createUserReq.Password == "" {

		errorText := fmt.Sprintf("A field is empty: %+v", createUserReq)
		err := echo.NewHTTPError(http.StatusBadRequest, errorText)

		return nil, err
	}


	res, err := h.UserService.CreateUser(c.Request().Context(), &createUserReq)
	if err != nil {
		err := echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprint("create user error: ", err))
		return nil, err
	}

	return res, nil

}

func (h *Handler) LoginUser(c echo.Context) (*model.LoginUserRes, *echo.HTTPError){


	var loginReq model.LoginUserReq
	bindErr := c.Bind(&loginReq); if bindErr != nil {
		errorText := fmt.Sprintf("Bad request: %s", bindErr.Error())
		err := echo.NewHTTPError(http.StatusBadRequest, errorText)
		return nil, err
	}


	if loginReq.Email == "" || loginReq.Password == "" {
		return nil, echo.NewHTTPError(http.StatusBadGateway, "A field is empty:", loginReq)
	}


	loginRes, loginErr := h.UserService.Login(c.Request().Context(), &loginReq)
	if loginErr != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "error logging in:", loginErr)
	}


	return loginRes, nil
}

func (h *UserHandler) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{Name: "jwt", Value: "", MaxAge: -1, Domain: "localhost", Secure: false, HttpOnly: true})
	return c.JSON(http.StatusOK, "Logout successful")
}


// USER ROUTES HANDLER
func getUserHandler() (*UserHandler, error) {
	dbConn, err := db.ConnectDatabase()
	if err != nil {
		log.Fatalf("Could not initialize postgres db connection: %s", err)
	}

	userRep := service.NewRepository(dbConn.GetDB())
	userSvc := service.NewService(userRep)
	userHandler := NewHandler(userSvc)
	return userHandler, nil
}

func HandleUserRegister(c echo.Context) error {
	userHandler, err := getUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}


	_, userErr := userHandler.CreateUser(c)
	if userErr != nil {
		return c.String(userErr.Code, userErr.Error())
	}

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusFound)
}

func HandleUserLogin(c echo.Context) error {
	userHandler, err := getUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}

	loginRes, loginErr := userHandler.LoginUser(c); if loginErr != nil {
		return c.String(loginErr.Code, loginErr.Error())
	}

	accessToken := loginRes.GetAccessToken()
	log.Println(accessToken)

	c.SetCookie(&http.Cookie{
		Name: "jwt", 
		Value: accessToken, 
		MaxAge: 3600, 
		Path: "/landing", 
		Domain: "localhost", 
		Secure: false, 
		HttpOnly: true,
	})

	c.Response().Header().Set("HX-Redirect", "/hub")
	return c.NoContent(http.StatusFound)
}

func HandleUserLogout(c echo.Context) error {
	userHandler, err := getUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}

	userHandler.Logout(c)
	c.Response().Header().Set("HX-Redirect", "/landing")
	return c.NoContent(http.StatusFound)
}
