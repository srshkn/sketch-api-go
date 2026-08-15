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

type Auth interface {
	Login(ctx context.Context, request v1Generated.LoginUserRequest) (v1Generated.TokensResponse, token.RefreshToken, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, request v1Generated.RefreshRequest) (v1Generated.TokensResponse, token.RefreshToken, error)
}

type authService struct {
	repository   repository.Auth
	tokenManager token.JWTManager
}

func NewAuthService(auth repository.Auth, manager token.JWTManager) *authService {
	return &authService{
		repository:   auth,
		tokenManager: manager,
	}
}

func (a *authService) Login(
	ctx context.Context,
	request v1Generated.LoginUserRequest,
) (
	v1Generated.TokensResponse,
	token.RefreshToken,
	error,
) {
	var response v1Generated.TokensResponse
	var refreshToken token.RefreshToken
	var err error

	user, err := a.repository.GetUserByLogin(ctx, strings.ToLower(string(request.Email)))
	if err != nil {
		return response, refreshToken, err
	}

	if flag, err := password.Compare(request.Password, user.PasswordHash); !flag {
		return response, refreshToken, errors.New("неверный пароль")
	} else if err != nil {
		return response, refreshToken, err
	}

	return a.issueTokens(ctx, user.ID)
}

func (a *authService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := a.tokenManager.HashToken(refreshToken)

	return a.repository.DeleteTokenHash(ctx, tokenHash)
}

func (a *authService) Refresh(
	ctx context.Context,
	request v1Generated.RefreshRequest,
) (
	v1Generated.TokensResponse,
	token.RefreshToken,
	error,
) {
	var response v1Generated.TokensResponse
	var refreshToken token.RefreshToken
	var err error

	tokenHash := a.tokenManager.HashToken(request.RefreshToken)

	stored, err := a.getValidRefresh(ctx, tokenHash)
	if err != nil {
		return response, refreshToken, err
	}

	user, err := a.getUserForToken(ctx, stored.UserID)
	if err != nil {
		return response, refreshToken, err
	}

	if err = a.repository.DeleteTokenHash(ctx, tokenHash); err != nil {
		return response, refreshToken, err
	}

	return a.issueTokens(ctx, user.ID)
}

func (a *authService) getValidRefresh(ctx context.Context, refreshToken string) (db.RefreshToken, error) {
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

func (a *authService) getUserForToken(ctx context.Context, userID uuid.UUID) (db.GetUserByIDRow, error) {
	user, err := a.repository.GetUserByID(ctx, userID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (a *authService) issueTokens(
	ctx context.Context,
	userID uuid.UUID,
) (
	v1Generated.TokensResponse,
	token.RefreshToken,
	error,
) {
	var response v1Generated.TokensResponse
	var refreshToken token.RefreshToken
	var err error

	response.AccessToken, err = a.tokenManager.CreateAccessToken(userID.String())
	if err != nil {
		return response, refreshToken, err
	}

	refreshToken.Token, err = a.tokenManager.GenerateRefreshToken()
	if err != nil {
		return response, refreshToken, err
	}

	hashRefreshToken := a.tokenManager.HashToken(refreshToken.Token)
	refreshToken.ExpiresAt = a.tokenManager.RefreshTokenExpiresAt()

	_, err = a.repository.CreateRefreshToken(
		ctx, db.CreateRefreshTokenParams{
			UserID:    userID,
			TokenHash: hashRefreshToken,
			ExpiresAt: refreshToken.ExpiresAt,
		},
	)
	if err != nil {
		return response, refreshToken, err
	}

	return response, refreshToken, nil
}
