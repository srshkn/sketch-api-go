package jwt

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

const (
	AccessToken  = "access"
	RefreshToken = "refresh"
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

type TokenManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewTokenManager(privatePEM, publicPEM []byte) (*TokenManager, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return nil, err
	}

	return &TokenManager{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (t *TokenManager) CreateAccessToken(userID string) (string, error) {
	claim := newClaims(userID, AccessToken)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	return token.SignedString(t.privateKey)
}

func (t *TokenManager) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return t.publicKey, nil
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

func (t *TokenManager) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 48)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (t *TokenManager) HashToken(refreshToken string) string {
	hash := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(hash[:])
}
