package dto

import "time"

type ServiceVariant struct {
	ID            int16     `db:"id" json:"id"`
	ServiceTypeID int16     `db:"service_type_id" json:"service_type_id"`
	ModalityID    int16     `db:"modality_id" json:"modality_id"`
	Code          string    `db:"code" json:"code"`
	Name          string    `db:"name" json:"name"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type ServiceVariantCreateRequest struct {
	ServiceTypeID int16  `json:"service_type_id" validate:"required"`
	ModalityID    int16  `json:"modality_id" validate:"required"`
	Code          string `json:"code" validate:"required"`
	Name          string `json:"name" validate:"required"`
	IsActive      *bool  `json:"is_active,omitempty"`
}

type ServiceVariantUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}
