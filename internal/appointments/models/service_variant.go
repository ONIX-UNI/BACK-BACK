package models

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/google/uuid"
)

type IServiceVariantRepository interface {
	Create(ctx context.Context, req dto.ServiceVariantCreateRequest) (*dto.ServiceVariant, error)
	Update(ctx context.Context, id uuid.UUID, req dto.ServiceVariantUpdateRequest) (*dto.ServiceVariant, error)
	List(ctx context.Context, limit, offset int) ([]dto.ServiceVariant, error)
	GetById(ctx context.Context, id uuid.UUID) (*dto.ServiceVariant, error)
	Delete(ctx context.Context, id uuid.UUID) (*dto.ServiceVariant, error)
}
