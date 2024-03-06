package model

type User struct {
	ID       int64  `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type CreateUserReq struct {
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type CreateUserRes struct {
	ID       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
}

type LoginUserReq struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type LoginUserRes struct {
	ID          string `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
	accessToken string
}

func (user User) DisplayName() string {
	return user.Username
}

func NewLoginUserRes(token, id, username string) LoginUserRes {
	return LoginUserRes{
		ID:          id,
		Username:    username,
		accessToken: token,
	}
}

func (r *LoginUserRes) GetAccessToken() string {
	return r.accessToken
}
