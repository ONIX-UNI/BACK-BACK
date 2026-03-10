package dto

import (
	"time"

	"github.com/google/uuid"
)

type Turn struct {
	ID         uuid.UUID `db:"id" json:"id"`
	PreturnoID uuid.UUID `db:"preturno_id" json:"preturno_id"`

	TurnDate    time.Time `db:"turn_date" json:"turn_date"`
	Consecutive int       `db:"consecutive" json:"consecutive"`

	Status       string     `db:"status" json:"status"`
	Priority     int16      `db:"priority" json:"priority"`
	EscritorioID *uuid.UUID `db:"escritorio_id" json:"escritorio_id,omitempty"`

	CreatedBy *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`

	CalledAt   *time.Time `db:"called_at" json:"called_at,omitempty"`
	AttendedAt *time.Time `db:"attended_at" json:"attended_at,omitempty"`
	FinishedAt *time.Time `db:"finished_at" json:"finished_at,omitempty"`
}

type CreateTurnRequest struct {
	PreturnoID uuid.UUID  `json:"preturno_id"`
	TurnDate   *time.Time `json:"turn_date,omitempty"`
	Priority   *int16     `json:"priority,omitempty"`
	CreatedBy  *uuid.UUID `json:"created_by,omitempty"`
	StudentID  *uuid.UUID `json:"student_id,omitempty"`
	AssignedTo *uuid.UUID `json:"assigned_to,omitempty"`
}

type UpdateTurnRequest struct {
	Status     *string    `json:"status,omitempty"`
	Priority   *int16     `json:"priority,omitempty"`
	CalledAt   *time.Time `json:"called_at,omitempty"`
	AttendedAt *time.Time `json:"attended_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

type SetTurnDesktopRequest struct {
	TurnID       uuid.UUID `json:"turn_id"`
	EscritorioID uuid.UUID `json:"escritorio_id"`
	Reason       *string   `json:"reason,omitempty"`
}

type SetTurnEscritorioResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
