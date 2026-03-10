package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq" 
)

// EmailTemplate
type EmailTemplate struct {
	Code       string         `db:"code" json:"code"`
	SubjectTpl string         `db:"subject_tpl" json:"subject_tpl"`
	BodyTpl    string         `db:"body_tpl" json:"body_tpl"`
	CCEmails   pq.StringArray `db:"cc_emails" json:"cc_emails,omitempty"` // text[]
	IsActive   bool           `db:"is_active" json:"is_active"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at" json:"updated_at"`
}

// EmailOutbox
type EmailOutbox struct {
	ID           uuid.UUID      `db:"id" json:"id"`
	TemplateCode *string        `db:"template_code" json:"template_code,omitempty"`
	ToEmails     pq.StringArray `db:"to_emails" json:"to_emails"`
	CCEmails     pq.StringArray `db:"cc_emails" json:"cc_emails,omitempty"`
	Subject      string         `db:"subject" json:"subject"`
	Body         string         `db:"body" json:"body"`
	Status       string         `db:"status" json:"status"` // PENDIENTE, ENVIANDO, ENVIADO, FALLIDO
	Attempts     int            `db:"attempts" json:"attempts"`
	LastError    *string        `db:"last_error" json:"last_error,omitempty"`
	ScheduledAt  time.Time      `db:"scheduled_at" json:"scheduled_at"`
	SentAt       *time.Time     `db:"sent_at" json:"sent_at,omitempty"`

	//Trazabilidad
	RelatedCaseID *uuid.UUID `db:"related_case_id" json:"related_case_id,omitempty"`
	RelatedTurnID *uuid.UUID `db:"related_turn_id" json:"related_turn_id,omitempty"`
	RelatedTermID *uuid.UUID `db:"related_term_id" json:"related_term_id,omitempty"`
	RelatedPQRSID *uuid.UUID `db:"related_pqrs_id" json:"related_pqrs_id,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateEmailOutboxRequest struct {
	TemplateCode *string   `json:"template_code,omitempty" db:"template_code"`
	ToEmails     []string  `json:"to_emails" db:"to_emails"`
	CCEmails     []string  `json:"cc_emails,omitempty" db:"cc_emails"`
	Subject      string    `json:"subject" db:"subject"`
	Body         string    `json:"body" db:"body"`
	ScheduledAt  time.Time `json:"scheduled_at" db:"scheduled_at"`
	
	RelatedCaseID *uuid.UUID `json:"related_case_id,omitempty" db:"related_case_id"`
	RelatedPQRSID *uuid.UUID `json:"related_pqrs_id,omitempty" db:"related_pqrs_id"`
}