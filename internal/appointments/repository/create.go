package repository

import (
	"context"
	"encoding/json"
	"errors"
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

	if err := ensurePreturnoTables(ctx, tx); err != nil {
		return CreateResult{}, err
	}

	var payload IntakePayload
	if raw := strings.TrimSpace(in.Payload); raw != "" {
		if err := json.Unmarshal([]byte(raw), &payload); err != nil {
			return CreateResult{}, err
		}
	}

	documentTypeID, err := resolveRequiredCatalogID(ctx, tx, "catalog_document_type", payload.DocumentType, "OTRO")
	if err != nil {
		return CreateResult{}, err
	}
	civilStatusID, err := resolveOptionalCatalogSelection(
		ctx,
		tx,
		"catalog_civil_status",
		payload.MaritalStatus,
		payload.OtherMaritalStatus,
	)
	if err != nil {
		return CreateResult{}, err
	}
	genderID, err := resolveOptionalCatalogID(ctx, tx, "catalog_gender", payload.Gender)
	if err != nil {
		return CreateResult{}, err
	}
	housingTypeID, err := resolveOptionalCatalogSelection(
		ctx,
		tx,
		"catalog_housing_type",
		payload.HousingType,
		payload.OtherHousingType,
	)
	if err != nil {
		return CreateResult{}, err
	}
	populationTypeID, err := resolveOptionalCatalogSelection(
		ctx,
		tx,
		"catalog_population_type",
		payload.PopulationType,
		payload.OtherPopulationType,
	)
	if err != nil {
		return CreateResult{}, err
	}
	educationLevelID, err := resolveOptionalCatalogSelection(
		ctx,
		tx,
		"catalog_education_level",
		payload.EducationLevel,
		payload.OtherEducationLevel,
	)
	if err != nil {
		return CreateResult{}, err
	}
	channelID, err := resolveRequiredCatalogID(ctx, tx, "catalog_channel", "VIRTUAL", "VIRTUAL")
	if err != nil {
		return CreateResult{}, err
	}

	normalizedFullName := strings.TrimSpace(payload.FullName)
	if normalizedFullName == "" {
		normalizedFullName = strings.TrimSpace(in.CitizenName)
	}
	if normalizedFullName == "" {
		normalizedFullName = "No informado"
	}

	normalizedDocumentNumber := strings.TrimSpace(payload.DocumentNumber)
	if normalizedDocumentNumber == "" {
		normalizedDocumentNumber = "NO-INFORMADO-" + shortID(in.ID)
	}

	birthDate := parseOptionalDate(payload.BirthDate)
	ageYears := parseOptionalSmallInt(payload.Age)
	stratum := normalizeStratum(payload.SocioEconomicStratum)

	phoneMobile := strings.TrimSpace(payload.MobilePhone)
	if phoneMobile == "" {
		phoneMobile = "NO_INFORMADO"
	}

	address := strings.TrimSpace(payload.Address)
	if address == "" {
		address = "NO_INFORMADA"
	}

	contactEmail := normalizeEmailOrPlaceholder(payload.Email, in.NotificationEmail, in.ID)
	situationStory := strings.TrimSpace(payload.CaseDescription)
	if situationStory == "" {
		situationStory = "No informado"
	}

	citizenID, err := upsertCitizen(ctx, tx, upsertCitizenInput{
		DocumentTypeID: int32(documentTypeID),
		DocumentNumber: normalizedDocumentNumber,
		FullName:       normalizedFullName,
		BirthDate:      birthDate,
		PhoneMobile:    phoneMobile,
		Email:          strings.TrimSpace(payload.Email),
		Address:        address,
	})
	if err != nil {
		return CreateResult{}, err
	}

	result := CreateResult{}
	var trackingSeq int64
	err = tx.QueryRow(ctx, `
		INSERT INTO sicou.preturno (
			id,
			citizen_id,
			contact_email,
			data_consent_text,
			data_consent_accepted,
			consultation_at,
			full_name_snapshot,
			document_type_id_snapshot,
			document_number_snapshot,
			birth_date_snapshot,
			age_years_snapshot,
			civil_status_id,
			gender_id,
			address_snapshot,
			housing_type_id,
			stratum_snapshot,
			sisben_category_snapshot,
			phone_mobile_snapshot,
			email_snapshot,
			population_type_id,
			head_of_household,
			occupation_snapshot,
			education_level_id,
			situation_story,
			notify_by_email_consent,
			channel_id,
			status
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24, $25, $26, $27
		)
		RETURNING id, status, created_at, tracking_seq
	`,
		in.ID,
		citizenID,
		contactEmail,
		"Acepta tratamiento de datos personales",
		in.AcceptsDataProcessing,
		in.ConsultationDate,
		normalizedFullName,
		documentTypeID,
		normalizedDocumentNumber,
		birthDate,
		ageYears,
		civilStatusID,
		genderID,
		address,
		housingTypeID,
		stratum,
		strings.TrimSpace(payload.SisbenCategory),
		phoneMobile,
		contactEmail,
		populationTypeID,
		in.HeadOfHousehold,
		strings.TrimSpace(payload.Occupation),
		educationLevelID,
		situationStory,
		in.AuthorizesNotification,
		channelID,
		StatusPendientePreturno,
	).Scan(&result.ID, &result.Status, &result.CreatedAt, &trackingSeq)
	if err != nil {
		return CreateResult{}, err
	}
	result.PreturnoNumber = formatPreturnoNumber(result.CreatedAt, trackingSeq)

	timelineSource := normalizeTimelineSource(in.EventSource)
	timelineDetail := buildPreturnoCreatedDetail(timelineSource)

	_, err = tx.Exec(ctx, `
		INSERT INTO sicou.preturno_timeline (
			preturno_id,
			event_type,
			title,
			detail,
			created_by,
			source
		)
		VALUES (
			$1::uuid,
			'CREATED',
			'Pre-turno registrado',
			$2,
			$3::uuid,
			$4
		)
	`, result.ID, timelineDetail, in.CreatedBy, timelineSource)
	if err != nil {
		return CreateResult{}, err
	}

	for _, attachment := range in.Attachments {
		fileObjectID := strings.TrimSpace(attachment.FileObjectID)
		if fileObjectID == "" {
			return CreateResult{}, errors.New("attachment file object id is required")
		}

		documentKindID, err := resolveDocumentKindID(ctx, tx, kindCodeFromAttachmentField(attachment.FieldName))
		if err != nil {
			return CreateResult{}, err
		}

		notes := strings.TrimSpace(attachment.FieldName)
		if notes != "" {
			notes = "Adjunto formulario legal-advise: " + notes
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO sicou.document (
				file_id,
				document_kind_id,
				preturno_id,
				notes
			)
			VALUES ($1::uuid, $2, $3::uuid, $4)
		`,
			fileObjectID,
			documentKindID,
			result.ID,
			notes,
		)
		if err != nil {
			return CreateResult{}, err
		}
	}

	if in.AuthorizesNotification && strings.TrimSpace(in.NotificationEmail) != "" {
		_, err = tx.Exec(ctx, `
			INSERT INTO sicou.email_outbox (
				to_emails,
				cc_emails,
				subject,
				body,
				status
			)
			VALUES ($1, $2, $3, $4, 'PENDIENTE')
		`,
			[]string{strings.TrimSpace(in.NotificationEmail)},
			[]string{},
			buildCitizenNotificationSubject(result.PreturnoNumber),
			buildCitizenNotificationBody(result.PreturnoNumber, in.CitizenName),
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
