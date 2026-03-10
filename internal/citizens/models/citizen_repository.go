package models

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/citizens/dto"
	"github.com/google/uuid"
)

type ICitizenRepository interface {
	Create(ctx context.Context, req dto.CreateCitizenRequest) (*dto.Citizen, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateCitizenRequest) (*dto.Citizen, error)
	List(ctx context.Context, limit, offset int) ([]dto.Citizen, error)
	GetById(ctx context.Context, id uuid.UUID) (*dto.Citizen, error)
	Delete(ctx context.Context, id uuid.UUID) (*dto.Citizen, error)
}
