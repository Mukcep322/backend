package repository

import (
	"context"
	"trainers-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepo struct {
	pool *pgxpool.Pool
}

func NewClientRepo(pool *pgxpool.Pool) *ClientRepo {
	return &ClientRepo{pool: pool}
}

func (r *ClientRepo) GetByID(ctx context.Context, id string) (*models.Client, error) {
	query := `SELECT id, user_id, trainer_id, date_of_birth, gender, height_cm, weight_kg, 
	          goal, medical_conditions, created_at, updated_at 
	          FROM clients WHERE id = $1`
	var client models.Client
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&client.ID, &client.UserID, &client.TrainerID, &client.DateOfBirth,
		&client.Gender, &client.HeightCm, &client.WeightKg,
		&client.Goal, &client.MedicalConditions,
		&client.CreatedAt, &client.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepo) GetByUserID(ctx context.Context, userID string) (*models.Client, error) {
	query := `SELECT id, user_id, trainer_id, date_of_birth, gender, height_cm, weight_kg, 
	          goal, medical_conditions, created_at, updated_at 
	          FROM clients WHERE user_id = $1`
	var client models.Client
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&client.ID, &client.UserID, &client.TrainerID, &client.DateOfBirth,
		&client.Gender, &client.HeightCm, &client.WeightKg,
		&client.Goal, &client.MedicalConditions,
		&client.CreatedAt, &client.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepo) GetAll(ctx context.Context, trainerID string, limit, offset int) ([]models.Client, int, error) {
	var clients []models.Client
	var total int

	countQuery := `SELECT COUNT(*) FROM clients WHERE trainer_id = $1`
	err := r.pool.QueryRow(ctx, countQuery, trainerID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, user_id, trainer_id, date_of_birth, gender, height_cm, weight_kg, 
	          goal, medical_conditions, created_at, updated_at 
	          FROM clients WHERE trainer_id = $1 
	          ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, trainerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var client models.Client
		if err := rows.Scan(&client.ID, &client.UserID, &client.TrainerID, &client.DateOfBirth,
			&client.Gender, &client.HeightCm, &client.WeightKg,
			&client.Goal, &client.MedicalConditions,
			&client.CreatedAt, &client.UpdatedAt); err != nil {
			return nil, 0, err
		}
		clients = append(clients, client)
	}
	return clients, total, nil
}

func (r *ClientRepo) Create(ctx context.Context, client *models.Client) error {
	query := `INSERT INTO clients (user_id, trainer_id, date_of_birth, gender, height_cm, weight_kg, goal, medical_conditions) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	          RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		client.UserID, client.TrainerID, client.DateOfBirth, client.Gender,
		client.HeightCm, client.WeightKg, client.Goal, client.MedicalConditions,
	).Scan(&client.ID, &client.CreatedAt, &client.UpdatedAt)
}

func (r *ClientRepo) Update(ctx context.Context, client *models.Client) error {
	query := `UPDATE clients SET date_of_birth=$1, gender=$2, height_cm=$3, weight_kg=$4, 
	          goal=$5, medical_conditions=$6, updated_at=NOW() WHERE id=$7`
	_, err := r.pool.Exec(ctx, query,
		client.DateOfBirth, client.Gender, client.HeightCm, client.WeightKg,
		client.Goal, client.MedicalConditions, client.ID,
	)
	return err
}

func (r *ClientRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM clients WHERE id = $1`, id)
	return err
}
