package users

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateUser(c echo.Context) {
	var u CreateUserReq
	if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
		c.JSONBlob(http.StatusBadRequest, []byte("create user error: "+err.Error()))
		return
	}

	res, err := h.Service.CreateUser(c.Request().Context(), &u)
	if err != nil {
		c.JSONBlob(http.StatusInternalServerError, []byte("create user error: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, res)

}

func (h *Handler) Login(c echo.Context) {
	var user LoginUserReq
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		c.JSONBlob(http.StatusBadRequest, []byte("login user error: "+err.Error()))
		return
	}

	u, err := h.Service.Login(c.Request().Context(), &user)
	if err != nil {
		c.JSONBlob(http.StatusInternalServerError, []byte("login user error: "+err.Error()))
		return
	}

	c.SetCookie(&http.Cookie{Name: "jwt", Value: u.accessToken, MaxAge: 3600, Path: "/landing", Domain: "localhost", Secure: false, HttpOnly: true})

	res := &LoginUserRes{
		Username: u.Username,
		ID:       u.ID,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Logout(c echo.Context) {
	c.SetCookie(&http.Cookie{Name: "jwt", Value: "", MaxAge: -1, Path: "/landing", Domain: "localhost", Secure: false, HttpOnly: true})
	c.JSON(http.StatusOK, "Logout successful")
}
