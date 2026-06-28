package service

import (
	"context"
	"trainers-backend/internal/models"
	"trainers-backend/internal/repository"
)

type WorkoutService struct {
	repo *repository.WorkoutRepo
}

func NewWorkoutService(repo *repository.WorkoutRepo) *WorkoutService {
	return &WorkoutService{repo: repo}
}

func (s *WorkoutService) GetAll(ctx context.Context, userID, role string, page, pageSize int) ([]models.Workout, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.GetAll(ctx, userID, role, pageSize, offset)
}

func (s *WorkoutService) GetByID(ctx context.Context, id string) (*models.Workout, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *WorkoutService) Create(ctx context.Context, w *models.Workout) error {
	return s.repo.Create(ctx, w)
}

func (s *WorkoutService) Update(ctx context.Context, w *models.Workout) error {
	return s.repo.Update(ctx, w)
}

func (s *WorkoutService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
