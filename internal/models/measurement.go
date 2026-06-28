package models

import "time"

type Measurement struct {
	ID                string    `json:"id" db:"id"`
	ClientID          string    `json:"client_id" db:"client_id"`
	MeasurementDate   string    `json:"measurement_date" db:"measurement_date"`
	Weight            float64   `json:"weight,omitempty" db:"weight"`
	BodyFatPercentage float64   `json:"body_fat_percentage,omitempty" db:"body_fat_percentage"`
	MuscleMass        float64   `json:"muscle_mass,omitempty" db:"muscle_mass"`
	ChestCm           float64   `json:"chest_cm,omitempty" db:"chest_cm"`
	WaistCm           float64   `json:"waist_cm,omitempty" db:"waist_cm"`
	HipsCm            float64   `json:"hips_cm,omitempty" db:"hips_cm"`
	BicepCm           float64   `json:"bicep_cm,omitempty" db:"bicep_cm"`
	ThighCm           float64   `json:"thigh_cm,omitempty" db:"thigh_cm"`
	Notes             string    `json:"notes,omitempty" db:"notes"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}
