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
