package user_model

import "github.com/rirachii/golivechat/model"

type CreateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserDTO struct {
	UserID   model.UserID
	Email    string
	Username string
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserDTO struct {
	UserID      model.UserID
	Username    string
	accessToken string
}

func NewUserLoggedIn(uid model.UserID, username string, signedToken string) LoginUserDTO {

	loggedInUserDTO := LoginUserDTO{
		UserID:      uid,
		Username:    username,
		accessToken: signedToken,
	}

	return loggedInUserDTO
}

func (r *LoginUserDTO) GetAccessToken() string {
	return r.accessToken
}
