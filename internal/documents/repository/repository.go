package repository

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/documents/dto"
)

type FileObjectRepository interface {
	Create(ctx context.Context, file dto.FileObject) (*dto.FileObject, error)
	UpdateFull(ctx context.Context, file dto.FileObject) (*dto.FileObject, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]dto.FileObject, error)
	GetById(ctx context.Context, id string) (*dto.FileObject, error)
	GetByName(ctx context.Context, name string) (*dto.FileObject, error)
	List(ctx context.Context, limit, offset int) ([]dto.FileObject, error)
}

type DocumentRepository interface {
	Create(ctx context.Context, req dto.CreateDocumentRequest) (*dto.Document, error)
	GetByID(ctx context.Context, id string) (*dto.Document, error)
	GetByFileID(ctx context.Context, fileID string) ([]dto.Document, error)
	GetByCaseID(ctx context.Context, caseID string) ([]dto.Document, error)
	Delete(ctx context.Context, id string) error
}
