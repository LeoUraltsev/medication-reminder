package services

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/LeoUraltsev/medication-reminder/internal/config"
	"github.com/LeoUraltsev/medication-reminder/internal/domain"
)

const (
	maxCountTakings          = 14
	maxMedicationTakingHours = maxCountTakings * time.Hour
)

type ScheduleSaver interface {
	InsertSchedule(ctx context.Context, schedule domain.Schedule, times []string) (int64, error)
	AllShedulesByUserID(ctx context.Context, userID int64) ([]int64, error)
	ScheduleByUserIDAndScheduleID(ctx context.Context, userID int64, scheduleID int64) (domain.Schedule, error)
	NextTakingByUserID(ctx context.Context, userID int64, duration time.Duration) ([]string, error)
}

type Service struct {
	ss  ScheduleSaver
	cfg *config.App
}

func New(sSaver ScheduleSaver, cfg *config.App) *Service {
	return &Service{
		ss:  nil,
		cfg: cfg,
	}
}

func (s *Service) AddSchedules(ctx context.Context, schedule domain.Schedule) (int64, error) {

	times, err := calculatingTimeTakingPills(time.Now(), schedule.DurationOfReceptions)
	if err != nil {
		return 0, fmt.Errorf("failed calculate schedules: %w", err)
	}

	id, err := s.ss.InsertSchedule(ctx, schedule, times)
	if err != nil {
		return 0, fmt.Errorf("failed creating schedules: %w", err)
	}
	return id, nil
}

func (s *Service) SchedulesByUserID(ctx context.Context, userID int64) ([]int64, error) {
	ids, err := s.ss.AllShedulesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed getting ids schedule: %w", err)
	}
	return ids, nil
}

func (s *Service) Schedule(ctx context.Context, userID int64, scheduleID int64) (domain.Schedule, error) {
	return s.ss.ScheduleByUserIDAndScheduleID(ctx, userID, scheduleID)
}

func (s *Service) NextTaking(ctx context.Context, userID int64) ([]string, error) {
	d := s.cfg.PeriodTime
	return s.ss.NextTakingByUserID(ctx, userID, d)
}

// CalculatingTimeTakingPills рассчитывает рассписание приёма таблеток с 8:00 - 22:00
func calculatingTimeTakingPills(createdTime time.Time, countTakings uint8) ([]string, error) {

	if countTakings == 0 {
		return nil, fmt.Errorf("count takings = 0")
	}
	if countTakings > maxCountTakings {
		return nil, fmt.Errorf("count takings > %d", maxCountTakings)
	}

	now := createdTime.UTC()
	now = truncate(now)

	times := make([]time.Time, 0, countTakings)
	times = append(times, now)
	d := maxMedicationTakingHours / time.Duration(countTakings)

	for i := uint8(1); i < countTakings; i++ {
		t := times[i-1].Add(d)

		ed := time.Date(t.Year(), t.Month(), t.Day(), 22, 0, 0, 0, time.UTC)
		sd := time.Date(t.Year(), t.Month(), t.Day(), 8, 0, 0, 0, time.UTC)

		if t.After(ed) || !t.After(sd) {
			t = t.Add(24*time.Hour - maxMedicationTakingHours)
		}
		if !(t.Minute()%15 == 0) {
			t = truncate(t)
		}
		times = append(times, t)
	}

	res := make([]string, 0, countTakings)
	for _, t := range times {
		res = append(res, t.Format(time.TimeOnly))
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})

	return res, nil
}

func truncate(t time.Time) time.Time {
	return t.Truncate(15 * time.Minute).Add(15 * time.Minute)
}
