package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/LeoUraltsev/medication-reminder/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestService_CalculatingTimeTakingPills(t *testing.T) {
	type args struct {
		createdTime time.Time
		count       uint8
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"test count 3",
			args{
				createdTime: time.Date(2025, time.March, 1, 10, 13, 0, 0, time.UTC),
				count:       3,
			},
			[]string{
				"10:15:00",
				"15:00:00",
				"19:45:00",
			},
		},
		{"test count 1",
			args{
				createdTime: time.Date(2025, time.March, 1, 10, 13, 0, 0, time.UTC),
				count:       1,
			},
			[]string{
				"10:15:00",
			},
		},
		{"test count 2",
			args{
				createdTime: time.Date(2025, time.March, 1, 10, 13, 0, 0, time.UTC),
				count:       2,
			},
			[]string{
				"10:15:00",
				"17:15:00",
			},
		},
		{"test count 14",
			args{
				createdTime: time.Date(2025, time.March, 1, 10, 13, 0, 0, time.UTC),
				count:       14,
			},
			[]string{
				"08:15:00",
				"09:15:00",
				"10:15:00",
				"11:15:00",
				"12:15:00",
				"13:15:00",
				"14:15:00",
				"15:15:00",
				"16:15:00",
				"17:15:00",
				"18:15:00",
				"19:15:00",
				"20:15:00",
				"21:15:00",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := calculatingTimeTakingPills(tt.args.createdTime, tt.args.count)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_FailedCalculatingTimeTakingPills(t *testing.T) {
	type args struct {
		createdTime time.Time
		count       uint8
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"test count 15",
			args{
				createdTime: time.Date(2025, time.March, 1, 10, 13, 0, 0, time.UTC),
				count:       15,
			},
			fmt.Errorf("count takings > %d", maxCountTakings),
		},
		{"test count 255",
			args{
				createdTime: time.Date(2025, time.March, 1, 10, 13, 0, 0, time.UTC),
				count:       15,
			},
			fmt.Errorf("count takings > %d", maxCountTakings),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := calculatingTimeTakingPills(tt.args.createdTime, tt.args.count)
			assert.EqualError(t, tt.want, gotErr.Error())
		})
	}
}

func TestService_AddSchedule(t *testing.T) {
	s := Service{}
	d := domain.Schedule{
		ID:                   0,
		UserID:               0,
		MedicineName:         "",
		DurationOfReceptions: 0,
		DurationOfTreatment:  time.Time{},
		TimesTakingPills:     []string{},
	}
	s.AddSchedules(context.Background(), &d)
	t.Errorf("%v", d)
}
