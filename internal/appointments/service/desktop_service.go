package service

import (
	"context"
	"errors"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/DuvanRozoParra/sicou/internal/appointments/models"
)

type SDesktopService struct {
	repo models.IDesktopRepository
}

func NewDesktopService(repo models.IDesktopRepository) *SDesktopService {
	return &SDesktopService{repo: repo}
}

func (s *SDesktopService) Create(ctx context.Context, req dto.CreateDesktopRequest) (*dto.CreateDesktopResponse, error) {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)

	if req.Code == "" {
		return nil, errors.New("code is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	if req.Location != nil {
		loc := strings.TrimSpace(*req.Location)
		if loc == "" {
			req.Location = nil
		} else {
			req.Location = &loc
		}
	}

	return s.repo.Create(ctx, req)
}

func (s *SDesktopService) UpdateByCode(ctx context.Context, code string, req dto.UpdateDesktopRequest) (*dto.UpdateDesktopResponse, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, errors.New("code is required")
	}

	// ✅ valida que venga ALGO
	if req.Name == nil && req.IsActive == nil && req.UserID == nil {
		return nil, errors.New("at least one field is required: name, is_active, user_id")
	}

	// normaliza name
	if req.Name != nil {
		n := strings.TrimSpace(*req.Name)
		if n == "" {
			return nil, errors.New("name cannot be empty")
		}
		req.Name = &n
	}

	if req.Reason != nil {
		r := strings.TrimSpace(*req.Reason)
		if r == "" {
			req.Reason = nil
		} else {
			req.Reason = &r
		}
	}

	return s.repo.UpdateByCode(ctx, code, req)
}

func (s *SDesktopService) List(ctx context.Context) ([]dto.Desktop, error) {
	return s.repo.List(ctx)
}

// ✅ NUEVO
func (s *SDesktopService) DeleteByCode(ctx context.Context, code string) error {
	code = strings.TrimSpace(code)
	if code == "" {
		return errors.New("code is required")
	}
	return s.repo.DeleteByCode(ctx, code)
}
