package main

import (
	"github.com/LeoUraltsev/medication-reminder/internal/config"
	"log/slog"
	"os"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("failed started app", slog.String("err", err.Error()))
		os.Exit(1)
	}

	log := initLogger(cfg.App.LogLevel)
	log.Info("startup application", slog.String("log level", cfg.App.LogLevel))
}

func initLogger(logLevel string) *slog.Logger {
	var log *slog.Logger

	switch logLevel {
	case "Debug":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case "Info":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	case "Warn":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}))
	case "Error":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelError,
		}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}
