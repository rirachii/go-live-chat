package user_service

import (
	"context"
	"log"

	pgx "github.com/jackc/pgx/v5"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userData RepoCreateUser) (dbUserInfo, error)
	GetUserByEmail(ctx context.Context, email string) (dbUser, error)
	GetUserByID(ctx context.Context, id int) (dbUser, error)
	GetUsernameByID(ctx context.Context, id int) (dbUsername, error)

}

type userRepository struct {
	db *pgx.Conn
}

func (r *userRepository) DB() *pgx.Conn { return r.db }

func NewUserRepository(db *pgx.Conn) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, userData RepoCreateUser) (dbUserInfo, error) {

	queryAccount := `INSERT INTO users (username, user_account) 
				VALUES ($3, ROW($1, $2)::USER_ACCOUNT)
				RETURNING id, username;`


	var username string
	var userID int
	err := r.DB().QueryRow(
		ctx,
		queryAccount,
		userData.Email,
		userData.HashedPassword,
		userData.Username,
	).Scan(&userID, &username)
	if err != nil {
		return dbUserInfo{}, err
	}


	res := dbUserInfo{
		ID: userID,
		Username: username,
	}


	return res, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (dbUser, error) {
	
	res := dbUser{}

	query := `SELECT id, username, (user_account).email, (user_account).hashed_password 
				FROM users
				WHERE (user_account).email = $1`

	row := r.DB().QueryRow(ctx, query, email)
	err := row.Scan(&res.ID, &res.Username, &res.Email, &res.HashedPassword)

	if err != nil {
		log.Println(err)
		return dbUser{}, err
	}


	return res, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID int) (dbUser, error) {

	res := dbUser{}

	query := `SELECT id, username, (user_account).email, (user_account).hashed_password 
				FROM users
				WHERE id = $1`

	row := r.DB().QueryRow(ctx, query, userID)
	err := row.Scan(&res.ID, &res.Username, &res.Email, &res.HashedPassword)

	if err != nil {
		log.Println(err)
		return dbUser{}, err
	}


	return res, nil

}

func (r *userRepository) GetUsernameByID(ctx context.Context, id int) (dbUsername, error){

	res := dbUsername{}

	query := `SELECT username 
				FROM users
				WHERE id = $1`

	row := r.DB().QueryRow(ctx, query, id)
	err := row.Scan(&res.Username)
	if err != nil {
		return dbUsername{}, err
	}

	return res, nil
	
}
