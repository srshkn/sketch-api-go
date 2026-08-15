package repository

import (
	"context"

	"sketch-api-go/internal/db"

	"github.com/google/uuid"
)

type User interface {
	CheckUserExists(ctx context.Context, arg db.CheckUserExistsParams) (db.CheckUserExistsRow, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (db.GetUserByIDRow, error)
}
