package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"sketch-api-go/internal/db"
	"sketch-api-go/internal/password"
	"sketch-api-go/internal/repository"
	"sketch-api-go/internal/token"

	"github.com/google/uuid"

	v1Generated "sketch-api-go/internal/generated/v1"
)

var ErrRefreshTokenExpired = errors.New("refresh token expired")

type AuthService struct {
	repository   repository.AuthRepository
	tokenManager token.JWTManager
}

func NewAuthService(auth repository.AuthRepository, manager token.JWTManager) *AuthService {
	return &AuthService{
		repository:   auth,
		tokenManager: manager,
	}
}

func (a *AuthService) Login(
	ctx context.Context,
	request v1Generated.LoginUserRequest,
) (token.TokenPair, error) {
	var tokenPair token.TokenPair

	user, err := a.repository.GetUserByEmail(ctx, strings.ToLower(string(request.Email)))
	if err != nil {
		return tokenPair, err
	}

	if flag, err := password.Compare(request.Password, user.PasswordHash); !flag {
		return tokenPair, errors.New("неверный пароль")
	} else if err != nil {
		return tokenPair, err
	}

	return a.issueTokens(ctx, user.ID)
}

func (a *AuthService) Refresh(
	ctx context.Context,
	request v1Generated.RefreshRequest,
) (token.TokenPair, error) {
	var tokenPair token.TokenPair

	tokenHash := a.tokenManager.HashToken(token.RefreshToken(request.RefreshToken))

	stored, err := a.getValidRefresh(ctx, tokenHash)
	if err != nil {
		return tokenPair, err
	}

	user, err := a.getUserForToken(ctx, stored.UserID)
	if err != nil {
		return tokenPair, err
	}

	if err = a.repository.DeleteTokenHash(ctx, tokenHash); err != nil {
		return tokenPair, err
	}

	return a.issueTokens(ctx, user.ID)
}

func (a *AuthService) getValidRefresh(ctx context.Context, refreshToken string) (db.RefreshToken, error) {
	stored, err := a.repository.GetTokenHash(ctx, refreshToken)
	if err != nil {
		return stored, err
	}

	nowTime := time.Now()
	if !stored.ExpiresAt.After(nowTime) {
		if err = a.repository.DeleteTokenHash(ctx, refreshToken); err != nil {
			return stored, err
		}
		return stored, ErrRefreshTokenExpired
	}

	return stored, nil
}

func (a *AuthService) getUserForToken(ctx context.Context, userID uuid.UUID) (db.GetUserByIDRow, error) {
	user, err := a.repository.GetUserByID(ctx, userID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (a *AuthService) issueTokens(ctx context.Context, userID uuid.UUID) (token.TokenPair, error) {
	var tokenPair token.TokenPair
	var err error

	tokenPair.AccessToken, err = a.tokenManager.CreateAccessToken(userID.String())
	if err != nil {
		return tokenPair, err
	}

	tokenPair.RefreshToken, err = a.tokenManager.GenerateRefreshToken()
	if err != nil {
		return tokenPair, err
	}

	hashRefreshToken := a.tokenManager.HashToken(tokenPair.RefreshToken)
	expiresAt := a.tokenManager.RefreshTokenExpiresAt()

	_, err = a.repository.CreateRefreshToken(
		ctx, db.CreateRefreshTokenParams{
			UserID:    userID,
			TokenHash: hashRefreshToken,
			ExpiresAt: expiresAt,
		},
	)
	if err != nil {
		return tokenPair, err
	}

	tokenPair.ExpiresAt = expiresAt

	return tokenPair, nil
}
