package service

import (
	"context"
	"sketch-api-go/internal/db"
	"sketch-api-go/internal/generated"
	"sketch-api-go/internal/repository"
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
	request generated.RegisterUserRequest,
) (db.CreateUserRow, error) {
	user := db.CreateUserParams{
		Name:         request.Username,
		Email:        string(request.Email),
		PasswordHash: request.Password,
	}
	createUser, err := u.repository.CreateUser(ctx, user)
	if err != nil {
		return createUser, err
	}

	return createUser, nil
}
