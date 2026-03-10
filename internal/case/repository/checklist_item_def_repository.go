package repository

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SCheckListItemDef struct {
	db *pgxpool.Pool
}

func NewCheckListItemDefRepository(db *pgxpool.Pool) *SCheckListItemDef {
	return &SCheckListItemDef{db: db}
}

func (r *SCheckListItemDef) Create(ctx context.Context, req dto.ChecklistItemDefCreateRequest) (*dto.ChecklistItemDef, error) {
	query := `
		INSERT INTO sicou.checklist_item_def (
			code,
			name,
			description,
			requires_document_kind_id,
			required_for_close,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5, COALESCE($6, true))
		RETURNING
			id,
			code,
			name,
			description,
			requires_document_kind_id,
			required_for_close,
			is_active,
			created_at,
			updated_at
	`

	var item dto.ChecklistItemDef

	err := r.db.QueryRow(ctx, query,
		req.Code,
		req.Name,
		req.Description,
		req.RequiresDocumentKindID,
		req.RequiredForClose,
		req.IsActive,
	).Scan(
		&item.ID,
		&item.Code,
		&item.Name,
		&item.Description,
		&item.RequiresDocumentKindID,
		&item.RequiredForClose,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &item, nil
}
func (r *SCheckListItemDef) Update(ctx context.Context, id uuid.UUID, req dto.ChecklistItemDefUpdateRequest) (*dto.ChecklistItemDef, error) {
	query := `
		UPDATE sicou.checklist_item_def
		SET
			code = COALESCE($2, code),
			name = COALESCE($3, name),
			description = COALESCE($4, description),
			requires_document_kind_id = COALESCE($5, requires_document_kind_id),
			required_for_close = COALESCE($6, required_for_close),
			is_active = COALESCE($7, is_active),
			updated_at = now()
		WHERE id = $1
		RETURNING
			id,
			code,
			name,
			description,
			requires_document_kind_id,
			required_for_close,
			is_active,
			created_at,
			updated_at
	`

	var item dto.ChecklistItemDef

	err := r.db.QueryRow(ctx, query,
		id,
		req.Code,
		req.Name,
		req.Description,
		req.RequiresDocumentKindID,
		req.RequiredForClose,
		req.IsActive,
	).Scan(
		&item.ID,
		&item.Code,
		&item.Name,
		&item.Description,
		&item.RequiresDocumentKindID,
		&item.RequiredForClose,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &item, nil
}
func (r *SCheckListItemDef) List(ctx context.Context, limit, offset int) ([]dto.ChecklistItemDef, error) {
	query := `
		SELECT
			id,
			code,
			name,
			description,
			requires_document_kind_id,
			required_for_close,
			is_active,
			created_at,
			updated_at
		FROM sicou.checklist_item_def
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.ChecklistItemDef

	for rows.Next() {
		var item dto.ChecklistItemDef
		if err := rows.Scan(
			&item.ID,
			&item.Code,
			&item.Name,
			&item.Description,
			&item.RequiresDocumentKindID,
			&item.RequiredForClose,
			&item.IsActive,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
func (r *SCheckListItemDef) GetById(ctx context.Context, id uuid.UUID) (*dto.ChecklistItemDef, error) {
	query := `
		SELECT
			id,
			code,
			name,
			description,
			requires_document_kind_id,
			required_for_close,
			is_active,
			created_at,
			updated_at
		FROM sicou.checklist_item_def
		WHERE id = $1
	`

	var item dto.ChecklistItemDef

	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.Code,
		&item.Name,
		&item.Description,
		&item.RequiresDocumentKindID,
		&item.RequiredForClose,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &item, nil
}
func (r *SCheckListItemDef) Delete(ctx context.Context, id uuid.UUID) (*dto.ChecklistItemDef, error) {
	query := `
		DELETE FROM sicou.checklist_item_def
		WHERE id = $1
		RETURNING
			id,
			code,
			name,
			description,
			requires_document_kind_id,
			required_for_close,
			is_active,
			created_at,
			updated_at
	`

	var item dto.ChecklistItemDef

	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.Code,
		&item.Name,
		&item.Description,
		&item.RequiresDocumentKindID,
		&item.RequiredForClose,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &item, nil
}
