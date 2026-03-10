package services

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/documents/dto"
	"github.com/DuvanRozoParra/sicou/internal/documents/repository"
	"github.com/DuvanRozoParra/sicou/pkg/storage"
	"github.com/google/uuid"
)

type FileObjectService struct {
	repo        repository.FileObjectRepository
	minioClient *storage.Client
	bucket      string
}

func NewFileObjectService(
	repo repository.FileObjectRepository,
	minioClient *storage.Client,
	bucket string,
) *FileObjectService {
	return &FileObjectService{
		repo:        repo,
		minioClient: minioClient,
		bucket:      bucket,
	}
}

func (s *FileObjectService) Create(
	ctx context.Context,
	file multipart.File,
	size int64,
	filename string,
	contentType string,
) (*dto.FileObject, error) {

	if file == nil {
		return nil, errors.New("file is required")
	}

	// 1️⃣ Generar ID para file_object
	fileID := uuid.New().String()

	// 2️⃣ Generar nombre único para MinIO
	ext := filepath.Ext(filename)
	objectName := fileID + ext

	// 3️⃣ Subir a MinIO
	_, err := s.minioClient.PutObject(
		ctx,
		s.bucket,
		objectName,
		file,
		size,
		storage.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return nil, err
	}

	// 4️⃣ Crear metadata para DB
	fileObj := dto.FileObject{
		ID:           fileID,
		StorageKey:   objectName,
		OriginalName: filename,
		MimeType:     contentType,
		SizeBytes:    size,
		CreatedAt:    time.Now(),
	}

	// 5️⃣ Guardar en DB
	return s.repo.Create(ctx, fileObj)
}

func (s *FileObjectService) Update(
	ctx context.Context,
	id string,
	file multipart.File,
	size int64,
	filename string,
	contentType string,
	req dto.UpdateFileObjectRequest,
) (*dto.FileObject, error) {

	if id == "" {
		return nil, errors.New("id is required")
	}

	existing, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	var oldStorageKey string

	if file != nil {

		ext := filepath.Ext(filename)
		newObjectName := id + ext

		// 1️⃣ subir nuevo primero
		_, err = s.minioClient.PutObject(
			ctx,
			s.bucket,
			newObjectName,
			file,
			size,
			storage.PutObjectOptions{
				ContentType: contentType,
			},
		)
		if err != nil {
			return nil, err
		}

		oldStorageKey = existing.StorageKey

		existing.StorageKey = newObjectName
		existing.SizeBytes = size
		existing.MimeType = contentType
		existing.OriginalName = filename
	}

	if req.OriginalName != nil {
		existing.OriginalName = *req.OriginalName
	}

	updated, err := s.repo.UpdateFull(ctx, *existing)
	if err != nil {

		// rollback del nuevo upload si DB falla
		if file != nil {
			_ = s.minioClient.RemoveObject(
				ctx,
				s.bucket,
				existing.StorageKey,
				storage.RemoveObjectOptions{},
			)
		}

		return nil, err
	}

	// 3️⃣ eliminar el anterior SOLO si todo salió bien
	if file != nil && oldStorageKey != "" && oldStorageKey != existing.StorageKey {
		_ = s.minioClient.RemoveObject(
			ctx,
			s.bucket,
			oldStorageKey,
			storage.RemoveObjectOptions{},
		)
	}

	return updated, nil
}

func (s *FileObjectService) GetById(
	ctx context.Context,
	id string,
) (*dto.FileObject, error) {

	return s.repo.GetById(ctx, id)
}

func (s *FileObjectService) GetAll(ctx context.Context) ([]dto.FileObject, error) {
	return s.repo.GetAll(ctx)
}

func (s *FileObjectService) Delete(
	ctx context.Context,
	id string,
) error {

	if id == "" {
		return errors.New("id is required")
	}

	existing, err := s.repo.GetById(ctx, id)
	if err != nil {
		return err
	}

	// Intentar borrar objeto (no bloquea si falla)
	_ = s.minioClient.RemoveObject(
		ctx,
		s.bucket,
		existing.StorageKey,
		storage.RemoveObjectOptions{},
	)

	return s.repo.Delete(ctx, id)
}

func (s *FileObjectService) GetFileStream(
	ctx context.Context,
	id string,
) (*storage.Object, *dto.FileObject, error) {

	if id == "" {
		return nil, nil, errors.New("id is required")
	}

	fileObj, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	object, err := s.minioClient.GetObject(
		ctx,
		s.bucket,
		fileObj.StorageKey,
		storage.GetObjectOptions{},
	)
	if err != nil {
		return nil, nil, err
	}

	return object, fileObj, nil
}

func (s *FileObjectService) GetPresignedURL(
	ctx context.Context,
	id string,
) (string, error) {

	fileObj, err := s.repo.GetById(ctx, id)
	if err != nil {
		return "", err
	}

	url, err := s.minioClient.PresignedGetObject(
		ctx,
		s.bucket,
		fileObj.StorageKey,
		time.Minute*15,
		nil,
	)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (s *FileObjectService) ListDocuments(
	ctx context.Context,
	page int,
	limit int,
) ([]dto.FileObject, error) {

	if page < 1 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	return s.repo.List(ctx, limit, offset)
}
