package repository

import (
	"context"

	"sketch-api-go/internal/db"

	"github.com/google/uuid"
)

type User interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error)
	GetUserByEmailOrUsername(ctx context.Context, arg db.GetUserByEmailOrUsernameParams) (db.GetUserByEmailOrUsernameRow, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (db.GetUserByIDRow, error)
}
