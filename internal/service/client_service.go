package service

import (
	"context"
	"trainers-backend/internal/models"
	"trainers-backend/internal/repository"
)

type ClientService struct {
	clientRepo *repository.ClientRepo
	userRepo   *repository.UserRepo
}

func NewClientService(clientRepo *repository.ClientRepo, userRepo *repository.UserRepo) *ClientService {
	return &ClientService{clientRepo: clientRepo, userRepo: userRepo}
}

func (s *ClientService) GetAll(ctx context.Context, trainerID string, page, pageSize int) ([]models.Client, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.clientRepo.GetAll(ctx, trainerID, pageSize, offset)
}

func (s *ClientService) GetByID(ctx context.Context, id string) (*models.Client, error) {
	return s.clientRepo.GetByID(ctx, id)
}

func (s *ClientService) Create(ctx context.Context, client *models.Client) error {
	return s.clientRepo.Create(ctx, client)
}

func (s *ClientService) Update(ctx context.Context, client *models.Client) error {
	return s.clientRepo.Update(ctx, client)
}

func (s *ClientService) Delete(ctx context.Context, id string) error {
	return s.clientRepo.Delete(ctx, id)
}
