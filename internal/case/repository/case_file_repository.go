package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SCaseFile struct {
	db *pgxpool.Pool
}

func NewSCaseFileRepository(db *pgxpool.Pool) *SCaseFile {
	return &SCaseFile{db: db}
}

func (r *SCaseFile) Create(ctx context.Context, req dto.CreateCaseFileRequest) (*dto.CaseFile, error) {
	query := `
	INSERT INTO sicou.case_file (
		citizen_id,
		preturno_id,
		service_type_id,
		legal_area_id,
		status,
		current_responsible,
		supervisor_user
	)
	VALUES ($1,$2,$3,$4,COALESCE($5,'ABIERTO'),$6,$7)
	RETURNING *
	`

	var cf dto.CaseFile

	err := r.db.QueryRow(ctx, query,
		req.CitizenID,
		req.PreturnoID,
		req.ServiceTypeID,
		req.LegalAreaID,
		req.Status,
		req.CurrentResponsible,
		req.SupervisorUser,
	).Scan(
		&cf.ID,
		&cf.CitizenID,
		&cf.PreturnoID,
		&cf.TurnID,
		&cf.ServiceTypeID,
		&cf.LegalAreaID,
		&cf.Status,
		&cf.CurrentResponsible,
		&cf.SupervisorUser,
		&cf.OpenedAt,
		&cf.ClosedAt,
		&cf.CloseNotes,
		&cf.CreatedBy,
		&cf.UpdatedBy,
		&cf.CreatedAt,
		&cf.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &cf, nil
}
func (r *SCaseFile) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCaseFileRequest) (*dto.CaseFile, error) {
	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if req.TurnID != nil {
		setClauses = append(setClauses, fmt.Sprintf("turn_id = $%d", argPos))
		args = append(args, *req.TurnID)
		argPos++
	}

	if req.ServiceTypeID != nil {
		setClauses = append(setClauses, fmt.Sprintf("service_type_id = $%d", argPos))
		args = append(args, *req.ServiceTypeID)
		argPos++
	}

	if req.LegalAreaID != nil {
		setClauses = append(setClauses, fmt.Sprintf("legal_area_id = $%d", argPos))
		args = append(args, *req.LegalAreaID)
		argPos++
	}

	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *req.Status)
		argPos++
	}

	if req.CurrentResponsible != nil {
		setClauses = append(setClauses, fmt.Sprintf("current_responsible = $%d", argPos))
		args = append(args, *req.CurrentResponsible)
		argPos++
	}

	if req.SupervisorUser != nil {
		setClauses = append(setClauses, fmt.Sprintf("supervisor_user = $%d", argPos))
		args = append(args, *req.SupervisorUser)
		argPos++
	}

	if req.ClosedAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("closed_at = $%d", argPos))
		args = append(args, *req.ClosedAt)
		argPos++
	}

	if req.CloseNotes != nil {
		setClauses = append(setClauses, fmt.Sprintf("close_notes = $%d", argPos))
		args = append(args, *req.CloseNotes)
		argPos++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	setClauses = append(setClauses, "updated_at = now()")

	query := fmt.Sprintf(`
	UPDATE sicou.case_file
	SET %s
	WHERE id = $%d
	RETURNING *
	`, strings.Join(setClauses, ","), argPos)

	args = append(args, id)

	var cf dto.CaseFile

	err := r.db.QueryRow(ctx, query, args...).Scan(
		&cf.ID,
		&cf.CitizenID,
		&cf.PreturnoID,
		&cf.TurnID,
		&cf.ServiceTypeID,
		&cf.LegalAreaID,
		&cf.Status,
		&cf.CurrentResponsible,
		&cf.SupervisorUser,
		&cf.OpenedAt,
		&cf.ClosedAt,
		&cf.CloseNotes,
		&cf.CreatedBy,
		&cf.UpdatedBy,
		&cf.CreatedAt,
		&cf.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &cf, nil
}
func (r *SCaseFile) List(ctx context.Context, limit, offset int) ([]dto.CaseFile, error) {

	query := `
	SELECT *
	FROM sicou.case_file
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []dto.CaseFile

	for rows.Next() {
		var cf dto.CaseFile

		err := rows.Scan(
			&cf.ID,
			&cf.CitizenID,
			&cf.PreturnoID,
			&cf.TurnID,
			&cf.ServiceTypeID,
			&cf.LegalAreaID,
			&cf.Status,
			&cf.CurrentResponsible,
			&cf.SupervisorUser,
			&cf.OpenedAt,
			&cf.ClosedAt,
			&cf.CloseNotes,
			&cf.CreatedBy,
			&cf.UpdatedBy,
			&cf.CreatedAt,
			&cf.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, cf)
	}

	return list, nil
}
func (r *SCaseFile) GetById(ctx context.Context, id uuid.UUID) (*dto.CaseFile, error) {
	query := `
	SELECT *
	FROM sicou.case_file
	WHERE id = $1
	`

	var cf dto.CaseFile

	err := r.db.QueryRow(ctx, query, id).Scan(
		&cf.ID,
		&cf.CitizenID,
		&cf.PreturnoID,
		&cf.TurnID,
		&cf.ServiceTypeID,
		&cf.LegalAreaID,
		&cf.Status,
		&cf.CurrentResponsible,
		&cf.SupervisorUser,
		&cf.OpenedAt,
		&cf.ClosedAt,
		&cf.CloseNotes,
		&cf.CreatedBy,
		&cf.UpdatedBy,
		&cf.CreatedAt,
		&cf.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &cf, nil
}
func (r *SCaseFile) Delete(ctx context.Context, id uuid.UUID) (*dto.CaseFile, error) {
	query := `
	DELETE FROM sicou.case_file
	WHERE id = $1
	RETURNING *
	`

	var cf dto.CaseFile

	err := r.db.QueryRow(ctx, query, id).Scan(
		&cf.ID,
		&cf.CitizenID,
		&cf.PreturnoID,
		&cf.TurnID,
		&cf.ServiceTypeID,
		&cf.LegalAreaID,
		&cf.Status,
		&cf.CurrentResponsible,
		&cf.SupervisorUser,
		&cf.OpenedAt,
		&cf.ClosedAt,
		&cf.CloseNotes,
		&cf.CreatedBy,
		&cf.UpdatedBy,
		&cf.CreatedAt,
		&cf.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &cf, nil
}
func (r *SCaseFile) ListCases(ctx context.Context) ([]dto.CaseListItem, error) {

	query := `
SELECT
    'EXP-' 
    || EXTRACT(YEAR FROM cf.opened_at) 
    || '-' 
    || LPAD(COALESCE(pt.tracking_seq, 0)::text, 4, '0') 
    AS expediente_code,
    
    cf.id,
    cf.status,
    cf.opened_at,
    cf.closed_at,

    cf.supervisor_user,
    su.display_name AS supervisor_name,

    c.full_name,
    c.document_number,

    la.name  AS legal_area,
    st.name  AS service_type,

    u.display_name AS responsible_name,

    t.consecutive,
    t.turn_date,

    folder_info.folders AS folders,

    /* Adjuntos */
    COALESCE(att.attachments_count, 0) AS attachments_count,

    /* Comentarios */
    COALESCE(att.comments_count, 0) AS comments_count,

    /* Progreso (%) */
    COALESCE(
        ROUND(
            (
                COALESCE(prog.completed_kinds, 0)::decimal
                / NULLIF(req.total_required, 0)
            ) * 100
        )::int,
    0) AS progress

FROM sicou.case_file cf

JOIN sicou.citizen c 
    ON c.id = cf.citizen_id

LEFT JOIN sicou.legal_area la 
    ON la.id = cf.legal_area_id

LEFT JOIN sicou.service_type st 
    ON st.id = cf.service_type_id

LEFT JOIN sicou.app_user u 
    ON u.id = cf.current_responsible

LEFT JOIN sicou.turn t 
    ON t.id = cf.turn_id

LEFT JOIN sicou.preturno pt 
    ON pt.id = cf.preturno_id

LEFT JOIN sicou.app_user su 
    ON su.id = cf.supervisor_user

/* Conteo adjuntos + notas */
LEFT JOIN (
    SELECT 
        case_id,
        COUNT(*) AS attachments_count,
        COUNT(*) FILTER (
            WHERE notes IS NOT NULL 
            AND TRIM(notes) <> ''
        ) AS comments_count
    FROM sicou.document
    WHERE deleted_at IS NULL
      AND is_current = true
      AND case_id IS NOT NULL
    GROUP BY case_id
) att 
    ON att.case_id = cf.id

/* Conteo tipos obligatorios completados */
LEFT JOIN (
    SELECT
        d.case_id,
        COUNT(DISTINCT d.document_kind_id) FILTER (
            WHERE dk.code <> 'OTRO'
              AND dk.is_active = true
        ) AS completed_kinds
    FROM sicou.document d
    JOIN sicou.document_kind dk 
        ON dk.id = d.document_kind_id
    WHERE d.deleted_at IS NULL
      AND d.is_current = true
      AND d.case_id IS NOT NULL
    GROUP BY d.case_id
) prog 
    ON prog.case_id = cf.id

/* Carpetas asociadas al expediente */
LEFT JOIN (
    SELECT
        d.case_id,
        STRING_AGG(DISTINCT f.name, ', ') AS folders
    FROM sicou.document d
    JOIN sicou.folder_document fd
        ON fd.document_id = d.id
    JOIN sicou.folder f
        ON f.id = fd.folder_id
       AND f.deleted_at IS NULL
    WHERE d.deleted_at IS NULL
      AND d.is_current = true
      AND d.case_id IS NOT NULL
    GROUP BY d.case_id
) folder_info
    ON folder_info.case_id = cf.id

/* Total tipos obligatorios activos */
CROSS JOIN (
    SELECT COUNT(*) AS total_required
    FROM sicou.document_kind
    WHERE code <> 'OTRO'
      AND is_active = true
) req

ORDER BY cf.created_at DESC;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []dto.CaseListItem

	for rows.Next() {
		var item dto.CaseListItem

		err := rows.Scan(
			&item.ExpedienteCode,

			&item.ID,

			&item.Status,
			&item.OpenedAt,
			&item.ClosedAt,

			&item.SupervisorID,
			&item.SupervisorName,

			&item.FullName,
			&item.DocumentNumber,

			&item.LegalArea,
			&item.ServiceType,

			&item.ResponsibleName,

			&item.Consecutive,
			&item.TurnDate,

			&item.Folders,

			&item.AttachmentsCount,
			&item.CommentsCount,
			&item.Progress,
		)

		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return list, nil
}

func (r *SCaseFile) ListCasesByEmail(
	ctx context.Context,
	email string,
) ([]dto.CaseListItem, error) {

	// 1️⃣ Resolver citizen_id
	citizenID, err := r.GetCitizenIDByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	query := `
	SELECT
		'EXP-' 
		|| EXTRACT(YEAR FROM cf.opened_at) 
		|| '-' 
		|| LPAD(COALESCE(pt.tracking_seq, 0)::text, 4, '0') 
		AS expediente_code,
		
		cf.id,
		cf.status,
		cf.opened_at,
		cf.closed_at,

		cf.supervisor_user,
		su.display_name AS supervisor_name,

		c.full_name,
		c.document_number,

		la.name  AS legal_area,
		st.name  AS service_type,

		u.display_name AS responsible_name,

		t.consecutive,
		t.turn_date,

		COALESCE(att.attachments_count, 0) AS attachments_count,
		COALESCE(att.comments_count, 0) AS comments_count,

		COALESCE(
			ROUND(
				(
					COALESCE(prog.completed_kinds, 0)::decimal
					/ NULLIF(req.total_required, 0)
				) * 100
			),
		0) AS progress

	FROM sicou.case_file cf

	JOIN sicou.citizen c 
		ON c.id = cf.citizen_id

	LEFT JOIN sicou.legal_area la 
		ON la.id = cf.legal_area_id

	LEFT JOIN sicou.service_type st 
		ON st.id = cf.service_type_id

	LEFT JOIN sicou.app_user u 
		ON u.id = cf.current_responsible

	LEFT JOIN sicou.turn t 
		ON t.id = cf.turn_id

	LEFT JOIN sicou.preturno pt 
		ON pt.id = cf.preturno_id

	LEFT JOIN sicou.app_user su 
		ON su.id = cf.supervisor_user

	LEFT JOIN (
		SELECT 
			case_id,
			COUNT(*) AS attachments_count,
			COUNT(*) FILTER (
				WHERE notes IS NOT NULL 
				AND TRIM(notes) <> ''
			) AS comments_count
		FROM sicou.document
		WHERE deleted_at IS NULL
		AND is_current = true
		AND case_id IS NOT NULL
		GROUP BY case_id
	) att 
		ON att.case_id = cf.id

	LEFT JOIN (
		SELECT
			d.case_id,
			COUNT(DISTINCT d.document_kind_id) FILTER (
				WHERE dk.code <> 'OTRO'
				AND dk.is_active = true
			) AS completed_kinds
		FROM sicou.document d
		JOIN sicou.document_kind dk 
			ON dk.id = d.document_kind_id
		WHERE d.deleted_at IS NULL
		AND d.is_current = true
		AND d.case_id IS NOT NULL
		GROUP BY d.case_id
	) prog 
		ON prog.case_id = cf.id

	CROSS JOIN (
		SELECT COUNT(*) AS total_required
		FROM sicou.document_kind
		WHERE code <> 'OTRO'
		AND is_active = true
	) req

	WHERE cf.citizen_id = $1

	ORDER BY cf.created_at DESC;
	`

	rows, err := r.db.Query(ctx, query, citizenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []dto.CaseListItem

	for rows.Next() {
		var item dto.CaseListItem

		err := rows.Scan(
			&item.ExpedienteCode,
			&item.ID,
			&item.Status,
			&item.OpenedAt,
			&item.ClosedAt,
			&item.SupervisorID,
			&item.SupervisorName,
			&item.FullName,
			&item.DocumentNumber,
			&item.LegalArea,
			&item.ServiceType,
			&item.ResponsibleName,
			&item.Consecutive,
			&item.TurnDate,
			&item.AttachmentsCount,
			&item.CommentsCount,
			&item.Progress,
		)

		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return list, nil
}

func (r *SCaseFile) GetCitizenIDByEmail(
	ctx context.Context,
	email string,
) (string, error) {

	query := `
		SELECT id
		FROM sicou.citizen
		WHERE email = $1
		AND deleted_at IS NULL
		LIMIT 1
	`

	var citizenID string

	err := r.db.QueryRow(ctx, query, email).Scan(&citizenID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("ciudadano no encontrado")
		}
		return "", err
	}

	return citizenID, nil
}

func (r *SCaseFile) SaveOtp(
	ctx context.Context,
	citizenID string,
	otpHash string,
	expiresAt time.Time,
) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1️⃣ Invalidar OTP activos anteriores
	invalidateQuery := `
		UPDATE sicou.case_file_otp
		SET used_at = now()
		WHERE citizen_id = $1
		AND used_at IS NULL
		AND expires_at > now()
	`

	_, err = tx.Exec(ctx, invalidateQuery, citizenID)
	if err != nil {
		return err
	}

	// 2️⃣ Insertar nuevo OTP
	insertQuery := `
		INSERT INTO sicou.case_file_otp (
			citizen_id,
			otp_hash,
			expires_at
		)
		VALUES ($1,$2,$3)
	`

	_, err = tx.Exec(ctx, insertQuery,
		citizenID,
		otpHash,
		expiresAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *SCaseFile) InsertEmailOutbox(
	ctx context.Context,
	to []string,
	subject string,
	body string,
) error {

	query := `
		INSERT INTO sicou.email_outbox (
			to_emails,
			subject,
			body
		)
		VALUES ($1,$2,$3)
	`

	_, err := r.db.Exec(ctx, query,
		to,
		subject,
		body,
	)

	return err
}

func (r *SCaseFile) VerifyOtp(
	ctx context.Context,
	citizenID string,
	otp string,
) error {

	query := `
		SELECT id, otp_hash
		FROM sicou.case_file_otp
		WHERE citizen_id = $1
		AND used_at IS NULL
		AND expires_at > now()
		ORDER BY expires_at DESC
		LIMIT 1
	`

	var (
		id      string
		otpHash string
	)

	err := r.db.QueryRow(ctx, query, citizenID).Scan(&id, &otpHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("otp inválido o expirado")
		}
		return err
	}

	// Comparar hash
	err = bcrypt.CompareHashAndPassword([]byte(otpHash), []byte(otp))
	if err != nil {
		return errors.New("otp incorrecto")
	}

	// Marcar como usado
	updateQuery := `
		UPDATE sicou.case_file_otp
		SET used_at = now()
		WHERE id = $1
	`

	_, err = r.db.Exec(ctx, updateQuery, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SCaseFile) FolderList(ctx context.Context) ([]dto.FolderResponse, error) {

	query := `
		SELECT
			id,
			name,
			description,
			created_at
		FROM sicou.folder
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []dto.FolderResponse

	for rows.Next() {
		var item dto.FolderResponse

		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return list, nil
}
