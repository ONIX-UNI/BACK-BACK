package repository

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SServiceType struct {
	db *pgxpool.Pool
}

func NewServiceTypeRepository(db *pgxpool.Pool) *SServiceType {
	return &SServiceType{db: db}
}

func (r *SServiceType) Create(ctx context.Context, req dto.ServiceTypeCreateRequest) (*dto.ServiceType, error) {
	query := `
		INSERT INTO sicou.service_type (code, name, is_active)
		VALUES ($1, $2, COALESCE($3, true))
		RETURNING id, code, name, is_active, created_at, updated_at;
	`

	var result dto.ServiceType

	err := r.db.QueryRow(ctx, query,
		req.Code,
		req.Name,
		req.IsActive,
	).Scan(
		&result.ID,
		&result.Code,
		&result.Name,
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SServiceType) Update(ctx context.Context, id int16, req dto.ServiceTypeUpdateRequest) (*dto.ServiceType, error) {
	query := `
		UPDATE sicou.service_type
		SET
			name = COALESCE($2, name),
			is_active = COALESCE($3, is_active),
			updated_at = now()
		WHERE id = $1
		RETURNING id, code, name, is_active, created_at, updated_at;
	`

	var result dto.ServiceType

	err := r.db.QueryRow(ctx, query,
		id,
		req.Name,
		req.IsActive,
	).Scan(
		&result.ID,
		&result.Code,
		&result.Name,
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SServiceType) List(ctx context.Context, limit, offset int) ([]dto.ServiceType, error) {
	query := `
		SELECT id, code, name, is_active, created_at, updated_at
		FROM sicou.service_type
		ORDER BY id
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.ServiceType

	for rows.Next() {
		var item dto.ServiceType
		if err := rows.Scan(
			&item.ID,
			&item.Code,
			&item.Name,
			&item.IsActive,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *SServiceType) GetById(ctx context.Context, id int16) (*dto.ServiceType, error) {
	query := `
		SELECT id, code, name, is_active, created_at, updated_at
		FROM sicou.service_type
		WHERE id = $1;
	`

	var result dto.ServiceType

	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.Code,
		&result.Name,
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SServiceType) Delete(ctx context.Context, id int16) (*dto.ServiceType, error) {
	query := `
		DELETE FROM sicou.service_type
		WHERE id = $1
		RETURNING id, code, name, is_active, created_at, updated_at;
	`

	var result dto.ServiceType

	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.Code,
		&result.Name,
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
