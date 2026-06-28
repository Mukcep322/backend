package service

import (
	"context"
	"trainers-backend/internal/models"
	"trainers-backend/internal/repository"
)

type ScheduleService struct {
	repo *repository.ScheduleRepo
}

func NewScheduleService(repo *repository.ScheduleRepo) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) GetAll(ctx context.Context, trainerID string) ([]models.Schedule, error) {
	return s.repo.GetAll(ctx, trainerID)
}

func (s *ScheduleService) Create(ctx context.Context, schedule *models.Schedule) error {
	return s.repo.Create(ctx, schedule)
}

func (s *ScheduleService) Update(ctx context.Context, schedule *models.Schedule) error {
	return s.repo.Update(ctx, schedule)
}

func (s *ScheduleService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ScheduleService) GetByID(ctx context.Context, id string) (*models.Schedule, error) {
	return s.repo.GetByID(ctx, id)
}
