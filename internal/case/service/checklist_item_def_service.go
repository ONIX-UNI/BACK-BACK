package service

import (
	"context"
	"fmt"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/DuvanRozoParra/sicou/internal/case/models"
	"github.com/google/uuid"
)

type SCheckListItemDefService struct {
	repo models.ICheckListItemDef
}

func NewCheckListItemDefService(repo models.ICheckListItemDef) *SCheckListItemDefService {
	return &SCheckListItemDefService{repo: repo}
}

func (s *SCheckListItemDefService) Create(ctx context.Context, req dto.ChecklistItemDefCreateRequest) (*dto.ChecklistItemDef, error) {
	if req.Code == "" {
		return nil, fmt.Errorf("code is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	return s.repo.Create(ctx, req)
}
func (s *SCheckListItemDefService) Update(ctx context.Context, id uuid.UUID, req dto.ChecklistItemDefUpdateRequest) (*dto.ChecklistItemDef, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid id")
	}

	if req.Code != nil && *req.Code == "" {
		return nil, fmt.Errorf("code cannot be empty")
	}
	if req.Name != nil && *req.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	return s.repo.Update(ctx, id, req)
}
func (s *SCheckListItemDefService) List(ctx context.Context, limit, offset int) ([]dto.ChecklistItemDef, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, limit, offset)
}
func (s *SCheckListItemDefService) GetById(ctx context.Context, id uuid.UUID) (*dto.ChecklistItemDef, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid id")
	}

	return s.repo.GetById(ctx, id)
}
func (s *SCheckListItemDefService) Delete(ctx context.Context, id uuid.UUID) (*dto.ChecklistItemDef, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid id")
	}

	return s.repo.Delete(ctx, id)
}
