package models

import "time"

type Note struct {
	ID          string    `json:"id" db:"id"`
	ClientID    string    `json:"client_id" db:"client_id"`
	AuthorID    string    `json:"author_id" db:"author_id"`
	Content     string    `json:"content" db:"content"`
	IsImportant bool      `json:"is_important" db:"is_important"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
