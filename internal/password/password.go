package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	memory      uint32 = 64 * 1024
	iterations  uint32 = 3
	parallelism uint8  = 4
	saltLength         = 16
	keyLength          = 32
)

var (
	ErrInvalidHash        = errors.New("invalid password hash")
	ErrUnsupportedHash    = errors.New("unsupported password hash")
	ErrUnsupportedVersion = errors.New("unsupported argon2 version")
)

func Hash(password string) (string, error) {
	salt := make([]byte, saltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memory,
		iterations,
		parallelism,
		saltEncoded,
		hashEncoded,
	)

	return encoded, nil
}

func Compare(password, encodedHash string) (bool, error) {
	params, salt, expectedHash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		uint32(len(expectedHash)),
	)

	return subtle.ConstantTimeCompare(actualHash, expectedHash) == 1, nil
}

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
}

func decodeHash(encodedHash string) (params, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return params{}, nil, nil, ErrInvalidHash
	}

	if parts[1] != "argon2id" {
		return params{}, nil, nil, ErrUnsupportedHash
	}

	var version int

	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return params{}, nil, nil, ErrInvalidHash
	}

	if version != argon2.Version {
		return params{}, nil, nil, ErrUnsupportedVersion
	}

	var p params

	if _, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&p.memory,
		&p.iterations,
		&p.parallelism,
	); err != nil {
		return params{}, nil, nil, ErrInvalidHash
	}

	if p.memory == 0 || p.iterations == 0 || p.parallelism == 0 {
		return params{}, nil, nil, ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return params{}, nil, nil, ErrInvalidHash
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return params{}, nil, nil, ErrInvalidHash
	}

	if len(salt) == 0 || len(expectedHash) == 0 {
		return params{}, nil, nil, ErrInvalidHash
	}

	return p, salt, expectedHash, nil
}
