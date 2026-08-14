package v1

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"sketch-api-go/internal/config"
	"sketch-api-go/internal/cookie"
	"sketch-api-go/internal/db"
	"sketch-api-go/internal/testutil"
	"sketch-api-go/internal/token"
)

var (
	testConfig        config.Config
	testDB            db.Querier
	testJWTManager    token.JWTManager
	testCookieManager cookie.Auth
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Configuration

	testConfig, err := testutil.NewTestConfig()
	if err != nil {
		slog.Error(
			"failed to load configuration",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	// -------------------------------------------------------------------------
	// JWT

	testJWTManager = token.New(testConfig.JWT())

	// -------------------------------------------------------------------------
	// Cookie

	testCookieManager = cookie.New(testConfig.Cookie())

	// -------------------------------------------------------------------------
	// Container PostgreSQL

	url, cleanup, err := testutil.NewTestContainer(ctx)
	if err != nil {
		if cleanup != nil {
			cleanup()
		}

		slog.Error(
			"failed to setup test PostgreSQL",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	// -------------------------------------------------------------------------
	// PostgreSQL

	pool, err := testutil.NewTestPostgres(ctx, testConfig.Postgres(), url)
	if err != nil {
		if pool != nil {
			pool.Close()
		}

		cleanup()
		slog.Error(
			"failed to setup test PostgreSQL",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	// -------------------------------------------------------------------------
	// Migrations

	if err = testutil.NewTestMigrations(url); err != nil {
		pool.Close()
		cleanup()
		slog.Error(
			"failed to setup test PostgreSQL",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	testDB = db.New(pool)

	// -------------------------------------------------------------------------
	// Tests

	code := m.Run()

	// -------------------------------------------------------------------------
	// Cleanup

	pool.Close()
	cleanup()

	os.Exit(code)
}
