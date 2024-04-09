package handler

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	model "github.com/rirachii/golivechat/model"
)

type JWTCookie struct {
	cookie *http.Cookie
}

func GetJWTCookie(c echo.Context) (*JWTCookie, error) {
	jwtCookie, err := c.Cookie("jwt")
	if err != nil {
		return nil, err
	}

	cookie := &JWTCookie{
		cookie: jwtCookie,
	}

	return cookie, nil
}

func (jwtCookie JWTCookie) Cookie() *http.Cookie {
	return jwtCookie.cookie
}

func (jwtCookie JWTCookie) JWT() string {
	return jwtCookie.cookie.Value
}

func (jwtCookie JWTCookie) Claims() (*model.JWTClaims, error) {

	tokenString := jwtCookie.JWT()

	validTokenClaims, validationErr := model.ValidateJWT(tokenString)
	if validationErr != nil {
		echo.New().Logger.Print(validationErr)
		return nil, validationErr
	}

	return validTokenClaims, nil
}


// calls `GetJWTCookie().Claims()` to centralize error handling
func GetJWTClaims(c echo.Context) (*model.JWTClaims, error){

	cookie, err := GetJWTCookie(c)
	if err != nil {
		return nil, err
	}

	claims, err := cookie.Claims(); if err != nil {
		return nil, err
	}

	return claims, nil

}

// calls `GetJWTClaims(c).UserID`
func GetJWTUserID(c echo.Context) (string, error) {

	claims, err := GetJWTClaims(c)
	if err != nil {
		return "", err
	}

	return claims.GetUID(), nil

}