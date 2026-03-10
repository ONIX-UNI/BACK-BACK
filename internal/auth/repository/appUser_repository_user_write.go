package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/DuvanRozoParra/sicou/internal/auth/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *AppUserInstance) Create(ctx context.Context, req dto.CreateAppUserRequest) (*dto.AppUser, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	query := `
		INSERT INTO sicou.app_user (
			email,
			display_name,
			password_hash,
			is_active
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, display_name, password_hash, is_active, last_access_at, created_at, updated_at
	`

	var user dto.AppUser
	err = scanAppUser(tx.QueryRow(ctx, query,
		req.Email,
		req.DisplayName,
		req.PasswordHash,
		isActive,
	), &user)
	if err != nil {
		return nil, err
	}

	var roleID int16
	err = tx.QueryRow(ctx, `
		SELECT id
		FROM sicou.role
		WHERE upper(code) = upper($1)
	`, req.Role).Scan(&roleID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrRoleNotFound
		}
		return nil, err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO sicou.user_role (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`, user.ID, roleID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AppUserInstance) Update(ctx context.Context, id uuid.UUID, req dto.UpdateAppUserRequest) (*dto.AppUser, error) {
	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if req.DisplayName != nil {
		setClauses = append(setClauses, fmt.Sprintf("display_name = $%d", argPos))
		args = append(args, *req.DisplayName)
		argPos++
	}

	if req.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *req.IsActive)
		argPos++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	setClauses = append(setClauses, "updated_at = now()")

	query := fmt.Sprintf(`
		UPDATE sicou.app_user
		SET %s
		WHERE id = $%d
		RETURNING id, email, display_name, password_hash, is_active, last_access_at, created_at, updated_at
	`, strings.Join(setClauses, ", "), argPos)

	args = append(args, id)

	var user dto.AppUser
	err := scanAppUser(r.db.QueryRow(ctx, query, args...), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AppUserInstance) ReplacePrimaryRole(ctx context.Context, userID uuid.UUID, roleCode string) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	normalizedRole := dto.NormalizeRoleCode(roleCode)
	if normalizedRole == "" {
		return fmt.Errorf("role code is required")
	}

	var roleID int16
	if err := tx.QueryRow(ctx, `
		SELECT id
		FROM sicou.role
		WHERE upper(code) = upper($1)
	`, normalizedRole).Scan(&roleID); err != nil {
		if err == pgx.ErrNoRows {
			return models.ErrRoleNotFound
		}
		return err
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM sicou.user_role
		WHERE user_id = $1
	`, userID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO sicou.user_role (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`, userID, roleID); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *AppUserInstance) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM sicou.app_user
		WHERE id = $1
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
