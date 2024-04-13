package model

import (
	errors "errors"
	fmt "fmt"

	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	secretKey = "TODO_change_to_something_better_secret"
)

type JWTClaims struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	DisplayName string `json:"display-name" db:"display_name"` 
	jwt.RegisteredClaims
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Check if the token is valid
	tokenClaims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		fmt.Println(err)
		return nil, errors.New("JWT TOKEN NOT VALID")
	}

	return tokenClaims, nil
}

func (claims JWTClaims) GetUID() string {
	return claims.ID
}
func (claims JWTClaims) GetUsername() string {
	return claims.Username
}
