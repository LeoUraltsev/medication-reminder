package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/LeoUraltsev/medication-reminder/internal/config"
	"github.com/LeoUraltsev/medication-reminder/internal/handlers"
	"github.com/LeoUraltsev/medication-reminder/internal/services"
	"github.com/go-chi/chi/v5"
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

	service := services.New(nil, &a.cfg.App)

	h := handlers.New(service, r, &a.cfg.App, a.log)
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
