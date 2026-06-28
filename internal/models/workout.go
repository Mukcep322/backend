package models

import "time"

type Workout struct {
	ID              string    `json:"id" db:"id"`
	ClientID        string    `json:"client_id" db:"client_id"`
	TrainerID       string    `json:"trainer_id" db:"trainer_id"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description,omitempty" db:"description"`
	WorkoutType     string    `json:"workout_type" db:"workout_type"`
	ScheduledAt     time.Time `json:"scheduled_at" db:"scheduled_at"`
	DurationMinutes int       `json:"duration_minutes" db:"duration_minutes"`
	Status          string    `json:"status" db:"status"`
	Exercises       string    `json:"exercises,omitempty" db:"exercises"`
	Notes           string    `json:"notes,omitempty" db:"notes"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
