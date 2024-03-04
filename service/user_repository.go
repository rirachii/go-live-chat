package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rirachii/golivechat/model"

)


type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

// flexibility can pass in transaction instead of db object
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
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
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return &model.User{}, nil
	}
	fmt.Print(u.Username)

	return &u, nil
}
