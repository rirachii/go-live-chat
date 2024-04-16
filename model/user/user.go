package user_model

import "github.com/rirachii/golivechat/model"

type UserInfo struct {
	ID       model.UserID
	Username string
}

type UserCreated struct {
	Success  bool
	ID       model.UserID
	Email    string
	Username string
}

type UserLoggedIn struct {
	Success     bool
	ID          model.UserID
	Username    string
	AccessToken string
}
