package services

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/documents/dto"
	"github.com/DuvanRozoParra/sicou/internal/documents/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/minio/minio-go/v7"
)

var (
	ErrDuplicateDocument = errors.New("duplicate_document_kind_for_case")
	ErrValidation        = errors.New("validation_error")

	ErrCaseNotFound  = errors.New("case_not_found")
	ErrInvalidCaseID = errors.New("invalid_case_id")

	ErrDocumentNotFound = errors.New("document not found")
)

type SDocFileObjectService struct {
	repo        repository.IDocFileObject
	minioClient *minio.Client
	bucket      string
}

func NewDocFileObjectService(repo repository.IDocFileObject, minioClient *minio.Client, bucket string) *SDocFileObjectService {
	return &SDocFileObjectService{repo: repo, minioClient: minioClient, bucket: bucket}
}

func (s *SDocFileObjectService) UploadDoc(
	ctx context.Context,
	file multipart.File,
	size int64,
	filename string,
	contentType string,
	req dto.CreateDocFileObjectRequest,
) (*dto.Document, error) {

	// 🔹 1. Validaciones del archivo
	if file == nil {
		return nil, errors.New("file is required")
	}

	if size <= 0 {
		return nil, errors.New("file size must be greater than zero")
	}

	if filename == "" {
		return nil, errors.New("filename is required")
	}

	if contentType == "" {
		return nil, errors.New("content type is required")
	}

	// 🔹 2. Validaciones del request
	if req.DocumentKindID == 0 {
		return nil, errors.New("document_kind_id is required")
	}

	// Al menos uno debe existir
	if req.PreturnoID == nil && req.CaseID == nil && req.CaseEventID == nil {
		return nil, errors.New("at least one of preturno_id, case_id or case_event_id is required")
	}

	// Si quieres hacerlo obligatorio:
	if req.UploadedBy == nil || *req.UploadedBy == "" {
		return nil, errors.New("uploaded_by is required")
	}

	// 🔹 3. Generar identificador trazable
	fileID := uuid.New().String()
	ext := filepath.Ext(filename)
	objectName := fileID + ext

	// 🔹 4. Subir a MinIO
	_, err := s.minioClient.PutObject(
		ctx,
		s.bucket,
		objectName,
		file,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return nil, err
	}

	// 🔹 5. Completar request interno
	req.FileID = fileID
	req.SizeBytes = size
	req.MimeType = contentType
	req.StorageKey = objectName
	req.OriginalName = filename

	// 🔹 6. Persistir en DB
	doc, err := s.repo.UploadDoc(ctx, req)
	if err != nil {

		// rollback archivo en MinIO
		_ = s.minioClient.RemoveObject(
			ctx,
			s.bucket,
			objectName,
			minio.RemoveObjectOptions{},
		)

		// 🔎 Detectar error de constraint única
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" &&
				pgErr.ConstraintName == "uq_case_doc_kind_current" {

				return nil, ErrDuplicateDocument
			}
		}

		return nil, err
	}

	return doc, nil
}

func (s *SDocFileObjectService) UploadDocFolder(
	ctx context.Context,
	file multipart.File,
	size int64,
	filename string,
	contentType string,
	req dto.CreateDocFileObjectFolderRequest,
) (*dto.Document, error) {

	// 1️⃣ Validaciones archivo
	if file == nil {
		return nil, errors.New("file is required")
	}

	if size <= 0 {
		return nil, errors.New("file size must be greater than zero")
	}

	if filename == "" {
		return nil, errors.New("filename is required")
	}

	if contentType == "" {
		return nil, errors.New("content type is required")
	}

	// 2️⃣ Validaciones request
	if req.DocumentKindID == 0 {
		return nil, errors.New("document_kind_id is required")
	}

	if req.FolderID == uuid.Nil {
		return nil, errors.New("folder_id is required")
	}

	if req.PreturnoID == nil && req.CaseID == nil && req.CaseEventID == nil {
		return nil, errors.New("at least one of preturno_id, case_id or case_event_id is required")
	}

	if req.UploadedBy == uuid.Nil {
		return nil, errors.New("uploaded_by is required")
	}

	// 3️⃣ Generar UUID real
	fileID := uuid.New()
	ext := filepath.Ext(filename)
	objectName := fileID.String() + ext

	// 4️⃣ Subir a MinIO
	_, err := s.minioClient.PutObject(
		ctx,
		s.bucket,
		objectName,
		file,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return nil, err
	}

	// 5️⃣ Completar request
	req.FileID = fileID
	req.SizeBytes = size
	req.MimeType = contentType
	req.StorageKey = objectName
	req.OriginalName = filename

	// 6️⃣ Persistir en DB
	doc, err := s.repo.UploadDocFolder(ctx, req)
	if err != nil {

		// rollback MinIO
		_ = s.minioClient.RemoveObject(
			ctx,
			s.bucket,
			objectName,
			minio.RemoveObjectOptions{},
		)

		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" &&
				pgErr.ConstraintName == "uq_case_doc_kind_current" {
				return nil, ErrDuplicateDocument
			}
		}

		return nil, err
	}

	return doc, nil
}

func (r *SDocFileObjectService) GetDocumentKinds(
	ctx context.Context,
) ([]dto.DocumentKind, error) {
	return r.repo.GetDocumentKinds(ctx)
}

func (s *SDocFileObjectService) GetCaseTimeLineItem(
	ctx context.Context,
	caseID uuid.UUID,
) ([]dto.CaseTimelineItem, error) {

	if caseID == uuid.Nil {
		return nil, errors.New("case_id is required")
	}

	exists, err := s.repo.Exists(ctx, caseID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("case not found")
	}

	return s.repo.GetCaseTimeLineItem(ctx, caseID)
}

func (s *SDocFileObjectService) ListCaseFiles(
	ctx context.Context,
	caseID uuid.UUID,
) ([]dto.CaseDocumentListItem, error) {

	if caseID == uuid.Nil {
		return nil, errors.New("case_id is required")
	}

	exists, err := s.repo.Exists(ctx, caseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("case not found")
	}

	return s.repo.ListCaseFiles(ctx, caseID)
}

func (s *SDocFileObjectService) GetFileStream(
	ctx context.Context,
	fileID uuid.UUID,
) (*dto.FileMetadata, *minio.Object, error) {

	meta, err := s.repo.GetFileMetadata(ctx, fileID)
	if err != nil {
		return nil, nil, err
	}

	object, err := s.minioClient.GetObject(
		ctx,
		"expedient-docs",
		meta.StorageKey,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, nil, err
	}

	return meta, object, nil
}

func (s *SDocFileObjectService) UpdateDocument(
	ctx context.Context,
	documentID uuid.UUID,
	req dto.UpdateDocumentRequest,
) error {

	// 1️⃣ Validar documentID
	if documentID == uuid.Nil {
		return errors.New("document_id is required")
	}

	// 2️⃣ Validar que al menos un campo venga en el PATCH
	// if req.OriginalName == nil && req.DocumentKindID == nil {
	// 	return errors.New("at least one field must be provided")
	// }

	// 3️⃣ Validación básica del nombre (si viene)
	if req.OriginalName != nil {
		if strings.TrimSpace(*req.OriginalName) == "" {
			return errors.New("original_name cannot be empty")
		}
	}

	// 4️⃣ Delegar al repository
	err := s.repo.UpdateDocument(ctx, documentID, req)
	if err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			return ErrDocumentNotFound
		}
		return err
	}

	return nil
}

func (s *SDocFileObjectService) DeleteDocument(
	ctx context.Context,
	documentID uuid.UUID,
) error {

	if documentID == uuid.Nil {
		return errors.New("document_id is required")
	}

	err := s.repo.DeleteDocument(ctx, documentID)
	if err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			return ErrDocumentNotFound
		}
		return err
	}

	return nil
}
