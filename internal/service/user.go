package service

import (
	"context"
	"errors"
	"strings"

	"sketch-api-go/internal/db"
	"sketch-api-go/internal/password"
	"sketch-api-go/internal/repository"

	v1Generated "sketch-api-go/internal/generated/v1"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User interface {
	Registration(ctx context.Context, request v1Generated.RegisterUserRequest) (v1Generated.UserResponse, error)
	GetUser(ctx context.Context, rowUserID string) (v1Generated.UserResponse, error)
}

type userService struct {
	repository repository.User
}

func NewUserService(user repository.User) *userService {
	return &userService{
		repository: user,
	}
}

func (u *userService) validateCreateUser(
	ctx context.Context,
	userRow v1Generated.RegisterUserRequest,
) error {
	emailUser := strings.ToLower(string(userRow.Email))
	user, _ := u.repository.CheckUserExists(
		ctx,
		db.CheckUserExistsParams{
			Username: userRow.Username,
			Email:    emailUser,
		},
	)

	if emailUser == user.Email {
		return errors.New("a user with this email already exists")
	}

	if userRow.Username == user.Username {
		return errors.New("a user with that username already exists")
	}

	return nil
}

func (u *userService) Registration(
	ctx context.Context,
	request v1Generated.RegisterUserRequest,
) (v1Generated.UserResponse, error) {
	var response v1Generated.UserResponse

	if err := u.validateCreateUser(ctx, request); err != nil {
		return response, err
	}

	hash, err := password.Hash(request.Password)
	if err != nil {
		return response, err
	}

	createUser, err := u.repository.CreateUser(
		ctx,
		db.CreateUserParams{
			Username:     request.Username,
			Email:        strings.ToLower(string(request.Email)),
			PasswordHash: hash,
		},
	)
	if err != nil {
		return response, err
	}

	response.Id = createUser.ID
	response.Username = createUser.Username

	return response, nil
}

func (u *userService) GetUser(
	ctx context.Context,
	rowUserID string,
) (v1Generated.UserResponse, error) {
	var response v1Generated.UserResponse

	userID, err := uuid.Parse(rowUserID)
	if err != nil {
		return response, err
	}

	user, err := u.repository.GetUserByID(ctx, uuid.UUID(userID))
	if err != nil {
		return response, err
	}

	response.Id = user.ID
	response.Username = user.Username

	return response, nil
}
