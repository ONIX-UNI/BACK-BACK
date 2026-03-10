package models

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/google/uuid"
)

type ITurnRepository interface {
	Create(ctx context.Context, req dto.CreateTurnRequest) (*dto.Turn, error)
	List(ctx context.Context) ([]dto.Turn, error)
	GetById(ctx context.Context, id uuid.UUID) (*dto.Turn, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateTurnRequest) (*dto.Turn, error)
	Delete(ctx context.Context, id uuid.UUID) (*dto.Turn, error)

	SetTurnDesktop(ctx context.Context, req dto.SetTurnDesktopRequest) (*dto.SetTurnEscritorioResponse, error)
}

type ITurnAuditRepository interface {
	GetTimeLine(ctx context.Context, id uuid.UUID) ([]dto.TurnAudit, error)
}
