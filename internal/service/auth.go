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
	tokenManager *token.Manager
}

func NewAuthService(auth repository.AuthRepository, manager *token.Manager) *AuthService {
	return &AuthService{
		repository:   auth,
		tokenManager: manager,
	}
}

func (a *AuthService) Login(
	ctx context.Context,
	request v1Generated.LoginUserRequest,
) (token.AccessToken, token.RefreshToken, error) {
	var accessToken token.AccessToken
	var refreshToken token.RefreshToken

	user, err := a.repository.GetUserByEmail(ctx, strings.ToLower(string(request.Email)))
	if err != nil {
		return accessToken, refreshToken, err
	}

	if flag, err := password.Compare(request.Password, user.PasswordHash); !flag {
		return accessToken, refreshToken, errors.New("неверный пароль")
	} else if err != nil {
		return accessToken, refreshToken, err
	}

	return a.issueTokens(ctx, user.ID)
}

func (a *AuthService) Refresh(
	ctx context.Context,
	request v1Generated.RefreshRequest,
) (token.AccessToken, token.RefreshToken, error) {
	var accessToken token.AccessToken
	var refreshToken token.RefreshToken

	tokenHash := a.tokenManager.HashToken(token.RefreshToken(request.RefreshToken))

	stored, err := a.getValidRefresh(ctx, tokenHash)
	if err != nil {
		return accessToken, refreshToken, err
	}

	user, err := a.getUserForToken(ctx, stored.UserID)
	if err != nil {
		return accessToken, refreshToken, err
	}

	if err = a.repository.DeleteTokenHash(ctx, tokenHash); err != nil {
		return accessToken, refreshToken, err
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

func (a *AuthService) issueTokens(ctx context.Context, userID uuid.UUID) (token.AccessToken, token.RefreshToken, error) {
	var accessToken token.AccessToken
	var refreshToken token.RefreshToken

	accessToken, err := a.tokenManager.CreateAccessToken(userID.String())
	if err != nil {
		return accessToken, refreshToken, err
	}

	refreshToken, err = a.tokenManager.GenerateRefreshToken()
	if err != nil {
		return accessToken, refreshToken, err
	}

	hashRefreshToken := a.tokenManager.HashToken(refreshToken)

	_, err = a.repository.CreateRefreshToken(
		ctx, db.CreateRefreshTokenParams{
			UserID:    userID,
			TokenHash: hashRefreshToken,
			ExpiresAt: a.tokenManager.RefreshTokenExpiresAt(),
		},
	)
	if err != nil {
		return accessToken, refreshToken, err
	}

	return accessToken, refreshToken, nil
}
