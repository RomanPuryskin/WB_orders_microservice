package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	logLevel string `env:"LOG_LEVEL" envDefault:"INFO"`
}

func InitLogger(cfg *Config) {

	var level slog.Level
	switch cfg.logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
