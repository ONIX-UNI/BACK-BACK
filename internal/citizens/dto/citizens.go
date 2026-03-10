package dto

import (
	"time"

	"github.com/google/uuid"
)

type Citizen struct {
	ID             uuid.UUID  `db:"id"`
	DocumentTypeID int16      `db:"document_type_id"`
	DocumentNumber string     `db:"document_number"`
	FullName       string     `db:"full_name"`
	BirthDate      *time.Time `db:"birth_date"`

	PhoneMobile *string `db:"phone_mobile"`
	Email       *string `db:"email"`
	Address     *string `db:"address"`

	CreatedBy *uuid.UUID `db:"created_by"`
	UpdatedBy *uuid.UUID `db:"updated_by"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type CreateCitizenRequest struct {
	DocumentTypeID int16      `json:"document_type_id" validate:"required"`
	DocumentNumber string     `json:"document_number" validate:"required"`
	FullName       string     `json:"full_name" validate:"required"`
	BirthDate      *time.Time `json:"birth_date,omitempty"`

	PhoneMobile *string `json:"phone_mobile,omitempty"`
	Email       *string `json:"email,omitempty"`
	Address     *string `json:"address,omitempty"`
}

type UpdateCitizenRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`

	DocumentTypeID *int16     `json:"document_type_id,omitempty"`
	DocumentNumber *string    `json:"document_number,omitempty"`
	FullName       *string    `json:"full_name,omitempty"`
	BirthDate      *time.Time `json:"birth_date,omitempty"`

	PhoneMobile *string `json:"phone_mobile,omitempty"`
	Email       *string `json:"email,omitempty"`
	Address     *string `json:"address,omitempty"`
}

type GetCitizenByIDRequest struct {
	ID uuid.UUID `params:"id" validate:"required,uuid"`
}

type DeleteCitizenRequest struct {
	ID uuid.UUID `params:"id" validate:"required,uuid"`
}
