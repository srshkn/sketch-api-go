package service

import (
	"context"
	"errors"
	"strings"

	"sketch-api-go/internal/db"
	"sketch-api-go/internal/password"
	"sketch-api-go/internal/repository"
	"sketch-api-go/internal/token"

	v1Generated "sketch-api-go/internal/generated/v1"
)

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

	accessToken, err = a.tokenManager.CreateAccessToken(user.ID.String())
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
			UserID:    user.ID,
			TokenHash: hashRefreshToken,
			ExpiresAt: a.tokenManager.RefreshTokenExpiresAt(),
		},
	)
	if err != nil {
		return accessToken, refreshToken, err
	}

	return accessToken, refreshToken, nil
}
