package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *AppUserInstance) EnsureSingleSuperAdmin(ctx context.Context, email, displayName, passwordHash string) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		INSERT INTO sicou.role(code, name)
		VALUES ('SUPER_ADMIN', 'Super Administrador (TI)')
		ON CONFLICT (code) DO NOTHING
	`); err != nil {
		return err
	}

	var existingCount int64
	if err := tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM sicou.user_role ur
		INNER JOIN sicou.role ro ON ro.id = ur.role_id
		WHERE ro.code = 'SUPER_ADMIN'
	`).Scan(&existingCount); err != nil {
		return err
	}

	if existingCount > 0 {
		return tx.Commit(ctx)
	}

	var userID uuid.UUID
	if err := tx.QueryRow(ctx, `
		INSERT INTO sicou.app_user (
			email,
			display_name,
			password_hash,
			is_active
		)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (email) DO UPDATE
		SET display_name = EXCLUDED.display_name,
			password_hash = EXCLUDED.password_hash,
			is_active = true,
			updated_at = now()
		RETURNING id
	`, email, displayName, passwordHash).Scan(&userID); err != nil {
		return err
	}

	var roleID int16
	if err := tx.QueryRow(ctx, `
		SELECT id
		FROM sicou.role
		WHERE code = 'SUPER_ADMIN'
	`).Scan(&roleID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO sicou.user_role (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`, userID, roleID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *AppUserInstance) EnsureAuthStorage(ctx context.Context) error {
	query := `
		ALTER TABLE sicou.app_user
		ADD COLUMN IF NOT EXISTS last_access_at timestamptz;

		CREATE TABLE IF NOT EXISTS sicou.auth_revoked_token (
			token_signature text PRIMARY KEY,
			expires_at timestamptz NOT NULL,
			revoked_at timestamptz NOT NULL DEFAULT now()
		);

		CREATE INDEX IF NOT EXISTS idx_auth_revoked_token_expires_at
		ON sicou.auth_revoked_token (expires_at);

		CREATE TABLE IF NOT EXISTS sicou.user_audit_log (
			id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			created_at         timestamptz NOT NULL DEFAULT now(),
			actor_user_id      uuid,
			actor_display_name text,
			actor_email        text,
			action             text NOT NULL,
			target_user_id     uuid,
			target_display_name text NOT NULL,
			target_email       text NOT NULL,
			detail             text NOT NULL,
			metadata           jsonb NOT NULL DEFAULT '{}'::jsonb
		);

		CREATE INDEX IF NOT EXISTS idx_user_audit_log_created_at
		ON sicou.user_audit_log (created_at DESC);

		CREATE INDEX IF NOT EXISTS idx_user_audit_log_action
		ON sicou.user_audit_log (action);

		CREATE INDEX IF NOT EXISTS idx_user_audit_log_actor_user
		ON sicou.user_audit_log (actor_user_id);

		CREATE INDEX IF NOT EXISTS idx_user_audit_log_target_user
		ON sicou.user_audit_log (target_user_id);

		CREATE UNIQUE INDEX IF NOT EXISTS uq_app_user_email
		ON sicou.app_user (email);
	`

	if _, err := r.db.Exec(ctx, query); err != nil {
		return err
	}

	_, _ = r.db.Exec(ctx, `
		DELETE FROM sicou.auth_revoked_token
		WHERE expires_at <= now()
	`)

	return nil
}
