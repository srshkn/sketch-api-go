package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

func newClaims(userID, tokenType string) *Claims {
	return &Claims{
		UserID:    userID,
		TokenType: tokenType,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sketch-api-go",
			Subject:   userID,
		},
	}
}

type Manager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewTokenManager(privatePEM, publicPEM []byte) (*Manager, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return nil, err
	}

	return &Manager{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (m *Manager) CreateAccessToken(userID string) (string, error) {
	claim := newClaims(userID, Access)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	return token.SignedString(m.privateKey)
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.publicKey, nil
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
