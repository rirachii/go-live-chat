package handler

import (
	"fmt"
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
	db "github.com/rirachii/golivechat/db"
	userService "github.com/rirachii/golivechat/internal/user"
	user_model "github.com/rirachii/golivechat/model/user"
)

func HandleUserRegister(c echo.Context) error {
	userHandler, err := createUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}

	_, userErr := userHandler.CreateUser(c)
	if userErr != nil {
		return userErr
	}

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusFound)
}

func HandleUserLogin(c echo.Context) error {
	userHandler, err := createUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}

	loginRes, loginErr := userHandler.LoginUser(c)
	if loginErr != nil {
		return c.String(loginErr.Code, loginErr.Error())
	}

	accessToken := loginRes.AccessToken


	c.SetCookie(newJWTCookie(accessToken))
	c.Response().Header().Set("HX-Redirect", "/hub")

	return c.NoContent(http.StatusFound)
}

func HandleUserLogout(c echo.Context) error {
	userHandler, err := createUserHandler()
	if err != nil {
		log.Fatalf("Could not get userHandler: %s", err)
	}

	_ = userHandler.Logout(c)

	c.Response().Header().Set("HX-Redirect", "/landing")
	return c.NoContent(http.StatusFound)
}

// handles db requests
type UserHandler struct {
	UserService userService.UserService
}

func NewHandler(s userService.UserService) *UserHandler {
	return &UserHandler{
		UserService: s,
	}
}

func createUserHandler() (*UserHandler, error) {
	dbConn, err := db.ConnectDatabase()
	if err != nil {
		log.Fatalf("Could not initialize postgres db connection: %s", err)
	}

	userRep := userService.NewUserRepository(dbConn.DB())
	userSvc := userService.NewUserService(userRep)
	userHandler := NewHandler(userSvc)
	return userHandler, nil
}

func (h *UserHandler) CreateUser(c echo.Context) (user_model.UserCreated, *echo.HTTPError) {

	userCreated := user_model.UserCreated{
		Success: false,
	}

	var createUserReq user_model.CreateUserRequest
	err := c.Bind(&createUserReq)
	if err != nil {
		errorText := fmt.Sprintf("Bad request: %s", err.Error())
		err := echo.NewHTTPError(http.StatusBadRequest, errorText)
		return userCreated, err
	}

	if createUserReq.Email == "" ||
		createUserReq.Username == "" ||
		createUserReq.Password == "" {

		errorText := fmt.Sprintf("A field is empty: %+v", createUserReq)
		err := echo.NewHTTPError(http.StatusBadRequest, errorText)

		return userCreated, err
	}

	res, err := h.UserService.CreateUser(c.Request().Context(), createUserReq)
	if err != nil {
		err := echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprint("create user error: ", err),
		)
		return userCreated, err
	}

	userCreated.Success = true
	userCreated.ID = res.UserID
	userCreated.Email = res.Email
	userCreated.Username = res.Username

	return userCreated, nil

}

func (h *UserHandler) LoginUser(c echo.Context) (user_model.UserLoggedIn, *echo.HTTPError) {

	userLoggedIn := user_model.UserLoggedIn{
		Success: false,
	}

	var loginReq user_model.LoginUserRequest
	bindErr := c.Bind(&loginReq)
	if bindErr != nil {
		errorText := fmt.Sprintf("Bad request: %s", bindErr.Error())
		err := echo.NewHTTPError(http.StatusBadRequest, errorText)
		return userLoggedIn, err
	}

	if loginReq.Email == "" || loginReq.Password == "" {
		return userLoggedIn, echo.NewHTTPError(http.StatusBadGateway, "A field is empty:", loginReq)
	}

	ctx := c.Request().Context()
	loginRes, loginErr := h.UserService.Login(ctx, loginReq)
	if loginErr != nil {
		return userLoggedIn, echo.NewHTTPError(http.StatusInternalServerError, "error logging in:", loginErr)
	}

	userLoggedIn.Success = true
	userLoggedIn.ID = loginRes.UserID
	userLoggedIn.Username = loginRes.Username
	userLoggedIn.AccessToken = loginRes.GetAccessToken()

	return userLoggedIn, nil
}

func (h *UserHandler) Logout(c echo.Context) error {
	c.SetCookie(deadJWTCookie())

	return nil
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
