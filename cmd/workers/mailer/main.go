package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/DuvanRozoParra/sicou/pkg/config"
	"github.com/DuvanRozoParra/sicou/pkg/database"
	"github.com/DuvanRozoParra/sicou/pkg/mail"
	"github.com/DuvanRozoParra/sicou/pkg/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	mailer, err := mail.NewSMTPClient(mail.SMTPConfig{
		Host:               cfg.SMTP.Host,
		Port:               cfg.SMTP.Port,
		Username:           cfg.SMTP.Username,
		Password:           cfg.SMTP.Password,
		AccessToken:        cfg.SMTP.AccessToken,
		AuthMethod:         mail.AuthMethod(cfg.SMTP.AuthMethod),
		FromName:           cfg.SMTP.FromName,
		FromAddress:        cfg.SMTP.FromAddress,
		UseStartTLS:        cfg.SMTP.UseStartTLS,
		UseImplicitTLS:     cfg.SMTP.UseImplicitTLS,
		RequireTLS:         cfg.SMTP.RequireTLS,
		InsecureSkipVerify: cfg.SMTP.InsecureSkipVerify,
		Timeout:            cfg.SMTP.Timeout,
	})
	if err != nil {
		log.Fatalf("invalid smtp configuration: %v", err)
	}

	log.Printf("mailer worker ready: host=%s port=%d auth=%s tls(starttls=%t implicit=%t)",
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.AuthMethod,
		cfg.SMTP.UseStartTLS,
		cfg.SMTP.UseImplicitTLS,
	)

	if testTo := strings.TrimSpace(os.Getenv("MAILER_TEST_TO")); testTo != "" {
		recipients := splitEmails(testTo)
		if len(recipients) == 0 {
			log.Fatal("MAILER_TEST_TO is set but no valid emails were found")
		}

		subject := strings.TrimSpace(os.Getenv("MAILER_TEST_SUBJECT"))
		if subject == "" {
			subject = "SICOU SMTP test"
		}

		textBody := strings.TrimSpace(os.Getenv("MAILER_TEST_TEXT"))
		if textBody == "" {
			textBody = "Correo de prueba SMTP desde SICOU."
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := mailer.Send(ctx, mail.Message{
			To:       recipients,
			Subject:  subject,
			TextBody: textBody,
		}); err != nil {
			log.Fatalf("smtp smoke test failed: %v", err)
		}

		log.Printf("smtp smoke test sent successfully to: %s", strings.Join(recipients, ", "))
		return
	}

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()

	db, err := database.NewPostgresPool(dbCtx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	minioClient := storage.NewMinioClient()
	pqrsBucket := pqrsBucketName()

	pollInterval := time.Duration(envInt("MAILER_POLL_INTERVAL_SECONDS", 5)) * time.Second
	sendTimeout := time.Duration(envInt("MAILER_SEND_TIMEOUT_SECONDS", 30)) * time.Second
	retryDelay := time.Duration(envInt("MAILER_RETRY_DELAY_SECONDS", 60)) * time.Second
	batchSize := envInt("MAILER_BATCH_SIZE", 20)
	maxAttempts := envInt("MAILER_MAX_ATTEMPTS", 5)

	log.Printf(
		"mailer outbox loop started: poll=%s batch=%d timeout=%s retry_delay=%s max_attempts=%d",
		pollInterval,
		batchSize,
		sendTimeout,
		retryDelay,
		maxAttempts,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	go func() {
		<-quit
		cancel()
	}()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		processed, err := processBatch(
			ctx,
			db,
			mailer,
			minioClient,
			pqrsBucket,
			batchSize,
			sendTimeout,
			retryDelay,
			maxAttempts,
		)
		if err != nil {
			log.Printf("mailer batch failed: %v", err)
		}

		if processed > 0 {
			continue
		}

		select {
		case <-ctx.Done():
			log.Println("mailer worker stopped")
			return
		case <-ticker.C:
		}
	}
}

type outboxEmail struct {
	ID            string
	ToEmails      []string
	CCEmails      []string
	Subject       string
	Body          string
	RelatedPQRSID *string
	Attempts      int
}

type pqrsAttachment struct {
	OriginalName string
	MimeType     string
	StorageKey   string
}

func processBatch(
	ctx context.Context,
	db *pgxpool.Pool,
	sender mail.Sender,
	minioClient *storage.Client,
	documentsBucket string,
	batchSize int,
	sendTimeout time.Duration,
	retryDelay time.Duration,
	maxAttempts int,
) (int, error) {
	emails, err := claimPendingEmails(ctx, db, batchSize)
	if err != nil {
		return 0, err
	}
	if len(emails) == 0 {
		return 0, nil
	}

	for _, email := range emails {
		sendCtx, cancel := context.WithTimeout(ctx, sendTimeout)
		message := mail.Message{
			To:      email.ToEmails,
			Cc:      email.CCEmails,
			Subject: email.Subject,
		}
		if looksLikeHTMLBody(email.Body) {
			message.HTMLBody = email.Body
			message.TextBody = htmlFallbackText(email.Body)
		} else {
			message.TextBody = email.Body
		}

		if email.RelatedPQRSID != nil && shouldAttachPQRSFiles(email.Subject) {
			pqrsID := strings.TrimSpace(*email.RelatedPQRSID)
			if pqrsID != "" {
				attachments, attachErr := loadPQRSAttachments(sendCtx, db, minioClient, documentsBucket, pqrsID)
				if attachErr != nil {
					cancel()
					if updateErr := markEmailFailed(ctx, db, email.ID, email.Attempts, maxAttempts, retryDelay, attachErr); updateErr != nil {
						log.Printf("mailer failed to persist attachment error for email %s: attach_err=%v persist_err=%v", email.ID, attachErr, updateErr)
						continue
					}
					log.Printf("mailer attachment load failed for email %s: %v", email.ID, attachErr)
					continue
				}
				message.Attachments = attachments
			}
		}

		err := sender.Send(sendCtx, message)
		cancel()

		if err != nil {
			if updateErr := markEmailFailed(ctx, db, email.ID, email.Attempts, maxAttempts, retryDelay, err); updateErr != nil {
				log.Printf("mailer failed to persist error for email %s: send_err=%v persist_err=%v", email.ID, err, updateErr)
				continue
			}
			log.Printf("mailer send failed for email %s: %v", email.ID, err)
			continue
		}

		if err := markEmailSent(ctx, db, email.ID); err != nil {
			log.Printf("mailer sent email %s but failed to persist status: %v", email.ID, err)
			continue
		}

		log.Printf("mailer sent email %s", email.ID)
	}

	return len(emails), nil
}

func claimPendingEmails(ctx context.Context, db *pgxpool.Pool, batchSize int) ([]outboxEmail, error) {
	rows, err := db.Query(ctx, `
		WITH claimed AS (
			SELECT id
			FROM sicou.email_outbox
			WHERE status = 'PENDIENTE'
				AND scheduled_at <= now()
			ORDER BY scheduled_at ASC, created_at ASC
			FOR UPDATE SKIP LOCKED
			LIMIT $1
		)
		UPDATE sicou.email_outbox eo
		SET
			status = 'ENVIANDO',
			attempts = eo.attempts + 1,
			last_error = NULL
		FROM claimed
		WHERE eo.id = claimed.id
		RETURNING eo.id, eo.to_emails, eo.cc_emails, eo.subject, eo.body, eo.related_pqrs_id, eo.attempts
	`, batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	emails := make([]outboxEmail, 0, batchSize)
	for rows.Next() {
		var email outboxEmail
		if err := rows.Scan(
			&email.ID,
			&email.ToEmails,
			&email.CCEmails,
			&email.Subject,
			&email.Body,
			&email.RelatedPQRSID,
			&email.Attempts,
		); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}

func markEmailSent(ctx context.Context, db *pgxpool.Pool, id string) error {
	_, err := db.Exec(ctx, `
		UPDATE sicou.email_outbox
		SET
			status = 'ENVIADO',
			sent_at = now(),
			last_error = NULL
		WHERE id = $1
	`, id)
	return err
}

func markEmailFailed(
	ctx context.Context,
	db *pgxpool.Pool,
	id string,
	attempts int,
	maxAttempts int,
	retryDelay time.Duration,
	sendErr error,
) error {
	errText := truncate(strings.TrimSpace(sendErr.Error()), 4000)
	if attempts >= maxAttempts {
		_, err := db.Exec(ctx, `
			UPDATE sicou.email_outbox
			SET
				status = 'FALLIDO',
				last_error = $2
			WHERE id = $1
		`, id, errText)
		return err
	}

	retryDelaySeconds := int(retryDelay.Seconds())
	if retryDelaySeconds < 1 {
		retryDelaySeconds = 1
	}

	_, err := db.Exec(ctx, `
		UPDATE sicou.email_outbox
		SET
			status = 'PENDIENTE',
			last_error = $2,
			scheduled_at = now() + ($3 * INTERVAL '1 second')
		WHERE id = $1
	`, id, errText, retryDelaySeconds)
	return err
}

func splitEmails(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		email := strings.TrimSpace(part)
		if email == "" {
			continue
		}
		out = append(out, email)
	}
	return out
}

func looksLikeHTMLBody(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return false
	}

	lower := strings.ToLower(trimmed)
	return strings.Contains(lower, "<html") ||
		strings.Contains(lower, "<body") ||
		strings.Contains(lower, "<a ") ||
		strings.Contains(lower, "<p") ||
		strings.Contains(lower, "<div")
}

func htmlFallbackText(htmlBody string) string {
	re := regexp.MustCompile(`(?i)href=["']([^"']+)["']`)
	match := re.FindStringSubmatch(htmlBody)
	if len(match) == 2 && strings.TrimSpace(match[1]) != "" {
		return "Recibiste un correo en formato HTML.\nSi no ves el boton, abre este enlace:\n" + strings.TrimSpace(match[1])
	}

	return "Recibiste un correo en formato HTML."
}

func envInt(key string, defaultValue int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		log.Printf("invalid %s=%q, using default %d", key, raw, defaultValue)
		return defaultValue
	}
	return parsed
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max]
}

func loadPQRSAttachments(
	ctx context.Context,
	db *pgxpool.Pool,
	minioClient *storage.Client,
	bucket string,
	pqrsID string,
) ([]mail.Attachment, error) {
	if minioClient == nil {
		return nil, errors.New("minio client is not initialized")
	}
	if strings.TrimSpace(bucket) == "" {
		return nil, errors.New("documents bucket is not configured")
	}

	records, err := listPQRSAttachments(ctx, db, pqrsID)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}

	attachments := make([]mail.Attachment, 0, len(records))
	for _, record := range records {
		object, err := minioClient.GetObject(
			ctx,
			bucket,
			record.StorageKey,
			storage.GetObjectOptions{},
		)
		if err != nil {
			return nil, err
		}

		data, readErr := io.ReadAll(object)
		closeErr := object.Close()
		if readErr != nil {
			return nil, readErr
		}
		if closeErr != nil {
			return nil, closeErr
		}

		attachments = append(attachments, mail.Attachment{
			Filename:    record.OriginalName,
			ContentType: record.MimeType,
			Data:        data,
		})
	}

	return attachments, nil
}

func listPQRSAttachments(
	ctx context.Context,
	db *pgxpool.Pool,
	pqrsID string,
) ([]pqrsAttachment, error) {
	rows, err := db.Query(ctx, `
		SELECT
			fo.original_name,
			COALESCE(NULLIF(BTRIM(fo.mime_type), ''), 'application/octet-stream') AS mime_type,
			fo.storage_key
		FROM sicou.pqrs_attachment pa
		INNER JOIN sicou.file_object fo ON fo.id = pa.file_object_id
		WHERE pa.pqrs_id = $1
		ORDER BY pa.created_at ASC
	`, pqrsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]pqrsAttachment, 0)
	for rows.Next() {
		var record pqrsAttachment
		if err := rows.Scan(
			&record.OriginalName,
			&record.MimeType,
			&record.StorageKey,
		); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func pqrsBucketName() string {
	bucket := strings.TrimSpace(os.Getenv("PQRS_BUCKET"))
	if bucket == "" {
		bucket = "pqrsd"
	}
	return bucket
}

func shouldAttachPQRSFiles(subject string) bool {
	normalized := strings.ToLower(strings.TrimSpace(subject))
	return strings.Contains(normalized, "nueva pqrs recibida")
}
