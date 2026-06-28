package repository

import (
	"context"
	"trainers-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkoutRepo struct {
	pool *pgxpool.Pool
}

func NewWorkoutRepo(pool *pgxpool.Pool) *WorkoutRepo {
	return &WorkoutRepo{pool: pool}
}

func (r *WorkoutRepo) GetAll(ctx context.Context, userID, role string, limit, offset int) ([]models.Workout, int, error) {
	var workouts []models.Workout
	var total int

	var query string
	var args []interface{}

	if role == "trainer" {
		countQuery := `SELECT COUNT(*) FROM workouts WHERE trainer_id = $1`
		r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
		query = `SELECT id, client_id, trainer_id, title, description, workout_type, scheduled_at, 
		         duration_minutes, status, exercises, notes, created_at, updated_at 
		         FROM workouts WHERE trainer_id = $1 
		         ORDER BY scheduled_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{userID, limit, offset}
	} else {
		countQuery := `SELECT COUNT(*) FROM workouts WHERE client_id = $1`
		r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
		query = `SELECT id, client_id, trainer_id, title, description, workout_type, scheduled_at, 
		         duration_minutes, status, exercises, notes, created_at, updated_at 
		         FROM workouts WHERE client_id = $1 
		         ORDER BY scheduled_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{userID, limit, offset}
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var w models.Workout
		if err := rows.Scan(&w.ID, &w.ClientID, &w.TrainerID, &w.Title, &w.Description,
			&w.WorkoutType, &w.ScheduledAt, &w.DurationMinutes, &w.Status,
			&w.Exercises, &w.Notes, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, 0, err
		}
		workouts = append(workouts, w)
	}
	return workouts, total, nil
}

func (r *WorkoutRepo) GetByID(ctx context.Context, id string) (*models.Workout, error) {
	query := `SELECT id, client_id, trainer_id, title, description, workout_type, scheduled_at, 
	          duration_minutes, status, exercises, notes, created_at, updated_at 
	          FROM workouts WHERE id = $1`
	var w models.Workout
	err := r.pool.QueryRow(ctx, query, id).Scan(&w.ID, &w.ClientID, &w.TrainerID, &w.Title, &w.Description,
		&w.WorkoutType, &w.ScheduledAt, &w.DurationMinutes, &w.Status,
		&w.Exercises, &w.Notes, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WorkoutRepo) Create(ctx context.Context, w *models.Workout) error {
	query := `INSERT INTO workouts (client_id, trainer_id, title, description, workout_type, scheduled_at, 
	          duration_minutes, exercises, notes) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
	          RETURNING id, status, created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		w.ClientID, w.TrainerID, w.Title, w.Description, w.WorkoutType,
		w.ScheduledAt, w.DurationMinutes, w.Exercises, w.Notes,
	).Scan(&w.ID, &w.Status, &w.CreatedAt, &w.UpdatedAt)
}

func (r *WorkoutRepo) Update(ctx context.Context, w *models.Workout) error {
	query := `UPDATE workouts SET title=$1, description=$2, workout_type=$3, scheduled_at=$4, 
	          duration_minutes=$5, status=$6, exercises=$7, notes=$8, updated_at=NOW() WHERE id=$9`
	_, err := r.pool.Exec(ctx, query,
		w.Title, w.Description, w.WorkoutType, w.ScheduledAt,
		w.DurationMinutes, w.Status, w.Exercises, w.Notes, w.ID,
	)
	return err
}

func (r *WorkoutRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM workouts WHERE id = $1`, id)
	return err
}
