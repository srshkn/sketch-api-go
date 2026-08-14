package testutil

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func NewTestContainer(ctx context.Context) (string, func(), error) {
	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return "", nil, fmt.Errorf(
			"start postgres container: %w",
			err,
		)
	}

	cleanup := func() {
		_ = container.Terminate(ctx)
	}

	url, err := container.ConnectionString(
		ctx,
		"sslmode=disable",
	)
	if err != nil {
		cleanup()

		return "", nil, fmt.Errorf(
			"get postgres connection string: %w",
			err,
		)
	}

	return url, cleanup, nil
}
