package service

import (
	"context"
	"errors"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/DuvanRozoParra/sicou/internal/appointments/models"
	"github.com/google/uuid"
)

type STurnAuditService struct {
	repo models.ITurnAuditRepository
}

func NewTurnAuditService(repo models.ITurnAuditRepository) *STurnAuditService {
	return &STurnAuditService{repo: repo}
}

func (s *STurnAuditService) GetTimeLine(ctx context.Context, id uuid.UUID) ([]dto.TurnAudit, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	timeline, err := s.repo.GetTimeLine(ctx, id)
	if err != nil {
		return nil, err
	}

	return timeline, nil
}
