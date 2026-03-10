package service

import (
	"context"
	"errors"

	"github.com/DuvanRozoParra/sicou/internal/citizens/dto"
	"github.com/DuvanRozoParra/sicou/internal/citizens/models"
	"github.com/google/uuid"
)

type SCitizenService struct {
	repo models.ICitizenRepository
}

func NewCitizenService(repo models.ICitizenRepository) *SCitizenService {
	return &SCitizenService{repo: repo}
}

func (s *SCitizenService) Create(ctx context.Context, req dto.CreateCitizenRequest) (*dto.Citizen, error) {
	if req.DocumentTypeID == 0 {
		return nil, errors.New("document_type_id is required")
	}

	if req.DocumentNumber == "" {
		return nil, errors.New("document_number is required")
	}

	if req.FullName == "" {
		return nil, errors.New("full_name is required")
	}

	return s.repo.Create(ctx, req)
}
func (s *SCitizenService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCitizenRequest) (*dto.Citizen, error) {
	existing, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("citizen not found")
	}

	return s.repo.Update(ctx, id, req)
}
func (s *SCitizenService) List(ctx context.Context, limit, offset int) ([]dto.Citizen, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, limit, offset)
}
func (s *SCitizenService) GetById(ctx context.Context, id uuid.UUID) (*dto.Citizen, error) {

	citizen, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if citizen == nil {
		return nil, errors.New("citizen not found")
	}

	return citizen, nil
}
func (s *SCitizenService) Delete(ctx context.Context, id uuid.UUID) (*dto.Citizen, error) {
	existing, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("citizen not found")
	}

	return s.repo.Delete(ctx, id)
}
