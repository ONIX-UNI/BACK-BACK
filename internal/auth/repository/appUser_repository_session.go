package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (r *AppUserInstance) UpdateLastAccess(ctx context.Context, id uuid.UUID, accessedAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE sicou.app_user
		SET last_access_at = $2
		WHERE id = $1
	`, id, accessedAt.UTC())
	return err
}

func (r *AppUserInstance) IsTokenRevoked(ctx context.Context, tokenSignature string) (bool, error) {
	var revoked bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM sicou.auth_revoked_token
			WHERE token_signature = $1
			  AND expires_at > now()
		)
	`, tokenSignature).Scan(&revoked)
	if err != nil {
		return false, err
	}

	return revoked, nil
}

func (r *AppUserInstance) RevokeToken(ctx context.Context, tokenSignature string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sicou.auth_revoked_token (
			token_signature,
			expires_at
		)
		VALUES ($1, $2)
		ON CONFLICT (token_signature) DO UPDATE
		SET expires_at = EXCLUDED.expires_at,
			revoked_at = now()
	`, tokenSignature, expiresAt)
	return err
}
