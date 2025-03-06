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

	_ = cfg
}
