package repository

import (
	"context"
	"trainers-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	query := `SELECT id, telegram_id, username, first_name, last_name, phone, role, is_active, created_at, updated_at 
	          FROM users WHERE telegram_id = $1`
	var user models.User
	err := r.pool.QueryRow(ctx, query, telegramID).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName,
		&user.LastName, &user.Phone, &user.Role, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, telegram_id, username, first_name, last_name, phone, role, is_active, created_at, updated_at 
	          FROM users WHERE id = $1`
	var user models.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName,
		&user.LastName, &user.Phone, &user.Role, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (telegram_id, username, first_name, last_name, phone, role) 
	          VALUES ($1, $2, $3, $4, $5, $6) 
	          RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		user.TelegramID, user.Username, user.FirstName,
		user.LastName, user.Phone, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET username=$1, first_name=$2, last_name=$3, phone=$4, updated_at=NOW() 
	          WHERE id=$5`
	_, err := r.pool.Exec(ctx, query,
		user.Username, user.FirstName, user.LastName, user.Phone, user.ID,
	)
	return err
}
