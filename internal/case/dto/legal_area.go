package dto

import "time"

// CREATE TABLE IF NOT EXISTS sicou.legal_area (
//   id          smallserial PRIMARY KEY,
//   code        text NOT NULL UNIQUE,
//   name        text NOT NULL,
//   is_active   boolean NOT NULL DEFAULT true,
//   created_at  timestamptz NOT NULL DEFAULT now(),
//   updated_at  timestamptz NOT NULL DEFAULT now()
// );

type LegalArea struct {
	ID        int16     `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateLegalAreaRequest struct {
	Code     string `json:"code" validate:"required"`
	Name     string `json:"name" validate:"required"`
	IsActive *bool  `json:"is_active,omitempty"`
}

type UpdateLegalAreaRequest struct {
	Code     *string `json:"code,omitempty"`
	Name     *string `json:"name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}
