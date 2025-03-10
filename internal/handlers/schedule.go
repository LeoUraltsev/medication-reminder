package handlers

import (
	"context"
	"errors"
	"github.com/LeoUraltsev/medication-reminder/internal/domain"
	resp "github.com/LeoUraltsev/medication-reminder/internal/lib/response"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Services interface {
	AddSchedules(ctx context.Context, schedule domain.Schedule) (int64, error)
	SchedulesByUserID(ctx context.Context, userID int64) ([]int64, error)
	Schedule(ctx context.Context, userID int64, scheduleID int64) (domain.Schedule, error)
	NextTaking(ctx context.Context, userID int64) ([]string, error)
}

type ScheduleRequest struct {
	UserID               int64  `json:"user_id"`
	MedicineName         string `json:"medicine_name"`
	DurationOfReceptions uint8  `json:"duration_of_receptions"`
	DurationOfTreatment  string `json:"duration_of_treatment"`
}

type CreateResponse struct {
	resp.Response
	ScheduleID int64 `json:"schedule_id,omitempty"`
}

type SchedulesIdsResponse struct {
	resp.Response
	SchedulesIds []int64 `json:"schedules_ids,omitempty"`
}

type ScheduleResponse struct {
	resp.Response
	Schedule domain.Schedule
}

type NextTakingResponse struct {
	resp.Response
	Pills []string `json:"pills,omitempty"`
}

func (h Handler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.CreateSchedule"
	log := h.log.With(slog.String("op", op))
	log.Info("creating schedule")

	var req ScheduleRequest

	err := render.DecodeJSON(r.Body, &req)
	if errors.Is(err, io.EOF) {
		log.Error("request body empty", slog.String("err", err.Error()))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &CreateResponse{
			Response: resp.Error("request empty"),
		})
		return
	}

	if err != nil {
		log.Error("failed decode json", slog.String("err", err.Error()))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &CreateResponse{
			Response: resp.Error("failed decode json"),
		})
		return
	}
	//todo: сделать валидацию
	if req.UserID == 0 || req.MedicineName == "" || req.DurationOfReceptions == 0 || req.DurationOfTreatment == "" {
		log.Error("failed validate request", slog.Any("req", req))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &CreateResponse{
			Response: resp.Error("failed validate request"),
		})
		return
	}

	t, err := time.Parse(time.DateOnly, req.DurationOfTreatment)
	if err != nil {
		log.Error("failed parse date", slog.String("date", req.DurationOfTreatment))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &CreateResponse{
			Response: resp.Error("failed parse date duration of treatment"),
		})

		return
	}

	id, err := h.s.AddSchedules(context.Background(), domain.Schedule{
		UserID:               req.UserID,
		MedicineName:         req.MedicineName,
		DurationOfReceptions: req.DurationOfReceptions,
		DurationOfTreatment:  t,
	})

	if err != nil {
		log.Error("failed create schedule", slog.String("err", err.Error()))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &CreateResponse{
			Response: resp.Error("failed create schedule"),
		})

		return
	}

	log.Info("success create schedule", slog.Int64("id", id))

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, &CreateResponse{
		Response:   resp.OK(),
		ScheduleID: id,
	})
}

func (h Handler) SchedulesIds(w http.ResponseWriter, r *http.Request) {
	const op = "handlets.SchedulesIds"

	log := h.log.With(slog.String("op", op))

	log.Info("getting schedules")

	q := r.URL.Query().Get("user_id")
	if q == "" {
		log.Error("empty query param")

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, SchedulesIdsResponse{
			Response: resp.Error("query param empty"),
		})

		return
	}

	userID, err := strconv.ParseInt(q, 10, 64)
	if err != nil {
		log.Error("user id not int64", slog.String("err", err.Error()))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, SchedulesIdsResponse{Response: resp.Error("user_id param not number")})
		return
	}

	ids, err := h.s.SchedulesByUserID(context.Background(), userID)
	if errors.Is(err, domain.ErrScheduleNotFound) {
		log.Error("schedules for user not found", slog.Int64("user_id", userID))
		render.Status(r, http.StatusNotFound)
	}

	if err != nil {
		log.Error("failed getting schedules", slog.String("err", err.Error()))

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, SchedulesIdsResponse{Response: resp.Error("failed getting schedules")})
		return
	}

	log.Info("success getting schedules")
	render.Status(r, http.StatusOK)
	render.JSON(w, r, SchedulesIdsResponse{
		Response:     resp.OK(),
		SchedulesIds: ids,
	})
}

func (h Handler) Schedule(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.Schedule"

	log := h.log.With(slog.String("op", op))

	log.Info("getting schedule reception")

	u := r.URL.Query().Get("user_id")
	s := r.URL.Query().Get("schedule_id")

	userID, err := strconv.ParseInt(u, 10, 64)
	if err != nil {
		log.Error("failed parse users id in query param", slog.String("err", err.Error()))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ScheduleResponse{
			Response: resp.Error("failed query param"),
		})

		return
	}
	scheduleID, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Error("failed parse schedule id in query param", slog.String("err", err.Error()))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ScheduleResponse{
			Response: resp.Error("failed query param"),
		})
		return
	}

	schedule, err := h.s.Schedule(context.Background(), userID, scheduleID)
	if errors.Is(err, domain.ErrScheduleNotFound) {
		log.Error(
			"schedule not found",
			slog.Int64("user_id", userID),
			slog.Int64("schedule_id", scheduleID),
		)

		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, ScheduleResponse{
			Response: resp.Error("schedule not found"),
		})

		return
	}

	if err != nil {
		log.Error("failed getting schedule", slog.String("err", err.Error()))

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, SchedulesIdsResponse{Response: resp.Error("failed getting reception")})
		return
	}

	log.Info("success getting schedule reception")
	render.Status(r, http.StatusOK)
	render.JSON(w, r, ScheduleResponse{
		Response: resp.OK(),
		Schedule: schedule,
	})
}

func (h Handler) NextTakings(w http.ResponseWriter, r *http.Request) {
	var op = "handlers.NextTakings"
	log := h.log.With(slog.String("op", op), slog.String("url", r.URL.Path))
	u := r.URL.Query().Get("user_id")
	if u == "" {
		log.Error("empty query param")

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Error("query param empty"))
		return
	}

	userID, err := strconv.ParseInt(u, 10, 64)
	if err != nil {
		log.Error("user id not int64", slog.String("err", err.Error()))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, SchedulesIdsResponse{Response: resp.Error("user_id param not number")})
		return
	}

	pills, err := h.s.NextTaking(context.Background(), userID)
	if err != nil {
		log.Error("failed getting next taking", slog.String("err", err.Error()))

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed getting next taking pills"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, NextTakingResponse{
		Response: resp.OK(),
		Pills:    pills,
	})
}
