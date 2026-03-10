package repository

import (
	"context"
	"errors"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SDesktop struct {
	db *pgxpool.Pool
}

func NewDesktopRepository(db *pgxpool.Pool) *SDesktop {
	return &SDesktop{db: db}
}

func (r *SDesktop) Create(ctx context.Context, req dto.CreateDesktopRequest) (*dto.CreateDesktopResponse, error) {
	var id string

	query := `SELECT sicou.fn_escritorio_create_empty($1, $2, $3)`
	err := r.db.QueryRow(ctx, query, req.Code, req.Name, req.Location).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &dto.CreateDesktopResponse{ID: id}, nil
}

func (r *SDesktop) List(ctx context.Context) ([]dto.Desktop, error) {
	query := `
		SELECT 
			id, code, name, location, is_active,
			current_user_id, assigned_at, assigned_by,
			created_at, updated_at
		FROM sicou.escritorio
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var desktops []dto.Desktop
	for rows.Next() {
		var d dto.Desktop
		if err := rows.Scan(
			&d.ID, &d.Code, &d.Name, &d.Location, &d.IsActive,
			&d.CurrentUserID, &d.AssignedAt, &d.AssignedBy,
			&d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		desktops = append(desktops, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return desktops, nil
}

// ✅ Update flexible: name / is_active / assign user
func (r *SDesktop) UpdateByCode(ctx context.Context, code string, req dto.UpdateDesktopRequest) (*dto.UpdateDesktopResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	updated := 0

	if req.Name != nil {
		tag, err := tx.Exec(ctx, `UPDATE sicou.escritorio SET name=$2, updated_at=now() WHERE code=$1`, code, *req.Name)
		if err != nil {
			return nil, err
		}
		if tag.RowsAffected() > 0 {
			updated++
		}
	}

	if req.IsActive != nil {
		tag, err := tx.Exec(ctx, `UPDATE sicou.escritorio SET is_active=$2, updated_at=now() WHERE code=$1`, code, *req.IsActive)
		if err != nil {
			return nil, err
		}
		if tag.RowsAffected() > 0 {
			updated++
		}
	}

	// asignación opcional (usa tu función existente)
	if req.UserID != nil {
		var escritorioID string
		if err := tx.QueryRow(ctx, `SELECT id FROM sicou.escritorio WHERE code=$1`, code).Scan(&escritorioID); err != nil {
			return nil, err
		}

		_, err := tx.Exec(ctx, `SELECT sicou.fn_escritorio_assign_user($1, $2, $3)`, escritorioID, *req.UserID, req.Reason)
		if err != nil {
			return nil, err
		}
		updated++
	}

	if updated == 0 {
		return nil, errors.New("no changes applied (desk not found or empty patch)")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &dto.UpdateDesktopResponse{Success: true, Message: "updated"}, nil
}

func (r *SDesktop) DeleteByCode(ctx context.Context, code string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM sicou.escritorio WHERE code=$1`, code)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("desk not found")
	}
	return nil
}