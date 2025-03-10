package domain

import (
	"time"
)

type Schedule struct {
	ID           int64
	UserID       int64
	MedicineName string
	//От 1 до 14 приёмов в день
	DurationOfReceptions uint8
	//От 1 дня до бесконечности
	DurationOfTreatment time.Time
	//Расписание приёма
	TimesTakingPills []time.Time
}
