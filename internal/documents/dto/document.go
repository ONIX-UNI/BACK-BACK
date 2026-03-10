package dto

import "time"

type Document struct {
	ID             string     `json:"id"`
	FileID         string     `json:"file_id"`
	DocumentKindID int16      `json:"document_kind_id"`
	PreturnoID     *string    `json:"preturno_id,omitempty"`
	CaseID         *string    `json:"case_id,omitempty"`
	CaseEventID    *string    `json:"case_event_id,omitempty"`
	Version        int        `json:"version"`
	IsCurrent      bool       `json:"is_current"`
	Notes          *string    `json:"notes,omitempty"`
	UploadedBy     *string    `json:"uploaded_by,omitempty"`
	UploadedAt     time.Time  `json:"uploaded_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type CreateDocumentRequest struct {
	FileID         string  `json:"file_id" validate:"required,uuid"`
	DocumentKindID int16   `json:"document_kind_id" validate:"required"`
	PreturnoID     *string `json:"preturno_id,omitempty"`
	CaseID         *string `json:"case_id,omitempty"`
	CaseEventID    *string `json:"case_event_id,omitempty"`
	Notes          *string `json:"notes,omitempty"`
	UploadedBy     *string `json:"uploaded_by,omitempty"`
}

type GetDocumentByIDRequest struct {
	ID string `params:"id" validate:"required,uuid"`
}

type GetDocumentsByFileIDRequest struct {
	FileID string `query:"file_id" validate:"required,uuid"`
}

type GetDocumentsByCaseRequest struct {
	CaseID string `query:"case_id" validate:"required,uuid"`
}

type DeleteDocumentRequest struct {
	ID string `params:"id" validate:"required,uuid"`
}
