package repository

import (
	"context"
	"errors"

	"github.com/DuvanRozoParra/sicou/internal/documents/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDocumentRepository struct {
	db *pgxpool.Pool
}

func NewPostgresDocumentRepository(db *pgxpool.Pool) *PostgresDocumentRepository {
	return &PostgresDocumentRepository{db: db}
}

func (r *PostgresDocumentRepository) Create(
	ctx context.Context,
	req dto.CreateDocumentRequest,
) (*dto.Document, error) {

	query := `
		INSERT INTO sicou.document
		(file_id, document_kind_id, preturno_id, case_id, case_event_id,
		 version, is_current, notes, uploaded_by, uploaded_at)
		VALUES ($1,$2,$3,$4,$5,1,true,$6,$7,now())
		RETURNING id, file_id, document_kind_id, preturno_id, case_id,
		          case_event_id, version, is_current, notes,
		          uploaded_by, uploaded_at, deleted_at
	`

	row := r.db.QueryRow(ctx, query,
		req.FileID,
		req.DocumentKindID,
		req.PreturnoID,
		req.CaseID,
		req.CaseEventID,
		req.Notes,
		req.UploadedBy,
	)

	var doc dto.Document

	err := row.Scan(
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

	return &doc, nil
}

func (r *PostgresDocumentRepository) GetByID(
	ctx context.Context,
	id string,
) (*dto.Document, error) {

	if id == "" {
		return nil, errors.New("document ID is required")
	}

	query := `
		SELECT id, file_id, document_kind_id, preturno_id, case_id,
		       case_event_id, version, is_current, notes,
		       uploaded_by, uploaded_at, deleted_at
		FROM sicou.document
		WHERE id=$1 AND deleted_at IS NULL
	`

	row := r.db.QueryRow(ctx, query, id)

	var doc dto.Document

	err := row.Scan(
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

	return &doc, nil
}

func (r *PostgresDocumentRepository) GetByFileID(
	ctx context.Context,
	fileID string,
) ([]dto.Document, error) {

	if fileID == "" {
		return nil, errors.New("file ID is required")
	}

	query := `
		SELECT id, file_id, document_kind_id, preturno_id, case_id,
		       case_event_id, version, is_current, notes,
		       uploaded_by, uploaded_at, deleted_at
		FROM sicou.document
		WHERE file_id=$1 AND deleted_at IS NULL
		ORDER BY version DESC
	`

	rows, err := r.db.Query(ctx, query, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []dto.Document

	for rows.Next() {
		var doc dto.Document

		err := rows.Scan(
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

		documents = append(documents, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *PostgresDocumentRepository) GetByCaseID(
	ctx context.Context,
	caseID string,
) ([]dto.Document, error) {

	if caseID == "" {
		return nil, errors.New("case ID is required")
	}

	query := `
		SELECT id, file_id, document_kind_id, preturno_id, case_id,
		       case_event_id, version, is_current, notes,
		       uploaded_by, uploaded_at, deleted_at
		FROM sicou.document
		WHERE case_id=$1 AND deleted_at IS NULL
		ORDER BY uploaded_at DESC
	`

	rows, err := r.db.Query(ctx, query, caseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []dto.Document

	for rows.Next() {
		var doc dto.Document

		err := rows.Scan(
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

		documents = append(documents, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *PostgresDocumentRepository) Delete(
	ctx context.Context,
	id string,
) error {

	if id == "" {
		return errors.New("document ID is required")
	}

	// Soft delete
	query := `
		UPDATE sicou.document
		SET deleted_at = now(), is_current = false
		WHERE id=$1 AND deleted_at IS NULL
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return nil
	}

	return nil
}
