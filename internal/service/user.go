package service

import (
	"context"
	"errors"

	"sketch-api-go/internal/db"
	"sketch-api-go/internal/password"
	"sketch-api-go/internal/repository"

	v1Generated "sketch-api-go/internal/generated/v1"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(user repository.UserRepository) *UserService {
	return &UserService{
		repository: user,
	}
}

func (u *UserService) Registration(
	ctx context.Context,
	request v1Generated.RegisterUserRequest,
) (db.CreateUserRow, error) {
	var createUser db.CreateUserRow

	hash, err := password.Hash(request.Password)
	if err != nil {
		return createUser, err
	}

	user := db.CreateUserParams{
		Name:         request.Username,
		Email:        string(request.Email),
		PasswordHash: hash,
	}

	createUser, err = u.repository.CreateUser(ctx, user)
	if err != nil {
		return createUser, err
	}

	return createUser, nil
}

func (u *UserService) Login(
	ctx context.Context,
	request v1Generated.LoginUserRequest,
) (db.GetUserByEmailRow, error) {

	user, err := u.repository.GetUserByEmail(ctx, string(request.Email))
	if err != nil {
		return user, err
	}

	if flag, err := password.Compare(request.Password, user.PasswordHash); !flag {
		return user, errors.New("неверный пароль")
	} else if err != nil {
		return user, err
	}

	return user, nil
}
