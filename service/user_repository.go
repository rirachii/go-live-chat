package service

import (
	"context"
	"fmt"
	"log"

	pgx "github.com/jackc/pgx/v5"
	user "github.com/rirachii/golivechat/model/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *user.User) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	GetUserByID(ctx context.Context, email string) (*user.User, error)
}

type userRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *user.User) (*user.User, error) {
	var lastInsertId int

	// query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"

	const (
		cmd   = "INSERT INTO %s VALUES %s RETURNING id"
		table = "users"
		data  = "(($1, $2, $3)::UserAccount)"
	)

	query := fmt.Sprintf(cmd, table, data)
	err := r.db.QueryRow(
		ctx,
		query,
		user.Email,
		user.Username,
		user.Password,
	).Scan(&lastInsertId)
	if err != nil {
		return nil, err
	}

	user.ID = int64(lastInsertId)
	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	u := user.User{}

	// query := "SELECT id, email, username, password FROM users WHERE email = $1"

	const (
		cmd   = "SELECT %s FROM %s WHERE %s"
		table = "users"
		data  = "id, (user_data).email, (user_data).username, (user_data).hashed_password"
		cond  = "(user_data).email = $1"
	)

	query := fmt.Sprintf(cmd, data, table, cond)
	row := r.db.QueryRow(ctx, query, email)

	err := row.Scan(&u.ID, &u.Email, &u.Username, &u.Password)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Printf("%+v", u)

	return &u, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID string) (*user.User, error) {

	u := user.User{}

	query := "SELECT id, email, username, password FROM users WHERE id = $1"
	err := r.db.QueryRow(ctx, query, userID).
		Scan(&u.ID, &u.Email, &u.Username, &u.Password)

	if err != nil {
		return nil, err
	}
	fmt.Print(u.Username)

	return &u, nil

}
