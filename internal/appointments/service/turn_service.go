package service

import (
	"context"
	"errors"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/DuvanRozoParra/sicou/internal/appointments/models"
	"github.com/google/uuid"
)

type STurnService struct {
	repo models.ITurnRepository
}

func NewTurnService(repo models.ITurnRepository) *STurnService {
	return &STurnService{repo: repo}
}

func (s *STurnService) Create(ctx context.Context, req dto.CreateTurnRequest) (*dto.Turn, error) {
	if req.PreturnoID == uuid.Nil {
		return nil, errors.New("preturno_id is required")
	}

	return s.repo.Create(ctx, req)
}
func (s *STurnService) List(ctx context.Context) ([]dto.Turn, error) {
	return s.repo.List(ctx)
}
func (s *STurnService) GetById(ctx context.Context, id uuid.UUID) (*dto.Turn, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid id")
	}

	turn, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if turn == nil {
		return nil, errors.New("turn not found")
	}

	return turn, nil
}
func (s *STurnService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateTurnRequest) (*dto.Turn, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid id")
	}

	// Validación opcional de status permitido
	if req.Status != nil {
		allowed := map[string]bool{
			"EN_COLA":     true,
			"LLAMADO":     true,
			"EN_ATENCION": true,
			"FINALIZADO":  true,
			"ANULADO":     true,
		}

		status := strings.ToUpper(*req.Status)
		if !allowed[status] {
			return nil, errors.New("invalid status value")
		}

		req.Status = &status
	}

	turn, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	if turn == nil {
		return nil, errors.New("turn not found")
	}

	return turn, nil
}
func (s *STurnService) Delete(ctx context.Context, id uuid.UUID) (*dto.Turn, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid id")
	}

	turn, err := s.repo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	if turn == nil {
		return nil, errors.New("turn not found")
	}

	return turn, nil
}

func (s *STurnService) SetTurnDesktop(ctx context.Context, req dto.SetTurnDesktopRequest) (*dto.SetTurnEscritorioResponse, error) {

	if req.TurnID == uuid.Nil {
		return nil, errors.New("invalid turn_id")
	}

	if req.EscritorioID == uuid.Nil {
		return nil, errors.New("invalid escritorio_id")
	}

	// Validación opcional de reason (trim espacios)
	if req.Reason != nil {
		trimmed := strings.TrimSpace(*req.Reason)
		req.Reason = &trimmed
	}

	return s.repo.SetTurnDesktop(ctx, req)
}
