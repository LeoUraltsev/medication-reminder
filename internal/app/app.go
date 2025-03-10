package app

import (
	"fmt"
	"github.com/LeoUraltsev/medication-reminder/internal/config"
	"github.com/LeoUraltsev/medication-reminder/internal/handlers"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type App struct {
	cfg *config.Config
	log *slog.Logger
}

func NewApp(cfg *config.Config, log *slog.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

func (a App) Run() error {
	r := chi.NewRouter()

	h := handlers.New(nil, r, &a.cfg.App, a.log)
	h.InitHandlers()

	s := http.Server{
		Addr:         a.cfg.Http.Address,
		Handler:      r,
		ReadTimeout:  a.cfg.Http.ReadTimeout,
		WriteTimeout: a.cfg.Http.WriteTimeout,
		IdleTimeout:  a.cfg.Http.IdleTimeout,
	}

	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("failed sturtup server: %w", err)
	}

	return nil
}
