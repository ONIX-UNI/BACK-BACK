package models

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/google/uuid"
)

type ICaseEventRepository interface {
	Create(ctx context.Context, req dto.CreateCaseEventRequest) (*dto.CaseEvent, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateCaseEventRequest) (*dto.CaseEvent, error)
	List(ctx context.Context, limit, offset int) ([]dto.CaseEvent, error)
	GetById(ctx context.Context, id uuid.UUID) (*dto.CaseEvent, error)
	Delete(ctx context.Context, id uuid.UUID) (*dto.CaseEvent, error)
}

// type IDocumentKindRepository interface {
// 	Create(ctx context.Context, req dto.CreateDocumentKind) (*dto.DocumentKind, error)
// 	Update(ctx context.Context, req dto.updateDocumentKind) (*dto.DocumentKind, error)
// 	GetByID(ctx context.Context, )
// 	GetByCode
// 	Delete
// }
