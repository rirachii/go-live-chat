package handler

import (
	"fmt"
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
	user "github.com/rirachii/golivechat/model/user"
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

func (h *UserHandler) CreateUser(c echo.Context) (*user.CreateUserRes, *echo.HTTPError) {

	var createUserReq user.CreateUserReq
	err := c.Bind(&createUserReq)
	if err != nil {
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
		err := echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprint("create user error: ", err),
		)
		return nil, err
	}

	return res, nil

}

func (h *UserHandler) LoginUser(c echo.Context) (*user.LoginUserRes, *echo.HTTPError) {

	var loginReq user.LoginUserReq
	bindErr := c.Bind(&loginReq)
	if bindErr != nil {
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
	c.SetCookie(deadJWTCookie())

	return nil
}

// USER ROUTES HANDLER
func getUserHandler() (*UserHandler, error) {
	dbConn, err := db.ConnectDatabase()
	if err != nil {
		log.Fatalf("Could not initialize postgres db connection: %s", err)
	}

	userRep := service.NewUserRepository(dbConn.DB())
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

	loginRes, loginErr := userHandler.LoginUser(c)
	if loginErr != nil {
		return c.String(loginErr.Code, loginErr.Error())
	}

	accessToken := loginRes.GetAccessToken()

	log.Println("acess token: ", accessToken)

	c.SetCookie(newJWTCookie(accessToken))
	c.Response().Header().Set("HX-Redirect", "/hub")

	return c.NoContent(http.StatusFound)
}

func HandleUserLogout(c echo.Context) error {
	userHandler, err := getUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}

	_ = userHandler.Logout(c)

	c.Response().Header().Set("HX-Redirect", "/landing")
	return c.NoContent(http.StatusFound)
}

func newJWTCookie(jwt string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		MaxAge:   3600,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}
	return cookie
}

func deadJWTCookie() *http.Cookie {
	deadCookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		MaxAge:   -1,
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}
	return deadCookie
}
