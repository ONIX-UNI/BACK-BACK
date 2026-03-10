package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SLegalArea struct {
	db *pgxpool.Pool
}

func NewLegalAreaRepository(db *pgxpool.Pool) *SLegalArea {
	return &SLegalArea{db: db}
}

func (r *SLegalArea) Create(ctx context.Context, req dto.CreateLegalAreaRequest) (*dto.LegalArea, error) {

	query := `
		INSERT INTO sicou.legal_area (code, name, is_active)
		VALUES ($1, $2, COALESCE($3, true))
		RETURNING id, code, name, is_active, created_at, updated_at;
	`

	var result dto.LegalArea

	err := r.db.QueryRow(
		ctx,
		query,
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
func (r *SLegalArea) Update(ctx context.Context, id int16, req dto.UpdateLegalAreaRequest) (*dto.LegalArea, error) {

	setParts := []string{}
	args := []interface{}{}
	argPos := 1

	if req.Code != nil {
		setParts = append(setParts, fmt.Sprintf("code = $%d", argPos))
		args = append(args, *req.Code)
		argPos++
	}

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *req.Name)
		argPos++
	}

	if req.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *req.IsActive)
		argPos++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// siempre actualizar updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argPos))
	args = append(args, time.Now())
	argPos++

	query := fmt.Sprintf(`
		UPDATE sicou.legal_area
		SET %s
		WHERE id = $%d
		RETURNING id, code, name, is_active, created_at, updated_at;
	`, strings.Join(setParts, ", "), argPos)

	args = append(args, id)

	var result dto.LegalArea

	err := r.db.QueryRow(ctx, query, args...).Scan(
		&result.ID,
		&result.Code,
		&result.Name,
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("legal area not found")
		}
		return nil, err
	}

	return &result, nil
}
func (r *SLegalArea) List(ctx context.Context, limit, offset int) ([]dto.LegalArea, error) {

	query := `
		SELECT id, code, name, is_active, created_at, updated_at
		FROM sicou.legal_area
		ORDER BY id
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.LegalArea

	for rows.Next() {
		var area dto.LegalArea
		err := rows.Scan(
			&area.ID,
			&area.Code,
			&area.Name,
			&area.IsActive,
			&area.CreatedAt,
			&area.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, area)
	}

	return result, nil
}
func (r *SLegalArea) GetById(ctx context.Context, id int16) (*dto.LegalArea, error) {

	query := `
		SELECT id, code, name, is_active, created_at, updated_at
		FROM sicou.legal_area
		WHERE id = $1;
	`

	var area dto.LegalArea

	err := r.db.QueryRow(ctx, query, id).Scan(
		&area.ID,
		&area.Code,
		&area.Name,
		&area.IsActive,
		&area.CreatedAt,
		&area.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("legal area not found")
		}
		return nil, err
	}

	return &area, nil
}
func (r *SLegalArea) Delete(ctx context.Context, id int16) (*dto.LegalArea, error) {

	query := `
		UPDATE sicou.legal_area
		SET is_active = false,
		    updated_at = now()
		WHERE id = $1
		RETURNING id, code, name, is_active, created_at, updated_at;
	`

	var area dto.LegalArea

	err := r.db.QueryRow(ctx, query, id).Scan(
		&area.ID,
		&area.Code,
		&area.Name,
		&area.IsActive,
		&area.CreatedAt,
		&area.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("legal area not found")
		}
		return nil, err
	}

	return &area, nil
}
