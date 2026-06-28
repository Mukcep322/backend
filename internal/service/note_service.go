package service

import (
	"context"
	"trainers-backend/internal/models"
	"trainers-backend/internal/repository"
)

type NoteService struct {
	repo *repository.NoteRepo
}

func NewNoteService(repo *repository.NoteRepo) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) GetByClientID(ctx context.Context, clientID string, page, pageSize int) ([]models.Note, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.GetByClientID(ctx, clientID, pageSize, offset)
}

func (s *NoteService) Create(ctx context.Context, note *models.Note) error {
	return s.repo.Create(ctx, note)
}

func (s *NoteService) Update(ctx context.Context, note *models.Note) error {
	return s.repo.Update(ctx, note)
}

func (s *NoteService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
