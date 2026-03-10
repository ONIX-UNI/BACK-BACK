package repository

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SServiceVariantRepository struct {
	db *pgxpool.Pool
}

func NewServiceVariantRepository(db *pgxpool.Pool) *SServiceVariantRepository {
	return &SServiceVariantRepository{db: db}
}

func (r *SServiceVariantRepository) Create(ctx context.Context, req dto.ServiceVariantCreateRequest) (*dto.ServiceVariant, error) {
	query := `
		INSERT INTO sicou.service_variant (
			service_type_id,
			modality_id,
			code,
			name,
			is_active
		)
		VALUES ($1,$2,$3,$4,COALESCE($5,true))
		RETURNING id, service_type_id, modality_id, code, name, is_active, created_at, updated_at;
	`

	var result dto.ServiceVariant

	err := r.db.QueryRow(ctx, query,
		req.ServiceTypeID,
		req.ModalityID,
		req.Code,
		req.Name,
		req.IsActive,
	).Scan(
		&result.ID,
		&result.ServiceTypeID,
		&result.ModalityID,
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

func (r *SServiceVariantRepository) Update(ctx context.Context, id int16, req dto.ServiceVariantUpdateRequest) (*dto.ServiceVariant, error) {
	query := `
		UPDATE sicou.service_variant
		SET
			name = COALESCE($2, name),
			is_active = COALESCE($3, is_active),
			updated_at = now()
		WHERE id = $1
		RETURNING id, service_type_id, modality_id, code, name, is_active, created_at, updated_at;
	`

	var result dto.ServiceVariant

	err := r.db.QueryRow(ctx, query,
		id,
		req.Name,
		req.IsActive,
	).Scan(
		&result.ID,
		&result.ServiceTypeID,
		&result.ModalityID,
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

func (r *SServiceVariantRepository) List(ctx context.Context, limit, offset int) ([]dto.ServiceVariant, error) {
	query := `
		SELECT id, service_type_id, modality_id, code, name, is_active, created_at, updated_at
		FROM sicou.service_variant
		ORDER BY id
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.ServiceVariant

	for rows.Next() {
		var item dto.ServiceVariant

		if err := rows.Scan(
			&item.ID,
			&item.ServiceTypeID,
			&item.ModalityID,
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

	return result, rows.Err()
}

func (r *SServiceVariantRepository) GetById(ctx context.Context, id int16) (*dto.ServiceVariant, error) {
	query := `
		SELECT id, service_type_id, modality_id, code, name, is_active, created_at, updated_at
		FROM sicou.service_variant
		WHERE id = $1;
	`

	var result dto.ServiceVariant

	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.ServiceTypeID,
		&result.ModalityID,
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

func (r *SServiceVariantRepository) Delete(ctx context.Context, id int16) (*dto.ServiceVariant, error) {
	query := `
		DELETE FROM sicou.service_variant
		WHERE id = $1
		RETURNING id, service_type_id, modality_id, code, name, is_active, created_at, updated_at;
	`

	var result dto.ServiceVariant

	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.ServiceTypeID,
		&result.ModalityID,
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
