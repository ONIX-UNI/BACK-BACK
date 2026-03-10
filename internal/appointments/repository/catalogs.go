package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type upsertCitizenInput struct {
	DocumentTypeID int32
	DocumentNumber string
	FullName       string
	BirthDate      *time.Time
	PhoneMobile    string
	Email          string
	Address        string
}

func upsertCitizen(ctx context.Context, tx pgx.Tx, in upsertCitizenInput) (uuid.UUID, error) {
	var citizenID uuid.UUID
	err := tx.QueryRow(ctx, `
		INSERT INTO sicou.citizen (
			document_type_id,
			document_number,
			full_name,
			birth_date,
			phone_mobile,
			email,
			address
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (document_type_id, document_number)
		DO UPDATE SET
			full_name = EXCLUDED.full_name,
			birth_date = COALESCE(EXCLUDED.birth_date, sicou.citizen.birth_date),
			phone_mobile = COALESCE(NULLIF(EXCLUDED.phone_mobile, ''), sicou.citizen.phone_mobile),
			email = COALESCE(NULLIF(EXCLUDED.email, ''), sicou.citizen.email),
			address = COALESCE(NULLIF(EXCLUDED.address, ''), sicou.citizen.address),
			updated_at = now(),
			deleted_at = NULL
		RETURNING id
	`,
		in.DocumentTypeID,
		strings.TrimSpace(in.DocumentNumber),
		strings.TrimSpace(in.FullName),
		in.BirthDate,
		strings.TrimSpace(in.PhoneMobile),
		strings.TrimSpace(in.Email),
		strings.TrimSpace(in.Address),
	).Scan(&citizenID)
	return citizenID, err
}

func resolveRequiredCatalogID(ctx context.Context, tx pgx.Tx, table string, rawValue string, fallbackCode string) (int16, error) {
	if id, err := resolveOptionalCatalogID(ctx, tx, table, rawValue); err != nil {
		return 0, err
	} else if id != nil {
		return *id, nil
	}

	code := strings.ToUpper(strings.TrimSpace(fallbackCode))
	if code == "" {
		code = "OTRO"
	}
	var id int16
	query := fmt.Sprintf(`SELECT id FROM sicou.%s WHERE upper(code) = $1 LIMIT 1`, table)
	if err := tx.QueryRow(ctx, query, code).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func resolveOptionalCatalogID(ctx context.Context, tx pgx.Tx, table string, rawValue string) (*int16, error) {
	value := strings.TrimSpace(rawValue)
	if value == "" {
		return nil, nil
	}

	query := fmt.Sprintf(`
		SELECT id
		FROM sicou.%s
		WHERE is_active = true
		  AND (
			upper(code) = upper($1)
			OR unaccent(lower(name)) = unaccent(lower($1))
		  )
		ORDER BY id
		LIMIT 1
	`, table)

	var id int16
	if err := tx.QueryRow(ctx, query, value).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &id, nil
}

func resolveOptionalCatalogSelection(
	ctx context.Context,
	tx pgx.Tx,
	table string,
	selected string,
	other string,
) (*int16, error) {
	selected = strings.TrimSpace(selected)
	other = strings.TrimSpace(other)
	if selected == "" && other == "" {
		return nil, nil
	}

	if isOtherSelection(selected) {
		return resolveOptionalCatalogID(ctx, tx, table, "OTRO")
	}

	if selected != "" {
		id, err := resolveOptionalCatalogID(ctx, tx, table, selected)
		if err != nil {
			return nil, err
		}
		if id != nil {
			return id, nil
		}
	}

	if other != "" {
		id, err := resolveOptionalCatalogID(ctx, tx, table, other)
		if err != nil {
			return nil, err
		}
		if id != nil {
			return id, nil
		}
	}

	return resolveOptionalCatalogID(ctx, tx, table, "OTRO")
}

func isOtherSelection(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "otro", "otra", "other":
		return true
	default:
		return false
	}
}

func resolveDocumentKindID(ctx context.Context, tx pgx.Tx, code string) (int16, error) {
	normalizedCode := strings.ToUpper(strings.TrimSpace(code))
	if normalizedCode == "" {
		normalizedCode = "OTRO"
	}

	var id int16
	if err := tx.QueryRow(ctx, `
		SELECT id
		FROM sicou.document_kind
		WHERE upper(code) = $1
		LIMIT 1
	`, normalizedCode).Scan(&id); err == nil {
		return id, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}

	if err := tx.QueryRow(ctx, `
		SELECT id
		FROM sicou.document_kind
		WHERE upper(code) = 'OTRO'
		LIMIT 1
	`).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func kindCodeFromAttachmentField(fieldName string) string {
	switch strings.ToLower(strings.TrimSpace(fieldName)) {
	case "identitydocument":
		return "ID_DOC"
	case "utilitybill":
		return "FACTURA_SERVICIO"
	default:
		return "OTRO"
	}
}
