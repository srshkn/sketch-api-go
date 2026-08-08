package repository

import (
	"context"

	"github.com/google/uuid"

	"sketch-api-go/internal/db"
)

type AuthRepository interface {
	CreateRefreshToken(ctx context.Context, arg db.CreateRefreshTokenParams) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, lower string) (db.GetUserByEmailRow, error)
	DeleteTokenHash(ctx context.Context, tokenHash string) error
	GetTokenHash(ctx context.Context, tokenHash string) (db.RefreshToken, error)
	GetTokenUserID(ctx context.Context, userID uuid.UUID) (db.GetTokenUserIDRow, error)
}
