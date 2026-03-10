package service

import (
	"context"
	"errors"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/DuvanRozoParra/sicou/internal/case/models"
	"github.com/google/uuid"
)

type SCaseEventService struct {
	repo models.ICaseEventRepository
}

func NewCaseEventService(repo models.ICaseEventRepository) *SCaseEventService {
	return &SCaseEventService{repo: repo}
}

func (s *SCaseEventService) Create(ctx context.Context, req dto.CreateCaseEventRequest) (*dto.CaseEvent, error) {

	if req.CaseID == uuid.Nil {
		return nil, errors.New("case_id is required")
	}

	if strings.TrimSpace(req.EventType) == "" {
		return nil, errors.New("event_type is required")
	}

	return s.repo.Create(ctx, req)
}
func (s *SCaseEventService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCaseEventRequest) (*dto.CaseEvent, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid id")
	}

	// Verificar existencia
	existing, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("case_event not found")
	}

	return s.repo.Update(ctx, id, req)
}
func (s *SCaseEventService) List(ctx context.Context, limit, offset int) ([]dto.CaseEvent, error) {

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
func (s *SCaseEventService) GetById(ctx context.Context, id uuid.UUID) (*dto.CaseEvent, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid id")
	}

	event, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, errors.New("case_event not found")
	}

	return event, nil
}
func (s *SCaseEventService) Delete(ctx context.Context, id uuid.UUID) (*dto.CaseEvent, error) {

	if id == uuid.Nil {
		return nil, errors.New("invalid id")
	}

	// Verificar existencia antes de borrar
	existing, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("case_event not found")
	}

	return s.repo.Delete(ctx, id)
}
