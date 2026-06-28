package models

import "time"

type User struct {
	ID         string    `json:"id" db:"id"`
	TelegramID int64     `json:"telegram_id" db:"telegram_id"`
	Username   string    `json:"username,omitempty" db:"username"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name,omitempty" db:"last_name"`
	Phone      string    `json:"phone,omitempty" db:"phone"`
	Role       string    `json:"role" db:"role"`
	IsActive   bool      `json:"is_active" db:"is_active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
