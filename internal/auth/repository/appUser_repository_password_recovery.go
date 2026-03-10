package repository

import (
	"context"
	"strings"
)

func (r *AppUserInstance) EnqueuePasswordResetEmail(ctx context.Context, toEmail string, subject string, body string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sicou.email_outbox (
			to_emails,
			cc_emails,
			subject,
			body
		)
		VALUES ($1, $2, $3, $4)
	`,
		[]string{strings.TrimSpace(toEmail)},
		[]string{},
		strings.TrimSpace(subject),
		strings.TrimSpace(body),
	)

	return err
}
