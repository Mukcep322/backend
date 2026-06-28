package repository

import (
	"context"
	"trainers-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MeasurementRepo struct {
	pool *pgxpool.Pool
}

func NewMeasurementRepo(pool *pgxpool.Pool) *MeasurementRepo {
	return &MeasurementRepo{pool: pool}
}

func (r *MeasurementRepo) GetByClientID(ctx context.Context, clientID string, limit, offset int) ([]models.Measurement, int, error) {
	var measurements []models.Measurement
	var total int

	countQuery := `SELECT COUNT(*) FROM measurements WHERE client_id = $1`
	err := r.pool.QueryRow(ctx, countQuery, clientID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, client_id, measurement_date, weight, body_fat_percentage, muscle_mass, 
	          chest_cm, waist_cm, hips_cm, bicep_cm, thigh_cm, notes, created_at 
	          FROM measurements WHERE client_id = $1 
	          ORDER BY measurement_date DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, clientID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Measurement
		if err := rows.Scan(&m.ID, &m.ClientID, &m.MeasurementDate, &m.Weight,
			&m.BodyFatPercentage, &m.MuscleMass, &m.ChestCm, &m.WaistCm,
			&m.HipsCm, &m.BicepCm, &m.ThighCm, &m.Notes, &m.CreatedAt); err != nil {
			return nil, 0, err
		}
		measurements = append(measurements, m)
	}
	return measurements, total, nil
}

func (r *MeasurementRepo) GetByID(ctx context.Context, id string) (*models.Measurement, error) {
	query := `SELECT id, client_id, measurement_date, weight, body_fat_percentage, muscle_mass, 
	          chest_cm, waist_cm, hips_cm, bicep_cm, thigh_cm, notes, created_at 
	          FROM measurements WHERE id = $1`
	var m models.Measurement
	err := r.pool.QueryRow(ctx, query, id).Scan(&m.ID, &m.ClientID, &m.MeasurementDate, &m.Weight,
		&m.BodyFatPercentage, &m.MuscleMass, &m.ChestCm, &m.WaistCm,
		&m.HipsCm, &m.BicepCm, &m.ThighCm, &m.Notes, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MeasurementRepo) Create(ctx context.Context, m *models.Measurement) error {
	query := `INSERT INTO measurements (client_id, measurement_date, weight, body_fat_percentage, muscle_mass, 
	          chest_cm, waist_cm, hips_cm, bicep_cm, thigh_cm, notes) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
	          RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query,
		m.ClientID, m.MeasurementDate, m.Weight, m.BodyFatPercentage, m.MuscleMass,
		m.ChestCm, m.WaistCm, m.HipsCm, m.BicepCm, m.ThighCm, m.Notes,
	).Scan(&m.ID, &m.CreatedAt)
}

func (r *MeasurementRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM measurements WHERE id = $1`, id)
	return err
}
