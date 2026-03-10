package repository

import (
	"context"
	"errors"
	"strings"
)

func (r *PostgresRepository) List(ctx context.Context, in ListInput) (ListResult, error) {
	if r.db == nil {
		return ListResult{}, errors.New("database connection is not initialized")
	}

	page := in.Page
	if page < 1 {
		page = 1
	}

	limit := in.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	search := strings.TrimSpace(in.Search)
	tipo := strings.ToLower(strings.TrimSpace(in.Tipo))
	statuses := sanitizeStatuses(in.Statuses)
	applyStatuses := len(statuses) > 0

	rows, err := r.db.Query(ctx, `
		SELECT
			p.id,
			FORMAT('PQRS-%s-%s', EXTRACT(YEAR FROM p.created_at)::int, LPAD(COALESCE(p.radicado, 0)::text, 4, '0')) AS radicado,
			LOWER(COALESCE(NULLIF(BTRIM((p.body::jsonb ->> 'requestType')), ''), '')) AS tipo,
			p.subject,
			COALESCE(
				NULLIF(
					BTRIM(CONCAT_WS(' ',
						p.body::jsonb ->> 'firstName',
						p.body::jsonb ->> 'middleName',
						p.body::jsonb ->> 'firstLastName',
						p.body::jsonb ->> 'secondLastName'
					)),
					''
				),
				'No informado'
			) AS ciudadano,
			COALESCE(
				NULLIF(BTRIM(p.from_email::text), ''),
				NULLIF(BTRIM((p.body::jsonb ->> 'email')), ''),
				''
			) AS correo,
			p.status,
			COALESCE(NULLIF(BTRIM(au.display_name), ''), '') AS responsable,
			p.received_at
		FROM sicou.pqrs p
		LEFT JOIN sicou.app_user au ON au.id = p.assigned_to
		WHERE
			($1 = '' OR (
				FORMAT('PQRS-%s-%s', EXTRACT(YEAR FROM p.created_at)::int, LPAD(COALESCE(p.radicado, 0)::text, 4, '0')) ILIKE '%' || $1 || '%'
				OR p.subject ILIKE '%' || $1 || '%'
				OR COALESCE(BTRIM((p.body::jsonb ->> 'requestDescription')), '') ILIKE '%' || $1 || '%'
				OR COALESCE(BTRIM((p.body::jsonb ->> 'firstName')), '') ILIKE '%' || $1 || '%'
				OR COALESCE(BTRIM((p.body::jsonb ->> 'firstLastName')), '') ILIKE '%' || $1 || '%'
				OR COALESCE(BTRIM((p.body::jsonb ->> 'email')), '') ILIKE '%' || $1 || '%'
				OR COALESCE(BTRIM(p.from_email::text), '') ILIKE '%' || $1 || '%'
			))
			AND ($2 = '' OR LOWER(COALESCE(NULLIF(BTRIM((p.body::jsonb ->> 'requestType')), ''), '')) = $2)
			AND (NOT $3 OR p.status = ANY($4::text[]))
		ORDER BY p.created_at DESC
		LIMIT $5 OFFSET $6
	`, search, tipo, applyStatuses, statuses, limit, offset)
	if err != nil {
		return ListResult{}, err
	}
	defer rows.Close()

	items := make([]ListItem, 0, limit)
	for rows.Next() {
		var item ListItem
		if err := rows.Scan(
			&item.ID,
			&item.Radicado,
			&item.Tipo,
			&item.Asunto,
			&item.Ciudadano,
			&item.Correo,
			&item.EstadoDB,
			&item.Responsable,
			&item.ReceivedAt,
		); err != nil {
			return ListResult{}, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return ListResult{}, err
	}

	var total int64
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM sicou.pqrs p
		WHERE
			($1 = '' OR (
				FORMAT('PQRS-%s-%s', EXTRACT(YEAR FROM p.created_at)::int, LPAD(COALESCE(p.radicado, 0)::text, 4, '0')) ILIKE '%' || $1 || '%'
				OR p.subject ILIKE '%' || $1 || '%'
				OR COALESCE(BTRIM((p.body::jsonb ->> 'requestDescription')), '') ILIKE '%' || $1 || '%'
				OR COALESCE(BTRIM((p.body::jsonb ->> 'firstName')), '') ILIKE '%' || $1 || '%'
				OR COALESCE(BTRIM((p.body::jsonb ->> 'firstLastName')), '') ILIKE '%' || $1 || '%'
				OR COALESCE(BTRIM((p.body::jsonb ->> 'email')), '') ILIKE '%' || $1 || '%'
				OR COALESCE(BTRIM(p.from_email::text), '') ILIKE '%' || $1 || '%'
			))
			AND ($2 = '' OR LOWER(COALESCE(NULLIF(BTRIM((p.body::jsonb ->> 'requestType')), ''), '')) = $2)
			AND (NOT $3 OR p.status = ANY($4::text[]))
	`, search, tipo, applyStatuses, statuses).Scan(&total); err != nil {
		return ListResult{}, err
	}

	return ListResult{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func sanitizeStatuses(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.ToUpper(strings.TrimSpace(value))
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}
