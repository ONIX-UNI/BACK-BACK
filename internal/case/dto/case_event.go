package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CaseEvent struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	CaseID    uuid.UUID  `db:"case_id" json:"case_id"`
	EventType string     `db:"event_type" json:"event_type"`
	Title     *string    `db:"title" json:"title,omitempty"`
	Notes     *string    `db:"notes" json:"notes,omitempty"`
	Payload   []byte     `db:"payload" json:"payload,omitempty"`
	CreatedBy *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

type CreateCaseEventRequest struct {
	CaseID    uuid.UUID       `json:"case_id" validate:"required"`
	EventType string          `json:"event_type" validate:"required"`
	Title     *string         `json:"title,omitempty"`
	Notes     *string         `json:"notes,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	CreatedBy uuid.UUID       `json:"created_by,omitempty"`
}

type UpdateCaseEventRequest struct {
	EventType *string         `json:"created_by,omitempty"`
	Title     *string         `json:"title,omitempty"`
	Notes     *string         `json:"notes,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type GetByIdCaseEventRequest struct {
	ID uuid.UUID `json:"id"`
}

type DeleteCaseEventRequest struct {
	ID uuid.UUID `json:"id"`
}
