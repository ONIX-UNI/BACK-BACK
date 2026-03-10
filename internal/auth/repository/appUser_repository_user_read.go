package repository

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *AppUserInstance) GetByID(ctx context.Context, id uuid.UUID) (*dto.AppUser, error) {
	query := `
		SELECT id, email, display_name, password_hash, is_active, last_access_at, created_at, updated_at
		FROM sicou.app_user
		WHERE id = $1
	`

	var user dto.AppUser
	err := scanAppUser(r.db.QueryRow(ctx, query, id), &user)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *AppUserInstance) GetByEmail(ctx context.Context, email string) (*dto.AppUser, error) {
	query := `
		SELECT id, email, display_name, password_hash, is_active, last_access_at, created_at, updated_at
		FROM sicou.app_user
		WHERE email = $1
	`

	var user dto.AppUser
	err := scanAppUser(r.db.QueryRow(ctx, query, email), &user)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *AppUserInstance) GetRolesByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	query := `
		SELECT ro.code
		FROM sicou.user_role ur
		INNER JOIN sicou.role ro ON ro.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY ro.code ASC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]string, 0)
	for rows.Next() {
		var roleCode string
		if err := rows.Scan(&roleCode); err != nil {
			return nil, err
		}
		roles = append(roles, roleCode)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return roles, nil
}

func (r *AppUserInstance) RoleExists(ctx context.Context, roleCode string) (bool, error) {
	normalizedRole := dto.NormalizeRoleCode(roleCode)
	if normalizedRole == "" {
		return false, nil
	}

	var exists bool
	if err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM sicou.role
			WHERE upper(code) = $1
		)
	`, normalizedRole).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
