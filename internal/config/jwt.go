package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
)

const (
	pathPrivatKeyEnv        string = "JWT_PATH_PRIVATE_KEY"
	pathPublicKeyEnv        string = "JWT_PATH_PUBLIC_KEY"
	accessExpiresMinutesEnv string = "ACCESS_EXPIRES_MINUTES"
	refreshExpiresDaysEnv   string = "REFRESH_EXPIRES_DAYS"
	// refreshCookieNameEnv    string = "REFRESH_COOKIE_NAME"
	issuerEnv   string = "JWT_ISSUER"
	audienceEnv string = "JWT_AUDIENCE"
)

type JWT interface {
	PrivateKey() *rsa.PrivateKey
	PublicKey() *rsa.PublicKey
	AccessExpiresMinutes() int
	RefreshExpiresMinutes() int
	// RefreshCookieName() string
	Issuer() string
	Audience() string
}

type configJWT struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	accessExpiresMinutes  int
	refreshExpiresMinutes int

	// refreshCookieName string
	issuer   string
	audience string
}

func (j *configJWT) PrivateKey() *rsa.PrivateKey {
	return j.privateKey
}

func (j *configJWT) PublicKey() *rsa.PublicKey {
	return j.publicKey
}

func (j *configJWT) AccessExpiresMinutes() int {
	return j.accessExpiresMinutes
}

func (j *configJWT) RefreshExpiresMinutes() int {
	return j.refreshExpiresMinutes
}

// func (j *configJWT) RefreshCookieName() string {
// 	 return j.refreshCookieName
// }

func (j *configJWT) Issuer() string {
	return j.issuer
}

func (j *configJWT) Audience() string {
	return j.audience
}

func (j *configJWT) validateJWT() error {
	switch {
	// case j.refreshCookieName == "":
	//	 return fmt.Errorf("environment variable %q is required", refreshCookieNameEnv)
	case j.issuer == "":
		return fmt.Errorf("environment variable %q is required", issuerEnv)
	case j.audience == "":
		return fmt.Errorf("environment variable %q is required", audienceEnv)
	}

	return nil
}

func readKeyFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read key file %q: %w", path, err)
	}

	return data, nil
}

func parsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM")
	}

	if block.Type != "RSA PRIVATE KEY" && block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected PEM type %q", block.Type)
	}

	// PKCS#1
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// PKCS#8
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}

	return rsaKey, nil
}

func parsePublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM")
	}

	if block.Type != "PUBLIC KEY" && block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("unexpected PEM type %q", block.Type)
	}

	// PKIX
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		rsaKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key is not RSA")
		}

		return rsaKey, nil
	}

	// PKCS#1
	rsaKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	return rsaKey, nil
}

func (j *configJWT) loadKeys() error {
	privateKeyPath := os.Getenv(pathPrivatKeyEnv)
	publicKeyPath := os.Getenv(pathPublicKeyEnv)

	if privateKeyPath == "" {
		return fmt.Errorf(
			"environment variable %q is required",
			pathPrivatKeyEnv,
		)
	}

	if publicKeyPath == "" {
		return fmt.Errorf(
			"environment variable %q is required",
			pathPublicKeyEnv,
		)
	}

	privatePEM, err := readKeyFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("private key: %w", err)
	}

	publicPEM, err := readKeyFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("public key: %w", err)
	}

	j.privateKey, err = parsePrivateKey(privatePEM)
	if err != nil {
		return fmt.Errorf("private key: %w", err)
	}

	j.publicKey, err = parsePublicKey(publicPEM)
	if err != nil {
		return fmt.Errorf("public key: %w", err)
	}

	return nil
}

func (j *configJWT) loadMinutes() error {
	accessExpiresMinutes, err := strconv.Atoi(os.Getenv(accessExpiresMinutesEnv))
	if err != nil {
		return fmt.Errorf(
			"environment variable %q must be a number: %w",
			accessExpiresMinutesEnv,
			err,
		)
	}

	if accessExpiresMinutes <= 0 {
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			accessExpiresMinutesEnv,
		)
	}

	refreshExpiresDays, err := strconv.Atoi(os.Getenv(refreshExpiresDaysEnv))
	if err != nil {
		return fmt.Errorf(
			"environment variable %q must be a number: %w",
			refreshExpiresDaysEnv,
			err,
		)
	}

	if refreshExpiresDays <= 0 {
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			refreshExpiresDaysEnv,
		)
	}

	j.accessExpiresMinutes = accessExpiresMinutes
	j.refreshExpiresMinutes = refreshExpiresDays

	return nil
}

func newJWTConfig() (*configJWT, error) {
	token := configJWT{
		// refreshCookieName: os.Getenv(refreshCookieNameEnv),
		issuer:   os.Getenv(issuerEnv),
		audience: os.Getenv(audienceEnv),
	}

	if err := token.validateJWT(); err != nil {
		return &token, err
	}

	if err := token.loadMinutes(); err != nil {
		return &token, err
	}

	if err := token.loadKeys(); err != nil {
		return &token, err
	}

	return &token, nil
}
