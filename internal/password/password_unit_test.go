package password

import (
	"errors"
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	password := "my-secret-password"

	hash, err := Hash(password)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	if hash == "" {
		t.Fatal("Hash() returned empty hash")
	}

	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Errorf("Hash() = %q, want argon2id hash", hash)
	}

	ok, err := Compare(password, hash)
	if err != nil {
		t.Fatalf("Compare() error = %v", err)
	}

	if !ok {
		t.Error("Compare() = false, want true")
	}
}

func TestHash_ProducesDifferentHashes(t *testing.T) {
	password := "my-secret-password"

	hash1, err := Hash(password)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	hash2, err := Hash(password)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	if hash1 == hash2 {
		t.Error("Hash() produced identical hashes, salt is probably not random")
	}
}

func TestCompare(t *testing.T) {
	password := "my-secret-password"

	hash, err := Hash(password)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "correct password",
			password: password,
			want:     true,
		},
		{
			name:     "incorrect password",
			password: "wrong-password",
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compare(tt.password, hash)
			if err != nil {
				t.Fatalf("Compare() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompare_InvalidHash(t *testing.T) {
	tests := []struct {
		name string
		hash string
		want error
	}{
		{
			name: "empty hash",
			hash: "",
			want: ErrInvalidHash,
		},
		{
			name: "too few parts",
			hash: "$argon2id$v=19$m=65536,t=3,p=4",
			want: ErrInvalidHash,
		},
		{
			name: "unsupported algorithm",
			hash: "$bcrypt$v=19$m=65536,t=3,p=4$YWJj$YWJj",
			want: ErrUnsupportedHash,
		},
		{
			name: "invalid version format",
			hash: "$argon2id$version=19$m=65536,t=3,p=4$YWJj$YWJj",
			want: ErrInvalidHash,
		},
		{
			name: "unsupported version",
			hash: "$argon2id$v=999$m=65536,t=3,p=4$YWJj$YWJj",
			want: ErrUnsupportedVersion,
		},
		{
			name: "invalid parameters",
			hash: "$argon2id$v=19$invalid$YWJj$YWJj",
			want: ErrInvalidHash,
		},
		{
			name: "zero memory",
			hash: "$argon2id$v=19$m=0,t=3,p=4$YWJj$YWJj",
			want: ErrInvalidHash,
		},
		{
			name: "zero iterations",
			hash: "$argon2id$v=19$m=65536,t=0,p=4$YWJj$YWJj",
			want: ErrInvalidHash,
		},
		{
			name: "zero parallelism",
			hash: "$argon2id$v=19$m=65536,t=3,p=0$YWJj$YWJj",
			want: ErrInvalidHash,
		},
		{
			name: "invalid salt encoding",
			hash: "$argon2id$v=19$m=65536,t=3,p=4$!!!$YWJj",
			want: ErrInvalidHash,
		},
		{
			name: "invalid hash encoding",
			hash: "$argon2id$v=19$m=65536,t=3,p=4$YWJj$!!!",
			want: ErrInvalidHash,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Compare("password", tt.hash)

			if !errors.Is(err, tt.want) {
				t.Errorf("Compare() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestDecodeHash(t *testing.T) {
	hash, err := Hash("password")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	params, salt, expectedHash, err := decodeHash(hash)
	if err != nil {
		t.Fatalf("decodeHash() error = %v", err)
	}

	if params.memory != memory {
		t.Errorf("memory = %d, want %d", params.memory, memory)
	}

	if params.iterations != iterations {
		t.Errorf("iterations = %d, want %d", params.iterations, iterations)
	}

	if params.parallelism != parallelism {
		t.Errorf("parallelism = %d, want %d", params.parallelism, parallelism)
	}

	if len(salt) != saltLength {
		t.Errorf("salt length = %d, want %d", len(salt), saltLength)
	}

	if len(expectedHash) != keyLength {
		t.Errorf("hash length = %d, want %d", len(expectedHash), keyLength)
	}
}
