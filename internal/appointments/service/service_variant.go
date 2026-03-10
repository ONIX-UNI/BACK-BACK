package service

import (
	"context"
	"fmt"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/DuvanRozoParra/sicou/internal/appointments/models"
	"github.com/google/uuid"
)

type SServiceVariantService struct {
	repo models.IServiceVariantRepository
}

func NewServiceVariant(repo models.IServiceVariantRepository) *SServiceVariantService {
	return &SServiceVariantService{repo: repo}
}

func (s *SServiceVariantService) Create(ctx context.Context, req dto.ServiceVariantCreateRequest) (*dto.ServiceVariant, error) {
	if req.ServiceTypeID <= 0 {
		return nil, fmt.Errorf("service_type_id is required")
	}

	if req.ModalityID <= 0 {
		return nil, fmt.Errorf("modality_id is required")
	}

	if req.Code == "" || req.Name == "" {
		return nil, fmt.Errorf("code and name are required")
	}

	return s.repo.Create(ctx, req)
}

func (s *SServiceVariantService) Update(ctx context.Context, id uuid.UUID, req dto.ServiceVariantUpdateRequest) (*dto.ServiceVariant, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid id")
	}

	return s.repo.Update(ctx, id, req)
}

func (s *SServiceVariantService) List(ctx context.Context, limit, offset int) ([]dto.ServiceVariant, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *SServiceVariantService) GetById(ctx context.Context, id uuid.UUID) (*dto.ServiceVariant, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid id")
	}

	return s.repo.GetById(ctx, id)
}

func (s *SServiceVariantService) Delete(ctx context.Context, id uuid.UUID) (*dto.ServiceVariant, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid id")
	}

	return s.repo.Delete(ctx, id)
}
