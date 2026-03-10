package service

import (
	"context"
	"errors"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/pqrs/repository"
)

func (s *Service) List(ctx context.Context, in ListPQRSInput) (ListPQRSResult, error) {
	if s.repo == nil {
		return ListPQRSResult{}, errors.New("pqrs repository is not initialized")
	}

	page := in.Page
	if page < 1 {
		page = 1
	}

	limit := in.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	statuses := normalizeStatusFilter(in.Estado)
	repoResult, err := s.repo.List(ctx, repository.ListInput{
		Page:     page,
		Limit:    limit,
		Search:   strings.TrimSpace(in.Search),
		Statuses: statuses,
		Tipo:     strings.ToLower(strings.TrimSpace(in.Tipo)),
	})
	if err != nil {
		return ListPQRSResult{}, err
	}

	items := make([]PQRSListItem, 0, len(repoResult.Items))
	for _, row := range repoResult.Items {
		fechaIngreso := row.ReceivedAt.UTC()
		fechaLimite := computeDeadline(fechaIngreso)

		items = append(items, PQRSListItem{
			ID:           row.ID,
			Radicado:     row.Radicado,
			Tipo:         strings.ToLower(strings.TrimSpace(row.Tipo)),
			Asunto:       strings.TrimSpace(row.Asunto),
			Ciudadano:    strings.TrimSpace(row.Ciudadano),
			Correo:       strings.TrimSpace(row.Correo),
			Estado:       mapPublicListStatus(row.EstadoDB),
			Responsable:  strings.TrimSpace(row.Responsable),
			FechaIngreso: fechaIngreso,
			FechaLimite:  fechaLimite,
		})
	}

	return ListPQRSResult{
		Items: items,
		Total: repoResult.Total,
		Page:  repoResult.Page,
		Limit: repoResult.Limit,
	}, nil
}
