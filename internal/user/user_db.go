package user_service

type RepoCreateUser struct {
	Email          string
	Username       string
	HashedPassword string
}

type dbUserID struct {
	ID int `db:"id"`
}

type dbUsername struct {
	Username string `db:"username"`
}

type dbUser struct {
	ID             int    `db:"id"`
	Username       string `db:"username"`
	Email          string `db:"(user_account).email"`
	HashedPassword string `db:"(user_account).hashed_password"`
}

type dbUserInfo struct {
	ID       int    `db:"id"`
	Email    string `db:"(user_account).email"`
	Username string `db:"username"`
}

type dbUserProfile struct {
	Username string `db:"username"`
	Email    string `db:"(user_account).email"`
}

type dbUserAccount struct {
	Email          string `db:"(user_account).email"`
	HashedPassword string `db:"(user_account).hashed_password"`
}
