package dto

type TelegramAuthRequest struct {
	InitData string `json:"init_data" validate:"required"`
}

type UpdateClientRequest struct {
	DateOfBirth       string  `json:"date_of_birth,omitempty"`
	Gender            string  `json:"gender,omitempty"`
	HeightCm          float64 `json:"height_cm,omitempty"`
	WeightKg          float64 `json:"weight_kg,omitempty"`
	Goal              string  `json:"goal,omitempty"`
	MedicalConditions string  `json:"medical_conditions,omitempty"`
}

type CreateMeasurementRequest struct {
	MeasurementDate   string  `json:"measurement_date" validate:"required"`
	Weight            float64 `json:"weight,omitempty"`
	BodyFatPercentage float64 `json:"body_fat_percentage,omitempty"`
	MuscleMass        float64 `json:"muscle_mass,omitempty"`
	ChestCm           float64 `json:"chest_cm,omitempty"`
	WaistCm           float64 `json:"waist_cm,omitempty"`
	HipsCm            float64 `json:"hips_cm,omitempty"`
	BicepCm           float64 `json:"bicep_cm,omitempty"`
	ThighCm           float64 `json:"thigh_cm,omitempty"`
	Notes             string  `json:"notes,omitempty"`
}

type CreateNoteRequest struct {
	Content     string `json:"content" validate:"required"`
	IsImportant bool   `json:"is_important"`
}

type UpdateNoteRequest struct {
	Content     string `json:"content,omitempty"`
	IsImportant bool   `json:"is_important,omitempty"`
}

type CreateWorkoutRequest struct {
	ClientID        string `json:"client_id" validate:"required"`
	Title           string `json:"title" validate:"required"`
	Description     string `json:"description,omitempty"`
	WorkoutType     string `json:"workout_type" validate:"required"`
	ScheduledAt     string `json:"scheduled_at" validate:"required"`
	DurationMinutes int    `json:"duration_minutes" validate:"required"`
	Exercises       string `json:"exercises,omitempty"`
	Notes           string `json:"notes,omitempty"`
}

type UpdateWorkoutRequest struct {
	Title           string `json:"title,omitempty"`
	Description     string `json:"description,omitempty"`
	WorkoutType     string `json:"workout_type,omitempty"`
	ScheduledAt     string `json:"scheduled_at,omitempty"`
	DurationMinutes int    `json:"duration_minutes,omitempty"`
	Status          string `json:"status,omitempty"`
	Exercises       string `json:"exercises,omitempty"`
	Notes           string `json:"notes,omitempty"`
}

type CreateScheduleRequest struct {
	DayOfWeek   int    `json:"day_of_week" validate:"required"`
	StartTime   string `json:"start_time" validate:"required"`
	EndTime     string `json:"end_time" validate:"required"`
	IsAvailable bool   `json:"is_available"`
}

type UpdateScheduleRequest struct {
	DayOfWeek   int    `json:"day_of_week,omitempty"`
	StartTime   string `json:"start_time,omitempty"`
	EndTime     string `json:"end_time,omitempty"`
	IsAvailable bool   `json:"is_available,omitempty"`
}
