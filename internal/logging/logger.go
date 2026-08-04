package logging

import (
	"log/slog"
	"os"

	"sketch-api-go/internal/config"
)

func New(cfg config.LoggerConfig) *slog.Logger {
	var handler slog.Handler

	switch cfg.Format {
	case "dev", "local":
		handler = slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		)
	case "prod":
		handler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		)
	}

	return slog.New(handler)
}
