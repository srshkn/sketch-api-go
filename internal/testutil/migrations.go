package testutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
)

func projectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("project root not found")
		}

		dir = parent
	}
}

func NewTestMigrations(dsn string) error {
	root, err := projectRoot()
	if err != nil {
		return fmt.Errorf(
			"create migration instance: %w",
			err,
		)
	}

	migrationsPath := filepath.ToSlash(
		filepath.Join(root, "db", "migrations"),
	)

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		return fmt.Errorf(
			"create migration instance: %w",
			err,
		)
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {

		return fmt.Errorf(
			"run migrations: %w",
			err,
		)
	}

	return nil
}
