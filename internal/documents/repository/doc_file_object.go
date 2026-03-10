package repository

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/documents/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDocumentNotFound = errors.New("document not found")

type IDocFileObject interface {
	UploadDoc(ctx context.Context, doc dto.CreateDocFileObjectRequest) (*dto.Document, error)
	GetDocumentKinds(ctx context.Context) ([]dto.DocumentKind, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	GetCaseTimeLineItem(
		ctx context.Context,
		caseID uuid.UUID,
	) ([]dto.CaseTimelineItem, error)
	ListCaseFiles(
		ctx context.Context,
		caseID uuid.UUID,
	) ([]dto.CaseDocumentListItem, error)

	GetFileMetadata(
		ctx context.Context,
		fileID uuid.UUID,
	) (*dto.FileMetadata, error)

	UpdateDocument(
		ctx context.Context,
		documentID uuid.UUID,
		req dto.UpdateDocumentRequest,
	) error

	DeleteDocument(
		ctx context.Context,
		documentID uuid.UUID,
	) error

	UploadDocFolder(ctx context.Context,
		req dto.CreateDocFileObjectFolderRequest) (*dto.Document, error)
}

type SDocFileObjectRepository struct {
	db *pgxpool.Pool
}

func NewDocFileObjectRepository(db *pgxpool.Pool) *SDocFileObjectRepository {
	return &SDocFileObjectRepository{db: db}
}

func (r *SDocFileObjectRepository) UploadDoc(
	ctx context.Context,
	req dto.CreateDocFileObjectRequest,
) (*dto.Document, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 1️⃣ Insert file_object
	fileQuery := `
		INSERT INTO sicou.file_object
		(id, storage_key, original_name, mime_type, size_bytes, sha256, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,now())
	`

	_, err = tx.Exec(ctx, fileQuery,
		req.FileID,
		req.StorageKey,
		req.OriginalName,
		req.MimeType,
		req.SizeBytes,
		req.Sha256,
	)
	if err != nil {
		return nil, err
	}

	// 2️⃣ Insert document
	docQuery := `
		INSERT INTO sicou.document
		(file_id, document_kind_id, preturno_id, case_id, case_event_id,
		 version, is_current, notes, uploaded_by, uploaded_at)
		VALUES ($1,$2,$3,$4,$5,1,true,$6,$7,now())
		RETURNING id, file_id, document_kind_id, preturno_id, case_id,
		          case_event_id, version, is_current, notes,
		          uploaded_by, uploaded_at, deleted_at
	`

	var doc dto.Document

	err = tx.QueryRow(ctx, docQuery,
		req.FileID,
		req.DocumentKindID,
		req.PreturnoID,
		req.CaseID,
		req.CaseEventID,
		req.Notes,
		req.UploadedBy,
	).Scan(
		&doc.ID,
		&doc.FileID,
		&doc.DocumentKindID,
		&doc.PreturnoID,
		&doc.CaseID,
		&doc.CaseEventID,
		&doc.Version,
		&doc.IsCurrent,
		&doc.Notes,
		&doc.UploadedBy,
		&doc.UploadedAt,
		&doc.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &doc, nil
}

func (r *SDocFileObjectRepository) GetDocumentKinds(
	ctx context.Context,
) ([]dto.DocumentKind, error) {

	query := `
		SELECT id, code, name, is_active, created_at, updated_at
		FROM sicou.document_kind
		WHERE is_active = true
		ORDER BY id ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kinds []dto.DocumentKind

	for rows.Next() {
		var k dto.DocumentKind

		err := rows.Scan(
			&k.ID,
			&k.Code,
			&k.Name,
			&k.IsActive,
			&k.CreatedAt,
			&k.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		kinds = append(kinds, k)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return kinds, nil
}

func (r *SDocFileObjectRepository) GetCaseTimeLineItem(
	ctx context.Context,
	caseID uuid.UUID,
) ([]dto.CaseTimelineItem, error) {

	query := `
	SELECT
		d.id,
		d.version,
		d.is_current,
		d.uploaded_at,
		d.deleted_at,
		d.notes,

		dk.name AS document_kind,
		f.original_name,
		f.size_bytes,
		f.mime_type,

		u.display_name AS uploaded_by_name

	FROM sicou.document d
	JOIN sicou.file_object f ON f.id = d.file_id
	JOIN sicou.document_kind dk ON dk.id = d.document_kind_id
	LEFT JOIN sicou.app_user u ON u.id = d.uploaded_by

	WHERE d.case_id = $1

	ORDER BY COALESCE(d.deleted_at, d.uploaded_at) DESC
	`

	rows, err := r.db.Query(ctx, query, caseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]dto.CaseTimelineItem, 0)

	for rows.Next() {

		var item dto.CaseTimelineItem

		err := rows.Scan(
			&item.DocumentID,
			&item.Version,
			&item.IsCurrent,
			&item.UploadedAt,
			&item.DeletedAt,
			&item.Notes,
			&item.DocumentKind,
			&item.OriginalName,
			&item.SizeBytes,
			&item.MimeType,
			&item.UploadedByName,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *SDocFileObjectRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM sicou.case_file WHERE id = $1`

	var exists int
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)

	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *SDocFileObjectRepository) ListCaseFiles(
	ctx context.Context,
	caseID uuid.UUID,
) ([]dto.CaseDocumentListItem, error) {

	query := `
SELECT
	d.id,
	f.id,
	dk.name,
	f.original_name,
	f.size_bytes,
	f.mime_type,
	d.version,
	d.uploaded_at,
	u.id,
	u.display_name,

	/* Carpeta */
	STRING_AGG(DISTINCT fol.name, ', ') AS folder_name

FROM sicou.document d

JOIN sicou.file_object f 
	ON f.id = d.file_id

JOIN sicou.document_kind dk 
	ON dk.id = d.document_kind_id

LEFT JOIN sicou.app_user u 
	ON u.id = d.uploaded_by

LEFT JOIN sicou.folder_document fd
	ON fd.document_id = d.id

LEFT JOIN sicou.folder fol
	ON fol.id = fd.folder_id
	AND fol.deleted_at IS NULL

WHERE d.case_id = $1
AND d.deleted_at IS NULL

GROUP BY
	d.id,
	f.id,
	dk.name,
	f.original_name,
	f.size_bytes,
	f.mime_type,
	d.version,
	d.uploaded_at,
	u.id,
	u.display_name

ORDER BY d.uploaded_at DESC;
	`

	rows, err := r.db.Query(ctx, query, caseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]dto.CaseDocumentListItem, 0) // 👈 MEJOR PRÁCTICA

	for rows.Next() {
		var item dto.CaseDocumentListItem

		err := rows.Scan(
			&item.DocumentID,
			&item.FileID,
			&item.DocumentKind,
			&item.OriginalName,
			&item.SizeBytes,
			&item.MimeType,
			&item.Version,
			&item.UploadedAt,
			&item.ResponsibleID,
			&item.ResponsibleName,
			&item.FolderName,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *SDocFileObjectRepository) UpdateDocument(
	ctx context.Context,
	documentID uuid.UUID,
	req dto.UpdateDocumentRequest,
) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var fileID uuid.UUID

	// 1️⃣ Validar existencia
	err = tx.QueryRow(ctx, `
		SELECT file_id
		FROM sicou.document
		WHERE id = $1
		  AND deleted_at IS NULL
		  AND is_current = true
	`, documentID).Scan(&fileID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrDocumentNotFound
		}
		return err
	}

	// 2️⃣ Actualizar tipo si viene
	if req.DocumentKindID != nil {

		// Validación defensiva
		if *req.DocumentKindID <= 0 {
			return errors.New("invalid document_kind_id")
		}

		cmd, err := tx.Exec(ctx, `
			UPDATE sicou.document
			SET document_kind_id = $1
			WHERE id = $2
		`, *req.DocumentKindID, documentID)
		if err != nil {
			return err
		}

		if cmd.RowsAffected() == 0 {
			return ErrDocumentNotFound
		}
	}

	// 3️⃣ Actualizar nombre si viene
	if req.OriginalName != nil {

		newName := strings.TrimSpace(*req.OriginalName)
		if newName == "" {
			return errors.New("original_name cannot be empty")
		}

		// 1️⃣ Obtener nombre actual
		var currentName string
		err := tx.QueryRow(ctx, `
		SELECT original_name
		FROM sicou.file_object
		WHERE id = $1
	`, fileID).Scan(&currentName)
		if err != nil {
			return err
		}

		// 2️⃣ Extraer extensión actual
		currentExt := filepath.Ext(currentName)

		// 3️⃣ Si el nuevo nombre NO trae extensión, agregar la actual
		if filepath.Ext(newName) == "" && currentExt != "" {
			newName = newName + currentExt
		}

		// 4️⃣ Update
		cmd, err := tx.Exec(ctx, `
		UPDATE sicou.file_object
		SET original_name = $1
		WHERE id = $2
	`, newName, fileID)
		if err != nil {
			return err
		}

		if cmd.RowsAffected() == 0 {
			return ErrDocumentNotFound
		}
	}

	// 4️⃣ Commit
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *SDocFileObjectRepository) GetFileMetadata(
	ctx context.Context,
	fileID uuid.UUID,
) (*dto.FileMetadata, error) {

	query := `
	SELECT
		id,
		storage_key,
		original_name,
		mime_type
	FROM sicou.file_object
	WHERE id = $1
	`

	var meta dto.FileMetadata

	err := r.db.QueryRow(ctx, query, fileID).Scan(
		&meta.ID,
		&meta.StorageKey,
		&meta.OriginalName,
		&meta.MimeType,
	)

	if err != nil {
		return nil, err
	}

	return &meta, nil
}

func (r *SDocFileObjectRepository) DeleteDocument(
	ctx context.Context,
	documentID uuid.UUID,
) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	cmd, err := tx.Exec(ctx, `
		UPDATE sicou.document
		SET deleted_at = NOW(),
		    is_current = false
		WHERE id = $1
		  AND deleted_at IS NULL
	`, documentID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrDocumentNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *SDocFileObjectRepository) UploadDocFolder(ctx context.Context,
	req dto.CreateDocFileObjectFolderRequest) (*dto.Document, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 1️⃣ Insert file_object
	fileQuery := `
		INSERT INTO sicou.file_object
		(id, storage_key, original_name, mime_type, size_bytes, sha256, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,now())
	`

	_, err = tx.Exec(ctx, fileQuery,
		req.FileID,
		req.StorageKey,
		req.OriginalName,
		req.MimeType,
		req.SizeBytes,
		req.Sha256,
	)
	if err != nil {
		return nil, err
	}

	// 2️⃣ Insert document
	docQuery := `
		INSERT INTO sicou.document
		(file_id, document_kind_id, preturno_id, case_id, case_event_id,
		 version, is_current, notes, uploaded_by, uploaded_at)
		VALUES ($1,$2,$3,$4,$5,1,true,$6,$7,now())
		RETURNING id, file_id, document_kind_id, preturno_id, case_id,
		          case_event_id, version, is_current, notes,
		          uploaded_by, uploaded_at, deleted_at
	`

	var doc dto.Document

	err = tx.QueryRow(ctx, docQuery,
		req.FileID,
		req.DocumentKindID,
		req.PreturnoID,
		req.CaseID,
		req.CaseEventID,
		req.Notes,
		req.UploadedBy,
	).Scan(
		&doc.ID,
		&doc.FileID,
		&doc.DocumentKindID,
		&doc.PreturnoID,
		&doc.CaseID,
		&doc.CaseEventID,
		&doc.Version,
		&doc.IsCurrent,
		&doc.Notes,
		&doc.UploadedBy,
		&doc.UploadedAt,
		&doc.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Relacionar con carpeta
	folderDocQuery := `
		INSERT INTO sicou.folder_document
		(folder_id, document_id, added_at, added_by)
		VALUES ($1, $2, now(), $3)
	`

	_, err = tx.Exec(ctx, folderDocQuery,
		req.FolderID,
		doc.ID,
		req.UploadedBy,
	)
	if err != nil {
		return nil, err
	}

	// 4️⃣ Commit
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &doc, nil
}
