package repository

import (
	"context"
	"trainers-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepo struct {
	pool *pgxpool.Pool
}

func NewScheduleRepo(pool *pgxpool.Pool) *ScheduleRepo {
	return &ScheduleRepo{pool: pool}
}

func (r *ScheduleRepo) GetAll(ctx context.Context, trainerID string) ([]models.Schedule, error) {
	var schedules []models.Schedule
	query := `SELECT id, trainer_id, day_of_week, start_time, end_time, is_available, created_at 
	          FROM schedule WHERE trainer_id = $1 ORDER BY day_of_week, start_time`
	rows, err := r.pool.Query(ctx, query, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.TrainerID, &s.DayOfWeek, &s.StartTime, &s.EndTime, &s.IsAvailable, &s.CreatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *ScheduleRepo) GetByID(ctx context.Context, id string) (*models.Schedule, error) {
	query := `SELECT id, trainer_id, day_of_week, start_time, end_time, is_available, created_at 
	          FROM schedule WHERE id = $1`
	var s models.Schedule
	err := r.pool.QueryRow(ctx, query, id).Scan(&s.ID, &s.TrainerID, &s.DayOfWeek, &s.StartTime, &s.EndTime, &s.IsAvailable, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ScheduleRepo) Create(ctx context.Context, s *models.Schedule) error {
	query := `INSERT INTO schedule (trainer_id, day_of_week, start_time, end_time, is_available) 
	          VALUES ($1, $2, $3, $4, $5) 
	          RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, s.TrainerID, s.DayOfWeek, s.StartTime, s.EndTime, s.IsAvailable).
		Scan(&s.ID, &s.CreatedAt)
}

func (r *ScheduleRepo) Update(ctx context.Context, s *models.Schedule) error {
	query := `UPDATE schedule SET day_of_week=$1, start_time=$2, end_time=$3, is_available=$4 
	          WHERE id=$5`
	_, err := r.pool.Exec(ctx, query, s.DayOfWeek, s.StartTime, s.EndTime, s.IsAvailable, s.ID)
	return err
}

func (r *ScheduleRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM schedule WHERE id = $1`, id)
	return err
}
