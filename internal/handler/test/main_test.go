package test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"sketch-api-go/internal/config"
	"sketch-api-go/internal/cookie"
	"sketch-api-go/internal/db"
	"sketch-api-go/internal/logging"
	"sketch-api-go/internal/testutil"
	"sketch-api-go/internal/token"
)

var (
	testConfig        config.Config
	testServet        config.Server
	testLogger        *slog.Logger
	testDB            db.Querier
	testJWTManager    token.JWTManager
	testCookieManager cookie.Auth
	testCORS          config.CORS
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
	// Logger

	testLogger = logging.New(testConfig.Logger())

	// -------------------------------------------------------------------------
	// JWT

	testJWTManager = token.New(testConfig.JWT())

	// -------------------------------------------------------------------------
	// Cookie

	testCookieManager = cookie.New(testConfig.Cookie())

	// -------------------------------------------------------------------------
	// CORSCfg

	testCORS = testConfig.CORS()

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
	// ServerCofg

	testServet = testConfig.Server()

	// -------------------------------------------------------------------------
	// Tests

	code := m.Run()

	// -------------------------------------------------------------------------
	// Cleanup

	pool.Close()
	cleanup()

	os.Exit(code)
}
