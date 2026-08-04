package repository

import (
	"context"

	"github.com/google/uuid"

	"sketch-api-go/internal/db"
)

type UserRepository interface {
	ByEmail(ctx context.Context, lower string) (db.ByEmailRow, error)
	ByName(ctx context.Context, name string) (db.ByNameRow, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUserByEmail(ctx context.Context, lower string) (db.GetUserByEmailRow, error)
	ListUsers(ctx context.Context) ([]db.ListUsersRow, error)
}
