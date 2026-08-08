package service

import (
	"context"
	"errors"
	"strings"

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

func (u *UserService) validateCreateUser(
	ctx context.Context,
	userRow v1Generated.RegisterUserRequest,
) error {
	emailUser := strings.ToLower(string(userRow.Email))
	user, _ := u.repository.GetUserByEmailOrUsername(
		ctx,
		db.GetUserByEmailOrUsernameParams{
			Username: userRow.Username,
			Lower:    emailUser,
		},
	)

	if emailUser == user.Email {
		return errors.New("пользователь с такой почтой уже существует")
	}

	if userRow.Username == user.Username {
		return errors.New("пользователь с таким именем уже существует")
	}

	return nil
}

func (u *UserService) Registration(
	ctx context.Context,
	request v1Generated.RegisterUserRequest,
) (db.CreateUserRow, error) {
	var createUser db.CreateUserRow

	if err := u.validateCreateUser(ctx, request); err != nil {
		return createUser, err
	}

	hash, err := password.Hash(request.Password)
	if err != nil {
		return createUser, err
	}

	user := db.CreateUserParams{
		Username:     request.Username,
		Email:        strings.ToLower(string(request.Email)),
		PasswordHash: hash,
	}

	createUser, err = u.repository.CreateUser(ctx, user)
	if err != nil {
		return createUser, err
	}

	return createUser, nil
}
