package dto

import (
	"time"

	"github.com/google/uuid"
)

type SatisfactionSurvey struct {
	CaseID    uuid.UUID `db:"case_id" json:"case_id"`
	ChannelID *int16    `db:"channel_id" json:"channel_id,omitempty"`
	Score     int16     `db:"score" json:"score"` // 1 a 5
	Comments  *string   `db:"comments" json:"comments,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type SatisfactionException struct {
	CaseID       uuid.UUID `db:"case_id" json:"case_id"`
	Reason       string    `db:"reason" json:"reason"`
	AuthorizedBy uuid.UUID `db:"authorized_by" json:"authorized_by"`
	AuthorizedAt time.Time `db:"authorized_at" json:"authorized_at"`
}

type CreateSurveyRequest struct {
	CaseID    uuid.UUID `json:"case_id" validate:"required"`
	ChannelID *int16    `json:"channel_id,omitempty"`
	Score     int16     `json:"score" validate:"required,min=1,max=5"`
	Comments  *string   `json:"comments,omitempty"`
}

type CreateExceptionRequest struct {
	CaseID uuid.UUID `json:"case_id" validate:"required"`
	Reason string    `json:"reason" validate:"required,min=10"`
}