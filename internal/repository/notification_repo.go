package repository

import (
	"context"
	"trainers-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepo struct {
	pool *pgxpool.Pool
}

func NewNotificationRepo(pool *pgxpool.Pool) *NotificationRepo {
	return &NotificationRepo{pool: pool}
}

func (r *NotificationRepo) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Notification, int, error) {
	var notifications []models.Notification
	var total int

	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)

	query := `SELECT id, user_id, title, message, type, is_read, created_at 
	          FROM notifications WHERE user_id = $1 
	          ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, n)
	}
	return notifications, total, nil
}

func (r *NotificationRepo) GetByID(ctx context.Context, id string) (*models.Notification, error) {
	query := `SELECT id, user_id, title, message, type, is_read, created_at 
	          FROM notifications WHERE id = $1`
	var n models.Notification
	err := r.pool.QueryRow(ctx, query, id).Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *NotificationRepo) MarkAsRead(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE notifications SET is_read = true WHERE id = $1`, id)
	return err
}

func (r *NotificationRepo) Create(ctx context.Context, n *models.Notification) error {
	query := `INSERT INTO notifications (user_id, title, message, type) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, n.UserID, n.Title, n.Message, n.Type).
		Scan(&n.ID, &n.CreatedAt)
}
