package models

import "time"

type Schedule struct {
	ID          string    `json:"id" db:"id"`
	TrainerID   string    `json:"trainer_id" db:"trainer_id"`
	DayOfWeek   int       `json:"day_of_week" db:"day_of_week"`
	StartTime   string    `json:"start_time" db:"start_time"`
	EndTime     string    `json:"end_time" db:"end_time"`
	IsAvailable bool      `json:"is_available" db:"is_available"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
