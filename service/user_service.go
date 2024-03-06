package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	model "github.com/rirachii/golivechat/model"
)

type UserService interface {
	CreateUser(c context.Context, req *model.CreateUserReq) (*model.CreateUserRes, error)
	Login(c context.Context, req *model.LoginUserReq) (*model.LoginUserRes, error)
}

type service struct {
	UserRepository UserRepository
	timeout        time.Duration
}

func NewService(repository UserRepository) UserService {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, req *model.CreateUserReq) (*model.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashPassword,
	}

	r, err := s.UserRepository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &model.CreateUserRes{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

const (
	secretKey = "TODO_change_to_something_better_secret"
)

type MyJWTClaims struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	jwt.RegisteredClaims
}

// Valid implements jwt.Claims.
func (MyJWTClaims) Valid() error {
	//TODO Check if jwt in cookie is valid
	panic("unimplemented")
}

func (s *service) Login(c context.Context, req *model.LoginUserReq) (*model.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.UserRepository.GetUserByEmail(ctx, req.Email)
	fmt.Println("Login called")
	if err != nil {
		return &model.LoginUserRes{}, err
	}

	err = CheckPassword(req.Password, u.Password)
	if err != nil {
		return &model.LoginUserRes{}, err
	}

	//generate jwt golang package
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &model.LoginUserRes{}, err
	}

	LoginRes := model.NewLoginUserRes(ss, u.Username, strconv.Itoa(int(u.ID)))
	return &LoginRes, nil
}
