package service

import (
	"context"
	"trainers-backend/internal/models"
	"trainers-backend/internal/repository"
)

type MeasurementService struct {
	repo *repository.MeasurementRepo
}

func NewMeasurementService(repo *repository.MeasurementRepo) *MeasurementService {
	return &MeasurementService{repo: repo}
}

func (s *MeasurementService) GetByClientID(ctx context.Context, clientID string, page, pageSize int) ([]models.Measurement, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.GetByClientID(ctx, clientID, pageSize, offset)
}

func (s *MeasurementService) Create(ctx context.Context, m *models.Measurement) error {
	return s.repo.Create(ctx, m)
}

func (s *MeasurementService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
