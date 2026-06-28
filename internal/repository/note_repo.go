package repository

import (
	"context"
	"trainers-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteRepo struct {
	pool *pgxpool.Pool
}

func NewNoteRepo(pool *pgxpool.Pool) *NoteRepo {
	return &NoteRepo{pool: pool}
}

func (r *NoteRepo) GetByClientID(ctx context.Context, clientID string, limit, offset int) ([]models.Note, int, error) {
	var notes []models.Note
	var total int

	countQuery := `SELECT COUNT(*) FROM notes WHERE client_id = $1`
	err := r.pool.QueryRow(ctx, countQuery, clientID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, client_id, author_id, content, is_important, created_at, updated_at 
	          FROM notes WHERE client_id = $1 
	          ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, clientID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.ClientID, &n.AuthorID, &n.Content, &n.IsImportant, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, 0, err
		}
		notes = append(notes, n)
	}
	return notes, total, nil
}

func (r *NoteRepo) GetByID(ctx context.Context, id string) (*models.Note, error) {
	query := `SELECT id, client_id, author_id, content, is_important, created_at, updated_at 
	          FROM notes WHERE id = $1`
	var n models.Note
	err := r.pool.QueryRow(ctx, query, id).Scan(&n.ID, &n.ClientID, &n.AuthorID, &n.Content, &n.IsImportant, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *NoteRepo) Create(ctx context.Context, note *models.Note) error {
	query := `INSERT INTO notes (client_id, author_id, content, is_important) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query, note.ClientID, note.AuthorID, note.Content, note.IsImportant).
		Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
}

func (r *NoteRepo) Update(ctx context.Context, note *models.Note) error {
	query := `UPDATE notes SET content=$1, is_important=$2, updated_at=NOW() WHERE id=$3`
	_, err := r.pool.Exec(ctx, query, note.Content, note.IsImportant, note.ID)
	return err
}

func (r *NoteRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM notes WHERE id = $1`, id)
	return err
}
