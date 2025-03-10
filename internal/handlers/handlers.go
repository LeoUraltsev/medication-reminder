package handlers

import (
	"github.com/LeoUraltsev/medication-reminder/internal/config"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

type Handler struct {
	s      Services
	router chi.Router
	cfg    *config.App
	log    *slog.Logger
}

func New(s Services, r chi.Router, cfg *config.App, log *slog.Logger) *Handler {
	return &Handler{s: s, router: r, cfg: cfg, log: log}
}

func (h Handler) InitHandlers() {
	h.router.Route("/api/v1", func(r chi.Router) {
		r.Post("/schedule", h.CreateSchedule)
		r.Get("/schedules", h.SchedulesIds)
		r.Get("/schedule", h.Schedule)
		r.Get("/next_takings", h.NextTakings)
	})
}
