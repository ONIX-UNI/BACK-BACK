package dto

import (
	"time"

	"github.com/google/uuid"
)

type CitizenEligibility struct {
	ID                    uuid.UUID  `db:"id" json:"id"`
	CitizenID             uuid.UUID  `db:"citizen_id" json:"citizen_id"`
	Criterion             string     `db:"criterion" json:"criterion"`
	IsEligible            bool       `db:"is_eligible" json:"is_eligible"`
	SupportDocumentID     *uuid.UUID `db:"support_document_id" json:"support_document_id,omitempty"`
	Observation           *string    `db:"observation" json:"observation,omitempty"`
	ExceptionAuthorized   bool       `db:"exception_authorized" json:"exception_authorized"`
	ExceptionAuthorizedBy *uuid.UUID `db:"exception_authorized_by" json:"exception_authorized_by,omitempty"`
	ExceptionAuthorizedAt *time.Time `db:"exception_authorized_at" json:"exception_authorized_at,omitempty"`

	CreatedBy *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateCitizenEligibilityRequest struct {
	CitizenID         uuid.UUID  `json:"citizen_id" db:"citizen_id"`
	Criterion         string     `json:"criterion" db:"criterion"`
	IsEligible        bool       `json:"is_eligible" db:"is_eligible"`
	SupportDocumentID *uuid.UUID `json:"support_document_id,omitempty" db:"support_document_id"`
	Observation       *string    `json:"observation,omitempty" db:"observation"`
}

type UpdateCitizenEligibilityRequest struct {
	IsEligible            *bool      `json:"is_eligible,omitempty" db:"is_eligible"`
	Observation           *string    `json:"observation,omitempty" db:"observation"`
	ExceptionAuthorized   *bool      `json:"exception_authorized,omitempty" db:"exception_authorized"`
	ExceptionAuthorizedBy *uuid.UUID `json:"exception_authorized_by,omitempty" db:"exception_authorized_by"`
}