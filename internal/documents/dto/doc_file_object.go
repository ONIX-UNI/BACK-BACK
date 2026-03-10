package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDocFileObjectRequest struct {
	// 🔹 Interno (no viene del cliente)
	FileID string `json:"-"`

	// 🔹 Document
	DocumentKindID int16   `json:"document_kind_id" validate:"required"`
	PreturnoID     *string `json:"preturno_id,omitempty"`
	CaseID         *string `json:"case_id,omitempty"`
	CaseEventID    *string `json:"case_event_id,omitempty"`
	Notes          *string `json:"notes,omitempty"`

	// 🔹 File
	StorageKey   string `json:"storage_key" validate:"required"`
	OriginalName string `json:"original_name" validate:"required"`
	MimeType     string `json:"mime_type,omitempty"`
	SizeBytes    int64  `json:"size_bytes" validate:"required"`
	Sha256       string `json:"sha256,omitempty"`

	// 🔹 Auditoría
	UploadedBy *string `json:"uploaded_by,omitempty"`
}

type DocumentKind struct {
	ID        int16     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CaseTimelineItem struct {
	DocumentID uuid.UUID  `json:"document_id"`
	Version    int        `json:"version"`
	IsCurrent  bool       `json:"is_current"`
	UploadedAt time.Time  `json:"uploaded_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Notes      *string    `json:"notes,omitempty"`

	DocumentKind   string  `json:"document_kind"`
	OriginalName   string  `json:"original_name"`
	SizeBytes      int64   `json:"size_bytes"`
	MimeType       *string `json:"mime_type,omitempty"`
	UploadedByName *string `json:"uploaded_by_name,omitempty"`
}

type CaseDocumentListItem struct {
	DocumentID      uuid.UUID  `json:"document_id"`
	FileID          uuid.UUID  `json:"file_id"`
	DocumentKind    string     `json:"document_kind"`
	OriginalName    string     `json:"original_name"`
	SizeBytes       int64      `json:"size_bytes"`
	MimeType        *string    `json:"mime_type,omitempty"`
	Version         int        `json:"version"`
	UploadedAt      time.Time  `json:"uploaded_at"`
	ResponsibleID   *uuid.UUID `json:"responsible_id,omitempty"`
	ResponsibleName *string    `json:"responsible_name,omitempty"`

	FolderName *string `json:"folder_name,omitempty"`
}

type FileMetadata struct {
	ID           uuid.UUID
	StorageKey   string
	OriginalName string
	MimeType     *string
}

type UpdateDocumentRequest struct {
	OriginalName   *string `json:"original_name,omitempty"`
	DocumentKindID *int16  `json:"document_kind_id,omitempty"`
}

type CreateDocFileObjectFolderRequest struct {
	FileID       uuid.UUID
	StorageKey   string
	OriginalName string
	MimeType     string
	SizeBytes    int64
	Sha256       string

	DocumentKindID int16
	PreturnoID     *uuid.UUID
	CaseID         *uuid.UUID
	CaseEventID    *uuid.UUID
	Notes          *string
	UploadedBy     uuid.UUID

	FolderID uuid.UUID
}

type FolderWithDocuments struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description,omitempty"`
	Documents   []DocumentSimple `json:"documents"`
}

type DocumentSimple struct {
	ID             uuid.UUID  `json:"id"`
	FileID         uuid.UUID  `json:"file_id"`
	DocumentKindID int16      `json:"document_kind_id"`
	Notes          *string    `json:"notes,omitempty"`
	UploadedAt     time.Time  `json:"uploaded_at"`
	UploadedBy     *uuid.UUID `json:"uploaded_by,omitempty"`
}
