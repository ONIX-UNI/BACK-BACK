package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *PostgresRepository) List(ctx context.Context, in ListInput) (ListResult, error) {
	return r.list(ctx, in.Page, in.Limit, nil, false)
}

func (r *PostgresRepository) ListAssigned(ctx context.Context, in ListAssignedInput) (ListResult, error) {
	return r.list(ctx, in.Page, in.Limit, in.AssignedCoordinatorID, true)
}

func (r *PostgresRepository) list(
	ctx context.Context,
	page int,
	limit int,
	assignedCoordinatorID *uuid.UUID,
	assignedOnly bool,
) (ListResult, error) {
	if r.db == nil {
		return ListResult{}, fmt.Errorf("database connection is not initialized")
	}

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return ListResult{}, err
	}
	if err := ensurePreturnoTables(ctx, tx); err != nil {
		_ = tx.Rollback(ctx)
		return ListResult{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return ListResult{}, err
	}

	conditions := make([]string, 0, 3)
	args := make([]any, 0, 4)
	argPos := 1

	conditions = append(conditions, "1=1")
	if assignedOnly {
		conditions = append(conditions, "p.assigned_coordinator_id IS NOT NULL")
	}
	if assignedCoordinatorID != nil {
		conditions = append(conditions, fmt.Sprintf("p.assigned_coordinator_id = $%d::uuid", argPos))
		args = append(args, *assignedCoordinatorID)
		argPos++
	}

	whereClause := strings.Join(conditions, " AND ")
	query := fmt.Sprintf(`
		SELECT
			p.id::text AS id,
			COALESCE(tracking_seq, 0) AS tracking_seq,
			p.status,
			p.created_at,
			p.consultation_at,
			CASE WHEN p.data_consent_accepted THEN 'Si' ELSE 'No' END AS accepts_data_processing,
			COALESCE(p.full_name_snapshot, '') AS full_name,
			COALESCE(cdt.name, cdt.code, '') AS document_type,
			COALESCE(p.document_number_snapshot, '') AS document_number,
			COALESCE(to_char(p.birth_date_snapshot, 'YYYY-MM-DD'), '') AS birth_date,
			COALESCE(p.age_years_snapshot::text, '') AS age,
			COALESCE(ccs.name, ccs.code, '') AS civil_status,
			COALESCE(cg.name, cg.code, '') AS gender,
			COALESCE(p.address_snapshot, '') AS address,
			COALESCE(cht.name, cht.code, '') AS housing_type,
			COALESCE(p.stratum_snapshot::text, '') AS stratum,
			COALESCE(p.sisben_category_snapshot, '') AS sisben_category,
			COALESCE(p.phone_mobile_snapshot, '') AS mobile_phone,
			COALESCE(p.email_snapshot::text, '') AS email,
			COALESCE(cpt.name, cpt.code, '') AS population_type,
			CASE
				WHEN p.head_of_household IS TRUE THEN 'Si'
				WHEN p.head_of_household IS FALSE THEN 'No'
				ELSE ''
			END AS head_of_household,
			COALESCE(p.occupation_snapshot, '') AS occupation,
			COALESCE(cel.name, cel.code, '') AS education_level,
			COALESCE(p.situation_story, '') AS case_description,
			CASE WHEN p.notify_by_email_consent THEN 'Si' ELSE 'No' END AS authorizes_notification
		FROM sicou.preturno p
		LEFT JOIN sicou.catalog_document_type cdt ON cdt.id = p.document_type_id_snapshot
		LEFT JOIN sicou.catalog_civil_status ccs ON ccs.id = p.civil_status_id
		LEFT JOIN sicou.catalog_gender cg ON cg.id = p.gender_id
		LEFT JOIN sicou.catalog_housing_type cht ON cht.id = p.housing_type_id
		LEFT JOIN sicou.catalog_population_type cpt ON cpt.id = p.population_type_id
		LEFT JOIN sicou.catalog_education_level cel ON cel.id = p.education_level_id
		WHERE %s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	rows, err := r.db.Query(ctx, query, append(args, limit, offset)...)
	if err != nil {
		return ListResult{}, err
	}
	defer rows.Close()

	items := make([]ListItem, 0, limit)
	for rows.Next() {
		var (
			item                   ListItem
			seq                    int64
			acceptsDataProcessing  string
			fullName               string
			documentType           string
			documentNumber         string
			birthDate              string
			age                    string
			civilStatus            string
			gender                 string
			address                string
			housingType            string
			stratum                string
			sisbenCategory         string
			mobilePhone            string
			email                  string
			populationType         string
			headOfHousehold        string
			occupation             string
			educationLevel         string
			caseDescription        string
			authorizesNotification string
		)
		if err := rows.Scan(
			&item.ID,
			&seq,
			&item.Status,
			&item.CreatedAt,
			&item.ConsultationDate,
			&acceptsDataProcessing,
			&fullName,
			&documentType,
			&documentNumber,
			&birthDate,
			&age,
			&civilStatus,
			&gender,
			&address,
			&housingType,
			&stratum,
			&sisbenCategory,
			&mobilePhone,
			&email,
			&populationType,
			&headOfHousehold,
			&occupation,
			&educationLevel,
			&caseDescription,
			&authorizesNotification,
		); err != nil {
			return ListResult{}, err
		}

		if seq > 0 {
			item.PreturnoNumber = formatPreturnoNumber(item.CreatedAt, seq)
		} else {
			item.PreturnoNumber = fmt.Sprintf("PT-%d-0000", item.CreatedAt.Year())
		}

		item.Payload = IntakePayload{
			AcceptsDataProcessing:  acceptsDataProcessing,
			ConsultationDate:       item.ConsultationDate.Format(time.RFC3339),
			FullName:               fullName,
			DocumentType:           documentType,
			DocumentNumber:         documentNumber,
			BirthDate:              birthDate,
			Age:                    age,
			MaritalStatus:          civilStatus,
			Gender:                 gender,
			Address:                address,
			HousingType:            housingType,
			SocioEconomicStratum:   stratum,
			SisbenCategory:         sisbenCategory,
			MobilePhone:            mobilePhone,
			Email:                  email,
			PopulationType:         populationType,
			HeadOfHousehold:        headOfHousehold,
			Occupation:             occupation,
			EducationLevel:         educationLevel,
			CaseDescription:        caseDescription,
			AuthorizesNotification: authorizesNotification,
			SubmittedAt:            item.CreatedAt.Format(time.RFC3339),
		}

		attachments, err := r.listAttachmentsByPreturno(ctx, item.ID)
		if err != nil {
			return ListResult{}, err
		}
		item.Attachments = attachments

		timeline, err := r.listTimelineByPreturno(ctx, item.ID)
		if err != nil {
			return ListResult{}, err
		}
		item.Timeline = timeline

		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return ListResult{}, err
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM sicou.preturno p
		WHERE %s
	`, whereClause)

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return ListResult{}, err
	}

	return ListResult{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (r *PostgresRepository) listAttachmentsByPreturno(ctx context.Context, preturnoID string) ([]ListAttachment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			d.id::text AS id,
			fo.id::text AS file_object_id,
			CASE
				WHEN dk.code = 'ID_DOC' THEN 'identityDocument'
				WHEN dk.code = 'FACTURA_SERVICIO' THEN 'utilityBill'
				ELSE lower(dk.code)
			END AS field_name,
			fo.original_name,
			COALESCE(NULLIF(BTRIM(fo.mime_type), ''), 'application/octet-stream') AS mime_type,
			COALESCE(fo.size_bytes, 0) AS size_bytes,
			'/api/v1/documents/files/' || fo.id::text || '/download' AS file_url,
			d.uploaded_at AS created_at
		FROM sicou.document d
		INNER JOIN sicou.file_object fo ON fo.id = d.file_id
		INNER JOIN sicou.document_kind dk ON dk.id = d.document_kind_id
		WHERE d.preturno_id = $1::uuid
			AND d.deleted_at IS NULL
			AND d.is_current = true
		ORDER BY d.uploaded_at ASC
	`, preturnoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]ListAttachment, 0)
	for rows.Next() {
		var item ListAttachment
		if err := rows.Scan(
			&item.ID,
			&item.FileObjectID,
			&item.FieldName,
			&item.OriginalName,
			&item.MimeType,
			&item.SizeBytes,
			&item.URL,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *PostgresRepository) listTimelineByPreturno(ctx context.Context, preturnoID string) ([]ListTimelineItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			pt.id::text AS id,
			pt.title,
			COALESCE(pt.detail, '') AS detail,
			COALESCE(NULLIF(BTRIM(au.display_name), ''), '') AS created_by_name,
			COALESCE(NULLIF(BTRIM(pt.source), ''), 'lista_interna') AS source,
			pt.created_at
		FROM sicou.preturno_timeline pt
		LEFT JOIN sicou.app_user au ON au.id = pt.created_by
		WHERE pt.preturno_id = $1::uuid
		ORDER BY pt.created_at ASC, pt.id ASC
	`, preturnoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]ListTimelineItem, 0)
	for rows.Next() {
		var item ListTimelineItem
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Detail,
			&item.User,
			&item.Source,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
