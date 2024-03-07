package service

import (
	"context"
	"database/sql"
	"fmt"

	model "github.com/rirachii/golivechat/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, email string) (*model.User, error)
	GetUserByToken(ctx context.Context, email string) (*model.User, error)
}

// flexibility can pass in transaction instead of db object
// uses sql.DB methods/interfaces

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) UserRepository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	var lastInsertId int
	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&lastInsertId)
	if err != nil {
		return &model.User{}, err
	}

	user.ID = int64(lastInsertId)
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	u := model.User{}

	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&u.ID, &u.Email, &u.Username, &u.Password)

	if err != nil {
		return &model.User{}, nil
	}
	fmt.Print(u.Username)

	return &u, nil
}

func (r *repository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {

	u := model.User{}

	query := "SELECT id, email, username, password FROM users WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, userID).
		Scan(&u.ID, &u.Email, &u.Username, &u.Password)

	if err != nil {
		return &model.User{}, nil
	}
	fmt.Print(u.Username)

	return &u, nil

}

// TODO i think we should add a token field to user data or a new table of token -> userID
func (r *repository) GetUserByToken(ctx context.Context, token string) (*model.User, error) {

	u := model.User{}

	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, token).
		Scan(&u.ID, &u.Email, &u.Username, &u.Password)

	if err != nil {
		return &model.User{}, nil
	}
	fmt.Print(u.Username)

	return &u, nil

}
