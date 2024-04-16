package user_service

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	model "github.com/rirachii/golivechat/model"

	user_model "github.com/rirachii/golivechat/model/user"
)

type UserService interface {
	CreateUser(c context.Context, req user_model.CreateUserRequest) (user_model.CreateUserDTO, error)
	Login(c context.Context, req user_model.LoginUserRequest) (user_model.LoginUserDTO, error)
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


func (s *userService) CreateUser(c context.Context, 
	req user_model.CreateUserRequest,
) (user_model.CreateUserDTO, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return user_model.CreateUserDTO{}, err
	}

	createUserData := RepoCreateUser{
		Email:          req.Email,
		Username:       req.Username,
		HashedPassword: hashedPassword,
	}

	r, err := s.UserRepository.CreateUser(ctx, createUserData)
	if err != nil {
		return user_model.CreateUserDTO{}, err
	}

	res := user_model.CreateUserDTO{
		UserID:       model.IntToUID(r.ID),
		Email:    r.Email,
		Username: r.Username,
	}

	return res, nil
}

func (s *userService) Login(
	c context.Context, req user_model.LoginUserRequest,
) (user_model.LoginUserDTO, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	user, err := s.UserRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return user_model.LoginUserDTO{}, err
	}

	err = CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return user_model.LoginUserDTO{}, err
	}

	//generate jwt golang package
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JWTClaims{
		ID:          strconv.Itoa(user.ID),
		Username:    user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(user.ID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	secretKey := os.Getenv("JWT_SECRET_KEY")
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return user_model.LoginUserDTO{}, errors.New("failed to sign token")
	}

	uid :=  model.IntToUID(user.ID)
	username := user.Username

	userLoginData := user_model.NewUserLoggedIn(uid, username, signedToken)

	return userLoginData, nil
}
