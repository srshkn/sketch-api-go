package repository

import (
	"context"

	"sketch-api-go/internal/db"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error)
	GetUserByEmailOrUsername(ctx context.Context, arg db.GetUserByEmailOrUsernameParams) (db.GetUserByEmailOrUsernameRow, error)
}
