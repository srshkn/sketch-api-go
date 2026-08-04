package service

import (
	"context"

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
