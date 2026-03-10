package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (r *PostgresRepository) Create(ctx context.Context, in CreateRecord) (CreateResult, error) {
	if r.db == nil {
		return CreateResult{}, errors.New("database connection is not initialized")
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return CreateResult{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := ensurePQRSAttachmentTable(ctx, tx); err != nil {
		return CreateResult{}, err
	}

	result := CreateResult{}
	var radicadoSerial int64
	err = tx.QueryRow(ctx, `
		INSERT INTO sicou.pqrs (
			id,
			from_email,
			subject,
			body,
			received_at,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, radicado, status, created_at
	`,
		in.ID,
		in.FromEmail,
		in.Subject,
		in.Body,
		in.ReceivedAt,
		pqrsStatusDB,
	).Scan(&result.ID, &radicadoSerial, &result.EstadoDB, &result.CreatedAt)
	if err != nil {
		return CreateResult{}, err
	}
	result.Radicado = formatRadicado(result.CreatedAt, radicadoSerial)

	for _, attachment := range in.Attachments {
		fileObjectID := strings.TrimSpace(attachment.FileObjectID)
		if fileObjectID == "" {
			return CreateResult{}, errors.New("attachment file object id is required")
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO sicou.pqrs_attachment (
				pqrs_id,
				file_object_id,
				file_url
			)
			VALUES ($1, $2, $3)
		`,
			in.ID,
			fileObjectID,
			attachment.URL,
		)
		if err != nil {
			return CreateResult{}, err
		}
	}

	internalRecipients := sanitizeRecipients(in.Email.InternalRecipients)
	if len(internalRecipients) > 0 {
		_, err = tx.Exec(ctx, `
			INSERT INTO sicou.email_outbox (
				to_emails,
				cc_emails,
				subject,
				body,
				status,
				related_pqrs_id
			)
			VALUES ($1, $2, $3, $4, 'PENDIENTE', $5)
		`,
			internalRecipients,
			[]string{},
			fmt.Sprintf("Nueva PQRS recibida - Radicado %s", result.Radicado),
			buildInternalEmailBody(result.Radicado, in.Subject, in.Email.CitizenName, in.FromEmail, in.Body),
			in.ID,
		)
		if err != nil {
			return CreateResult{}, err
		}
	}

	if in.Email.NotifyCitizen && strings.TrimSpace(in.Email.CitizenEmail) != "" {
		_, err = tx.Exec(ctx, `
			INSERT INTO sicou.email_outbox (
				to_emails,
				cc_emails,
				subject,
				body,
				status,
				related_pqrs_id
			)
			VALUES ($1, $2, $3, $4, 'PENDIENTE', $5)
		`,
			[]string{strings.TrimSpace(in.Email.CitizenEmail)},
			[]string{},
			fmt.Sprintf("Acuse de recibido PQRS - Radicado %s", result.Radicado),
			buildCitizenAckBody(result.Radicado, in.Email.CitizenName),
			in.ID,
		)
		if err != nil {
			return CreateResult{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return CreateResult{}, err
	}

	return result, nil
}

func ensurePQRSAttachmentTable(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS sicou.pqrs_attachment (
			id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			pqrs_id        uuid NOT NULL REFERENCES sicou.pqrs(id) ON DELETE CASCADE,
			file_object_id uuid NOT NULL REFERENCES sicou.file_object(id) ON DELETE RESTRICT,
			file_url       text NOT NULL,
			created_at     timestamptz NOT NULL DEFAULT now()
		)
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_pqrs_attachment_pqrs
		ON sicou.pqrs_attachment(pqrs_id)
	`)
	return err
}

func sanitizeRecipients(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}
