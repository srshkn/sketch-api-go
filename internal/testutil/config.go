package testutil

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"sketch-api-go/internal/config"
)

func NewTestConfig() (config.Config, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("failed to get current file path")
	}

	root := filepath.Join(filepath.Dir(filename), "../..")

	envPath := filepath.Join(root, ".test.env")

	cfg, err := config.NewFromFile(envPath)
	if err != nil {
		return nil, fmt.Errorf("load test config: %w", err)
	}

	return cfg, nil
}
