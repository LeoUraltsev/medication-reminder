package handlers

import (
	"github.com/LeoUraltsev/medication-reminder/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type Handler struct {
	router chi.Router
	cfg    *config.App
}

func New(r chi.Router, cfg *config.App) *Handler {
	return &Handler{router: r, cfg: cfg}
}

func (h Handler) InitHandlers() {
	h.router.Route("/api/v1", func(r chi.Router) {
		r.Post("/schedule", h.CreateSchedule)
		r.Get("/schedules", h.SchedulesIds)
		r.Get("/schedule", h.ScheduleReception)
		r.Get("/next_takings", h.NextTakings)
	})
}

func (h Handler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "test create")
}

func (h Handler) SchedulesIds(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "test get ids schedules")
}

func (h Handler) ScheduleReception(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "test get schedule")
}

func (h Handler) NextTakings(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "test get next takings")
}
