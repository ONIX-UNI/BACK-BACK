package models

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
)

type ILegalArea interface {
	Create(ctx context.Context, req dto.CreateLegalAreaRequest) (*dto.LegalArea, error)
	Update(ctx context.Context, id int16, req dto.UpdateLegalAreaRequest) (*dto.LegalArea, error)
	List(ctx context.Context, limit, offset int) ([]dto.LegalArea, error)
	GetById(ctx context.Context, id int16) (*dto.LegalArea, error)
	Delete(ctx context.Context, id int16) (*dto.LegalArea, error)
}
