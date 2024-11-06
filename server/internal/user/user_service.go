package user

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"server/internal/util"
	"strconv"
	"time"
)

const (
	secretKey = "secret"
)

type service struct {
	Repository
	timeout time.Duration
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewUserService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(12) * time.Second,
	}
}

func (u *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, u.timeout)
	defer cancelFunc()

	//TODO: hash password
	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	data := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashPassword,
	}

	user, err := u.Repository.CreateUser(ctx, data)

	if err != nil {
		return nil, err
	}

	res := &CreateUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	return res, nil
}

func (s *service) Login(c context.Context, req *LoginUserRequest) (*LoginUserResponse, error) {
	ctx, cancelFunc := context.WithTimeout(c, s.timeout)
	defer cancelFunc()

	user, err := s.Repository.GetUserByEmail(ctx, req.Email)

	if err != nil {
		return &LoginUserResponse{}, err
	}

	err = util.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return &LoginUserResponse{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(user.ID)),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &LoginUserResponse{}, err
	}

	res := &LoginUserResponse{
		AccessToken: tokenString,
		ID:          user.ID,
		Username:    user.Username,
	}

	return res, nil
}
