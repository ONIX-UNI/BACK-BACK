package models

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/google/uuid"
)

type ICheckListItemDef interface {
	Create(ctx context.Context, req dto.ChecklistItemDefCreateRequest) (*dto.ChecklistItemDef, error)
	Update(ctx context.Context, id uuid.UUID, req dto.ChecklistItemDefUpdateRequest) (*dto.ChecklistItemDef, error)
	List(ctx context.Context, limit, offset int) ([]dto.ChecklistItemDef, error)
	GetById(ctx context.Context, id uuid.UUID) (*dto.ChecklistItemDef, error)
	Delete(ctx context.Context, id uuid.UUID) (*dto.ChecklistItemDef, error)
}
