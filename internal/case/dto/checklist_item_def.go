package dto

import (
	"time"
)

type ChecklistItemDef struct {
	ID                     int16     `json:"id" db:"id"`
	Code                   string    `json:"code" db:"code"`
	Name                   string    `json:"name" db:"name"`
	Description            *string   `json:"description,omitempty" db:"description"`
	RequiresDocumentKindID *int16    `json:"requires_document_kind_id,omitempty" db:"requires_document_kind_id"`
	RequiredForClose       bool      `json:"required_for_close" db:"required_for_close"`
	IsActive               bool      `json:"is_active" db:"is_active"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

type ChecklistItemDefCreateRequest struct {
	Code                   string  `json:"code" db:"code"`
	Name                   string  `json:"name" db:"name"`
	Description            *string `json:"description,omitempty" db:"description"`
	RequiresDocumentKindID *int16  `json:"requires_document_kind_id,omitempty" db:"requires_document_kind_id"`
	RequiredForClose       bool    `json:"required_for_close"`
	IsActive               *bool   `json:"is_active,omitempty" db:"is_active"`
}

type ChecklistItemDefUpdateRequest struct {
	ID                     int16   `json:"id" db:"id"`
	Code                   *string `json:"code,omitempty" db:"code"`
	Name                   *string `json:"name,omitempty" db:"name"`
	Description            *string `json:"description,omitempty" db:"description"`
	RequiresDocumentKindID *int16  `json:"requires_document_kind_id,omitempty" db:"requires_document_kind_id"`
	RequiredForClose       *bool   `json:"required_for_close,omitempty" db:"required_for_close"`
	IsActive               *bool   `json:"is_active,omitempty" db:"is_active"`
}
