package service

import (
	"context"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/rirachii/golivechat/model"
	user "github.com/rirachii/golivechat/model/user"
)

type UserService interface {
	CreateUser(c context.Context, req *user.CreateUserReq) (*user.CreateUserRes, error)
	Login(c context.Context, req *user.LoginUserReq) (*user.LoginUserRes, error)
}

type userService struct {
	UserRepository UserRepository
	timeout        time.Duration
}


func NewUserService(repository UserRepository) UserService {
	return &userService{
		repository,
		time.Duration(2) * time.Second,
	}
}

const (
	secretKey = "TODO_change_to_something_better_secret"
)


func (s *userService) CreateUser(c context.Context, req *user.CreateUserReq) (*user.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashPassword,
	}

	r, err := s.UserRepository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &user.CreateUserRes{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *userService) Login(c context.Context, req *user.LoginUserReq) (*user.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.UserRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = CheckPassword(req.Password, u.Password)
	if err != nil {
		return nil, err
	}

	//generate jwt golang package
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JWTClaims{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	LoginRes := user.NewLoginUserRes(ss, u.Username, strconv.Itoa(int(u.ID)))
	return &LoginRes, nil
}
