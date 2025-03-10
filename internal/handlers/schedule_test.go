package handlers

import (
	"context"
	"encoding/json"
	"github.com/LeoUraltsev/medication-reminder/internal/config"
	"github.com/LeoUraltsev/medication-reminder/internal/domain"
	mock_handlers "github.com/LeoUraltsev/medication-reminder/internal/handlers/mocks"
	resp "github.com/LeoUraltsev/medication-reminder/internal/lib/response"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCreateSchedule(t *testing.T) {

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	mock := mock_handlers.NewMockServices(gomock.NewController(t))

	handler := New(
		mock,
		chi.NewRouter(),
		&config.App{
			LogLevel:   "debug",
			PeriodTime: 1 * time.Hour,
		},
		log,
	)

	sr := ScheduleRequest{
		UserID:               1234,
		MedicineName:         "aspirine",
		DurationOfReceptions: 4,
		DurationOfTreatment:  "2022-01-30",
	}

	body, _ := json.Marshal(sr)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(
		http.MethodPost,
		"http://localhost:10000/api/v1/schedule",
		strings.NewReader(string(body)),
	)

	date, _ := time.Parse(time.DateOnly, sr.DurationOfTreatment)
	mock.
		EXPECT().
		AddSchedules(context.Background(), domain.Schedule{
			UserID:               sr.UserID,
			MedicineName:         sr.MedicineName,
			DurationOfReceptions: sr.DurationOfReceptions,
			DurationOfTreatment:  date,
		}).Return(int64(1), nil)

	handler.CreateSchedule(w, r)

	expRes := CreateResponse{
		Response: resp.Response{
			Status: "OK",
		},
		ScheduleID: int64(1),
	}
	actualRes := CreateResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), &actualRes)

	assert.Equal(t, expRes, actualRes)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestErrorCreateSchedule(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	mock := mock_handlers.NewMockServices(gomock.NewController(t))

	handler := New(
		mock,
		chi.NewRouter(),
		&config.App{
			LogLevel:   "debug",
			PeriodTime: 1 * time.Hour,
		},
		log,
	)

	testCases := []string{
		"",
		`{}`,
		`{"user_id":""}`,
		`{"user_id":"",}`,
		`{"user_id":"","medicine_name":"","duration_of_receptions":"","duration_of_treatment":""}`,
	}

	for _, v := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(
			http.MethodPost,
			"http://localhost:10000/api/v1/schedule",
			strings.NewReader(v),
		)
		handler.CreateSchedule(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	}

}
