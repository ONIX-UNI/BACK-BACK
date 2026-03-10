package service

import (
	"context"
	"errors"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/DuvanRozoParra/sicou/internal/case/models"
)

type SlegalAreaService struct {
	repo models.ILegalArea
}

func NewLegalAreaService(repo models.ILegalArea) *SlegalAreaService {
	return &SlegalAreaService{repo: repo}
}

func (s *SlegalAreaService) Create(ctx context.Context, req dto.CreateLegalAreaRequest) (*dto.LegalArea, error) {

	if strings.TrimSpace(req.Code) == "" {
		return nil, errors.New("code is required")
	}

	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("name is required")
	}

	return s.repo.Create(ctx, req)
}
func (s *SlegalAreaService) Update(ctx context.Context, id int16, req dto.UpdateLegalAreaRequest) (*dto.LegalArea, error) {

	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	// verificar existencia
	_, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	// validar que al menos venga un campo
	if req.Code == nil && req.Name == nil && req.IsActive == nil {
		return nil, errors.New("at least one field must be provided")
	}

	return s.repo.Update(ctx, id, req)
}
func (s *SlegalAreaService) List(ctx context.Context, limit, offset int) ([]dto.LegalArea, error) {

	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, limit, offset)
}
func (s *SlegalAreaService) GetById(ctx context.Context, id int16) (*dto.LegalArea, error) {

	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	return s.repo.GetById(ctx, id)
}
func (s *SlegalAreaService) Delete(ctx context.Context, id int16) (*dto.LegalArea, error) {

	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	// verificar existencia antes de eliminar
	_, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.repo.Delete(ctx, id)
}
