package models

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
)

type IServiceTypeRepository interface {
	Create(ctx context.Context, req dto.ServiceTypeCreateRequest) (*dto.ServiceType, error)
	Update(ctx context.Context, id int16, req dto.ServiceTypeUpdateRequest) (*dto.ServiceType, error)
	List(ctx context.Context, limit, offset int) ([]dto.ServiceType, error)
	GetById(ctx context.Context, id int16) (*dto.ServiceType, error)
	Delete(ctx context.Context, id int16) (*dto.ServiceType, error)
}
