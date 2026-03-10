package dto

import (
	"time"

	"github.com/google/uuid"
)

type CitizenSocioeconomic struct {
	ID                 uuid.UUID  `db:"id" json:"id"`
	CitizenID          uuid.UUID  `db:"citizen_id" json:"citizen_id"`
	HousingTypeID      *int16     `db:"housing_type_id" json:"housing_type_id,omitempty"`
	Stratum            *int16     `db:"stratum" json:"stratum,omitempty"`
	SisbenCategory     *string    `db:"sisben_category" json:"sisben_category,omitempty"`
	SisbenScore        *float64   `db:"sisben_score" json:"sisben_score,omitempty"` // numeric(6,2)
	VerificationStatus string     `db:"verification_status" json:"verification_status"`
	Observation        *string    `db:"observation" json:"observation,omitempty"`
	SupportDocumentID  *uuid.UUID `db:"support_document_id" json:"support_document_id,omitempty"`

	CreatedBy *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateCitizenSocioeconomicRequest struct {
	CitizenID       uuid.UUID `json:"citizen_id" db:"citizen_id"`
	HousingTypeID   *int16    `json:"housing_type_id,omitempty" db:"housing_type_id"`
	Stratum         *int16    `json:"stratum,omitempty" db:"stratum"`
	SisbenCategory  *string   `json:"sisben_category,omitempty" db:"sisben_category"`
	SisbenScore     *float64  `json:"sisben_score,omitempty" bd:"sisben_score"`
	Observation     *string   `json:"observation,omitempty" db:"observation"`
}

type UpdateCitizenSocioeconomicRequest struct {
	HousingTypeID      *int16     `json:"housing_type_id,omitempty" db:"housing_type_id"`
	Stratum            *int16     `json:"stratum,omitempty" db:"stratum"`
	SisbenCategory     *string    `json:"sisben_category,omitempty" db:"sisben_category"`
	SisbenScore        *float64   `json:"sisben_score,omitempty" db:"sisben_score"`
	VerificationStatus *string    `json:"verification_status,omitempty" db:"observation"`
	Observation        *string    `json:"observation,omitempty"`
	SupportDocumentID  *uuid.UUID `json:"support_document_id,omitempty" db:"support_document_id"`
}