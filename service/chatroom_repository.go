package service

import (
	"context"
	"fmt"

	model "github.com/rirachii/golivechat/model"
)

type ChatRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, email string) (*model.User, error)
}

type chatRepository struct {
	db DBTX
}

func NewChatRepository(db DBTX) UserRepository {
	return &userRepository{db: db}
}

func (r *chatRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	var lastInsertId int
	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&lastInsertId)
	if err != nil {
		return &model.User{}, err
	}

	user.ID = int64(lastInsertId)
	return user, nil
}

func (r *chatRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
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

func (r *chatRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {

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
