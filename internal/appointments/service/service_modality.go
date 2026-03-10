package service

import (
	"context"
	"fmt"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/DuvanRozoParra/sicou/internal/appointments/models"
	"github.com/google/uuid"
)

type SServiceModalityService struct {
	repo models.IServiceModalityRepository
}

func NewServiceModality(repo models.IServiceModalityRepository) *SServiceModalityService {
	return &SServiceModalityService{repo: repo}
}

func (s *SServiceModalityService) Create(ctx context.Context, req dto.ServiceModalityCreateRequest) (*dto.ServiceModality, error) {
	if req.Code == "" || req.Name == "" {
		return nil, fmt.Errorf("code and name are required")
	}

	return s.repo.Create(ctx, req)
}

func (s *SServiceModalityService) Update(ctx context.Context, id uuid.UUID, req dto.ServiceModalityUpdateRequest) (*dto.ServiceModality, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid id")
	}

	return s.repo.Update(ctx, id, req)
}

func (s *SServiceModalityService) List(ctx context.Context, limit, offset int) ([]dto.ServiceModality, error) {
	return s.repo.List(ctx, limit, offset)
}
func (s *SServiceModalityService) GetById(ctx context.Context, id uuid.UUID) (*dto.ServiceModality, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid id")
	}

	return s.repo.GetById(ctx, id)
}
func (s *SServiceModalityService) Delete(ctx context.Context, id uuid.UUID) (*dto.ServiceModality, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid id")
	}

	return s.repo.Delete(ctx, id)
}
