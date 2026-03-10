package repository

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SServiceModalityRepository struct {
	db *pgxpool.Pool
}

func NewServiceModalityRepository(db *pgxpool.Pool) *SServiceModalityRepository {
	return &SServiceModalityRepository{db: db}
}

func (r *SServiceModalityRepository) Create(ctx context.Context, req dto.ServiceModalityCreateRequest) (*dto.ServiceModality, error) {
	query := `
		INSERT INTO sicou.service_modality (
			code,
			name,
			is_active
		)
		VALUES ($1, $2, COALESCE($3, true))
		RETURNING id, code, name, is_active;
	`

	var result dto.ServiceModality

	err := r.db.QueryRow(ctx, query,
		req.Code,
		req.Name,
		req.IsActive,
	).Scan(
		&result.ID,
		&result.Code,
		&result.Name,
		&result.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SServiceModalityRepository) Update(ctx context.Context, id int16, req dto.ServiceModalityUpdateRequest) (*dto.ServiceModality, error) {
	query := `
		UPDATE sicou.service_modality
		SET
			name = COALESCE($2, name),
			is_active = COALESCE($3, is_active)
		WHERE id = $1
		RETURNING id, code, name, is_active;
	`

	var result dto.ServiceModality

	err := r.db.QueryRow(ctx, query,
		id,
		req.Name,
		req.IsActive,
	).Scan(
		&result.ID,
		&result.Code,
		&result.Name,
		&result.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SServiceModalityRepository) List(ctx context.Context, limit, offset int) ([]dto.ServiceModality, error) {
	query := `
		SELECT id, code, name, is_active
		FROM sicou.service_modality
		ORDER BY id
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.ServiceModality

	for rows.Next() {
		var item dto.ServiceModality

		if err := rows.Scan(
			&item.ID,
			&item.Code,
			&item.Name,
			&item.IsActive,
		); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, rows.Err()
}

func (r *SServiceModalityRepository) GetById(ctx context.Context, id int16) (*dto.ServiceModality, error) {
	query := `
		SELECT id, code, name, is_active
		FROM sicou.service_modality
		WHERE id = $1;
	`

	var result dto.ServiceModality

	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.Code,
		&result.Name,
		&result.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SServiceModalityRepository) Delete(ctx context.Context, id int16) (*dto.ServiceModality, error) {
	query := `
		DELETE FROM sicou.service_modality
		WHERE id = $1
		RETURNING id, code, name, is_active;
	`

	var result dto.ServiceModality

	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.Code,
		&result.Name,
		&result.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
