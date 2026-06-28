package models

import "time"

type Client struct {
	ID                string    `json:"id" db:"id"`
	UserID            string    `json:"user_id" db:"user_id"`
	TrainerID         string    `json:"trainer_id,omitempty" db:"trainer_id"`
	DateOfBirth       string    `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Gender            string    `json:"gender,omitempty" db:"gender"`
	HeightCm          float64   `json:"height_cm,omitempty" db:"height_cm"`
	WeightKg          float64   `json:"weight_kg,omitempty" db:"weight_kg"`
	Goal              string    `json:"goal,omitempty" db:"goal"`
	MedicalConditions string    `json:"medical_conditions,omitempty" db:"medical_conditions"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}
