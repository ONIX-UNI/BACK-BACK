package services

import (
	"context"
	"errors"

	"github.com/DuvanRozoParra/sicou/internal/documents/dto"
	"github.com/DuvanRozoParra/sicou/internal/documents/repository"
)

type DocumentService struct {
	repo repository.DocumentRepository
}

func NewDocumentService(repo repository.DocumentRepository) *DocumentService {
	return &DocumentService{repo: repo}
}

func (s *DocumentService) Create(
	ctx context.Context,
	req dto.CreateDocumentRequest,
) (*dto.Document, error) {

	if req.FileID == "" {
		return nil, errors.New("file_id is required")
	}

	if req.DocumentKindID == 0 {
		return nil, errors.New("document_kind_id is required")
	}

	return s.repo.Create(ctx, req)
}

func (s *DocumentService) GetByID(
	ctx context.Context,
	id string,
) (*dto.Document, error) {

	if id == "" {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *DocumentService) GetByFileID(
	ctx context.Context,
	fileID string,
) ([]dto.Document, error) {

	if fileID == "" {
		return nil, errors.New("file_id is required")
	}

	return s.repo.GetByFileID(ctx, fileID)
}

func (s *DocumentService) GetByCaseID(
	ctx context.Context,
	caseID string,
) ([]dto.Document, error) {

	if caseID == "" {
		return nil, errors.New("case_id is required")
	}

	return s.repo.GetByCaseID(ctx, caseID)
}

func (s *DocumentService) Delete(
	ctx context.Context,
	id string,
) error {

	if id == "" {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}
