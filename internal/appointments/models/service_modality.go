package models

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/google/uuid"
)

type IServiceModalityRepository interface {
	Create(ctx context.Context, req dto.ServiceModalityCreateRequest) (*dto.ServiceModality, error)
	Update(ctx context.Context, id uuid.UUID, req dto.ServiceModalityUpdateRequest) (*dto.ServiceModality, error)
	List(ctx context.Context, limit, offset int) ([]dto.ServiceModality, error)
	GetById(ctx context.Context, id uuid.UUID) (*dto.ServiceModality, error)
	Delete(ctx context.Context, id uuid.UUID) (*dto.ServiceModality, error)
}
