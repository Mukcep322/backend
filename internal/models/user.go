package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID         string      `json:"id" db:"id"`
	TelegramID int64       `json:"telegram_id" db:"telegram_id"`
	Username   pgtype.Text `json:"username,omitempty" db:"username"`
	FirstName  string      `json:"first_name" db:"first_name"`
	LastName   pgtype.Text `json:"last_name,omitempty" db:"last_name"`
	Phone      pgtype.Text `json:"phone,omitempty" db:"phone"`
	Role       string      `json:"role" db:"role"`
	IsActive   bool        `json:"is_active" db:"is_active"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at" db:"updated_at"`
}
