package models

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
)

type IDesktopRepository interface {
	Create(ctx context.Context, req dto.CreateDesktopRequest) (*dto.CreateDesktopResponse, error)
	List(ctx context.Context) ([]dto.Desktop, error)

	UpdateByCode(ctx context.Context, code string, req dto.UpdateDesktopRequest) (*dto.UpdateDesktopResponse, error)

	DeleteByCode(ctx context.Context, code string) error
}