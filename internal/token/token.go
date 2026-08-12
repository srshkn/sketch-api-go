package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"sketch-api-go/internal/config"
)

type AccessToken string
type RefreshToken string

const (
	Access  string = "access"
	Refresh string = "refresh"
)

type Claims struct {
	UserID    string `json:"uid"`
	TokenType string `json:"typ"`

	jwt.RegisteredClaims
}

func newClaims(userID, tokenType, issuer string, accessMinutes int) *Claims {
	return &Claims{
		UserID:    userID,
		TokenType: tokenType,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(accessMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   userID,
		},
	}
}

type Manager struct {
	config      config.JWT
	accessName  string
	refreshName string
}

func New(cfg config.JWT) *Manager {
	return &Manager{
		config:      cfg,
		accessName:  Access,
		refreshName: Refresh,
	}
}

func (m *Manager) CreateAccessToken(userID string) (AccessToken, error) {
	claim := newClaims(
		userID,
		m.accessName,
		m.config.Issuer(),
		m.config.AccessExpiresMinutes(),
	)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	accessToken, err := token.SignedString(m.config.PrivateKey)
	if err != nil {
		return AccessToken(accessToken), err
	}

	return AccessToken(accessToken), nil
}

func (m *Manager) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.config.PublicKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func (_ *Manager) GenerateRefreshToken() (RefreshToken, error) {
	bytes := make([]byte, 48)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return RefreshToken(base64.RawURLEncoding.EncodeToString(bytes)), nil
}

func (_ *Manager) HashToken(refreshToken RefreshToken) string {
	hash := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(hash[:])
}

func (m *Manager) RefreshTokenExpiresAt() time.Time {
	return time.Now().Add(time.Duration(m.config.RefreshExpiresMinutes()) * 24 * time.Hour)
}
