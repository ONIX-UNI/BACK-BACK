package service

import (
	"context"
	"fmt"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/DuvanRozoParra/sicou/internal/case/models"
)

type SServiceTypeService struct {
	repo models.IServiceTypeRepository
}

func NewServiceTypeService(repo models.IServiceTypeRepository) *SServiceTypeService {
	return &SServiceTypeService{repo: repo}
}

func (s *SServiceTypeService) Create(ctx context.Context, req dto.ServiceTypeCreateRequest) (*dto.ServiceType, error) {
	if req.Code == "" {
		return nil, fmt.Errorf("code is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	return s.repo.Create(ctx, req)
}

func (s *SServiceTypeService) Update(ctx context.Context, id int16, req dto.ServiceTypeUpdateRequest) (*dto.ServiceType, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}
	if req.Name == nil && req.IsActive == nil {
		return nil, fmt.Errorf("at least one field must be provided for update")
	}
	if req.Name != nil && *req.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	return s.repo.Update(ctx, id, req)
}

func (s *SServiceTypeService) List(ctx context.Context, limit, offset int) ([]dto.ServiceType, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	return s.repo.List(ctx, limit, offset)
}

func (s *SServiceTypeService) GetById(ctx context.Context, id int16) (*dto.ServiceType, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}

	return s.repo.GetById(ctx, id)
}

func (s *SServiceTypeService) Delete(ctx context.Context, id int16) (*dto.ServiceType, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}

	return s.repo.Delete(ctx, id)
}
