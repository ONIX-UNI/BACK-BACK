package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *AppUserInstance) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	cmdTag, err := r.db.Exec(ctx, `
		UPDATE sicou.app_user
		SET password_hash = $2,
			updated_at = now()
		WHERE id = $1
	`, id, passwordHash)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
